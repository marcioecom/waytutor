package store

import (
	"context"
	"fmt"
	"time"
)

type AdminAccount struct {
	ID            string
	Email         string
	Name          string
	Status        string
	CreatedAt     time.Time
	UsageChars    int64
	UsageRequests int64
	Summaries     int64
	LastSeenAt    *time.Time
}

func (s *Store) ListAccounts(ctx context.Context, status, email, period string, limit int) ([]AdminAccount, error) {
	const query = `
		SELECT a.id, a.email, a.name, a.status, a.created_at,
		       COALESCE(u.tts_chars, 0), COALESCE(u.tts_requests, 0), COALESCE(u.summaries, 0),
		       seen.last_seen_at
		FROM accounts a
		LEFT JOIN usage u ON u.account_id = a.id AND u.period = $3
		LEFT JOIN LATERAL (
			SELECT MAX(last_seen_at) AS last_seen_at
			FROM access_sessions
			WHERE account_id = a.id
		) seen ON true
		WHERE ($1 = '' OR a.status = $1)
		  AND ($2 = '' OR lower(a.email) = lower($2))
		ORDER BY a.created_at DESC
		LIMIT $4`

	rows, err := s.pool.Query(ctx, query, status, email, period, limit)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	defer rows.Close()

	var accounts []AdminAccount
	for rows.Next() {
		var a AdminAccount
		if err := rows.Scan(
			&a.ID, &a.Email, &a.Name, &a.Status, &a.CreatedAt,
			&a.UsageChars, &a.UsageRequests, &a.Summaries, &a.LastSeenAt,
		); err != nil {
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
