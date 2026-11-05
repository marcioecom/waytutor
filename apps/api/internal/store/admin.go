package store

import (
	"context"
	"fmt"
	"time"
)

type AdminAccount struct {
	ID        string
	Email     string
	Name      string
	Status    string
	CreatedAt time.Time
}

func (s *Store) ListAccounts(ctx context.Context, status, email string, limit int) ([]AdminAccount, error) {
	const query = `
		SELECT id, email, name, status, created_at
		FROM accounts
		WHERE ($1 = '' OR status = $1)
		  AND ($2 = '' OR lower(email) = lower($2))
		ORDER BY created_at DESC
		LIMIT $3`

	rows, err := s.pool.Query(ctx, query, status, email, limit)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []AdminAccount
	for rows.Next() {
		var a AdminAccount
		if err := rows.Scan(&a.ID, &a.Email, &a.Name, &a.Status, &a.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan account: %w", err)
		}
		accounts = append(accounts, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read accounts: %w", err)
	}
	return accounts, nil
}

func (s *Store) SetAccountStatus(ctx context.Context, id, status string) error {
	const query = `UPDATE accounts SET status = $2, updated_at = now() WHERE id = $1`

	tag, err := s.pool.Exec(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("set account status: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) RevokeSessions(ctx context.Context, accountID string) (int64, error) {
	const query = `DELETE FROM access_sessions WHERE account_id = $1`

	tag, err := s.pool.Exec(ctx, query, accountID)
	if err != nil {
		return 0, fmt.Errorf("revoke sessions: %w", err)
	}
	return tag.RowsAffected(), nil
}
