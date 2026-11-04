package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func (s *Store) CreateLoginCode(ctx context.Context, codeHash, accountID string, expiresAt time.Time) error {
	const query = `
		INSERT INTO login_codes (code_hash, account_id, expires_at)
		VALUES ($1, $2, $3)`

	if _, err := s.pool.Exec(ctx, query, codeHash, accountID, expiresAt); err != nil {
		return fmt.Errorf("create login code: %w", err)
	}
	return nil
}

func (s *Store) ConsumeLoginCode(ctx context.Context, codeHash string) (Account, error) {
	const query = `
		UPDATE login_codes AS l
		SET used_at = now()
		FROM accounts AS a
		WHERE l.code_hash = $1
		  AND l.used_at IS NULL
		  AND l.expires_at > now()
		  AND a.id = l.account_id
		RETURNING a.id, a.email, a.name, a.status`

	var a Account
	err := s.pool.QueryRow(ctx, query, codeHash).
		Scan(&a.ID, &a.Email, &a.Name, &a.Status)
	if errors.Is(err, pgx.ErrNoRows) {
		return Account{}, ErrNotFound
	}
	if err != nil {
		return Account{}, fmt.Errorf("consume login code: %w", err)
	}
	return a, nil
}
