package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

var ErrNotFound = errors.New("not found")

func (s *Store) CreateAccessSession(ctx context.Context, accountID, tokenHash string, expiresAt time.Time, userAgent string) error {
	const query = `
		INSERT INTO access_sessions (token_hash, account_id, user_agent, expires_at)
		VALUES ($1, $2, $3, $4)`

	if _, err := s.pool.Exec(ctx, query, tokenHash, accountID, userAgent, expiresAt); err != nil {
		return fmt.Errorf("create access session: %w", err)
	}
	return nil
}

func (s *Store) AccountBySession(ctx context.Context, tokenHash string, renewTo time.Time) (Account, error) {
	const query = `
		UPDATE access_sessions AS s
		SET last_seen_at = now(), expires_at = $2
		FROM accounts AS a
		WHERE s.token_hash = $1
		  AND s.account_id = a.id
		  AND s.expires_at > now()
		RETURNING a.id, a.email, a.name, a.status`

	var a Account
	err := s.pool.QueryRow(ctx, query, tokenHash, renewTo).
		Scan(&a.ID, &a.Email, &a.Name, &a.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return Account{}, ErrNotFound
	}
	if err != nil {
		return Account{}, fmt.Errorf("load access session: %w", err)
	}
	return a, nil
}

func (s *Store) DeleteAccessSession(ctx context.Context, tokenHash string) error {
	const query = `DELETE FROM access_sessions WHERE token_hash = $1`

	if _, err := s.pool.Exec(ctx, query, tokenHash); err != nil {
		return fmt.Errorf("delete access session: %w", err)
	}
	return nil
}
