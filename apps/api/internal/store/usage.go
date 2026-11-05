package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

var ErrQuotaExceeded = errors.New("quota exceeded")

type Usage struct {
	Chars    int64
	Requests int64
}

func (s *Store) ReserveChars(ctx context.Context, accountID, period string, chars, perAccount, global int64) (Usage, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return Usage{}, fmt.Errorf("begin quota reservation: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const insert = `INSERT INTO usage (account_id, period) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	if _, err := tx.Exec(ctx, insert, accountID, period); err != nil {
		return Usage{}, fmt.Errorf("create usage row: %w", err)
	}

	var total int64
	const globalSum = `SELECT COALESCE(SUM(tts_chars), 0) FROM usage WHERE period = $1`
	if err := tx.QueryRow(ctx, globalSum, period).Scan(&total); err != nil {
		return Usage{}, fmt.Errorf("sum period usage: %w", err)
	}
	if total+chars > global {
		return Usage{}, ErrQuotaExceeded
	}

	const reserve = `
		UPDATE usage
		SET tts_chars = tts_chars + $3, tts_requests = tts_requests + 1
		WHERE account_id = $1 AND period = $2 AND tts_chars + $3 <= $4
		RETURNING tts_chars, tts_requests`

	var usage Usage
	err = tx.QueryRow(ctx, reserve, accountID, period, chars, perAccount).
		Scan(&usage.Chars, &usage.Requests)
	if errors.Is(err, pgx.ErrNoRows) {
		return Usage{}, ErrQuotaExceeded
	}
	if err != nil {
		return Usage{}, fmt.Errorf("reserve chars: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Usage{}, fmt.Errorf("commit quota reservation: %w", err)
	}
	return usage, nil
}

func (s *Store) RefundChars(ctx context.Context, accountID, period string, chars int64) error {
	const query = `
		UPDATE usage
		SET tts_chars = GREATEST(tts_chars - $3, 0),
		    tts_requests = GREATEST(tts_requests - 1, 0)
		WHERE account_id = $1 AND period = $2`

	if _, err := s.pool.Exec(ctx, query, accountID, period, chars); err != nil {
		return fmt.Errorf("refund chars: %w", err)
	}
	return nil
}

func (s *Store) ReserveSummary(ctx context.Context, accountID, period string, limit int64) error {
	const query = `
		INSERT INTO usage (account_id, period, summaries) VALUES ($1, $2, 1)
		ON CONFLICT (account_id, period) DO UPDATE
		SET summaries = usage.summaries + 1
		WHERE usage.summaries + 1 <= $3`

	tag, err := s.pool.Exec(ctx, query, accountID, period, limit)
	if err != nil {
		return fmt.Errorf("reserve summary: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrQuotaExceeded
	}
	return nil
}

func (s *Store) RefundSummary(ctx context.Context, accountID, period string) error {
	const query = `
		UPDATE usage
		SET summaries = GREATEST(summaries - 1, 0)
		WHERE account_id = $1 AND period = $2`

	if _, err := s.pool.Exec(ctx, query, accountID, period); err != nil {
		return fmt.Errorf("refund summary: %w", err)
	}
	return nil
}
