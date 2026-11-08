package store

import (
	"context"
	"fmt"
	"time"
)

type DailyUsage struct {
	Day      string
	Chars    int64
	Requests int64
}

type AccountUsage struct {
	Email    string
	Name     string
	Chars    int64
	Requests int64
}

type UsageTotals struct {
	Chars     int64
	Requests  int64
	Summaries int64
}

type AdminUsage struct {
	Day           string
	Period        string
	Today         UsageTotals
	TodayAccounts int64
	Month         UsageTotals
	Daily         []DailyUsage
	Top           []AccountUsage
	Active        int64
	Pending       int64
	Paused        int64
}

// AdminUsage reads the panel for one period, with the last `days` days of the
// daily rollup.
func (s *Store) AdminUsage(ctx context.Context, period string, days int) (AdminUsage, error) {
	now := time.Now().In(usageDay)
	usage := AdminUsage{Day: currentDay(), Period: period}

	const today = `
		SELECT COALESCE(SUM(tts_chars), 0), COALESCE(SUM(tts_requests), 0),
		       COALESCE(SUM(summaries), 0), COUNT(*)
		FROM usage_daily WHERE day = $1::date`

	if err := s.pool.QueryRow(ctx, today, usage.Day).Scan(
		&usage.Today.Chars, &usage.Today.Requests, &usage.Today.Summaries, &usage.TodayAccounts,
	); err != nil {
		return AdminUsage{}, fmt.Errorf("sum daily usage: %w", err)
	}

	const month = `
		SELECT COALESCE(SUM(tts_chars), 0), COALESCE(SUM(tts_requests), 0), COALESCE(SUM(summaries), 0)
		FROM usage WHERE period = $1`

	if err := s.pool.QueryRow(ctx, month, period).Scan(
		&usage.Month.Chars, &usage.Month.Requests, &usage.Month.Summaries,
	); err != nil {
		return AdminUsage{}, fmt.Errorf("sum monthly usage: %w", err)
	}

	const statuses = `
		SELECT COUNT(*) FILTER (WHERE status = 'active'),
		       COUNT(*) FILTER (WHERE status = 'pending'),
		       COUNT(*) FILTER (WHERE status = 'paused')
		FROM accounts`

	if err := s.pool.QueryRow(ctx, statuses).Scan(&usage.Active, &usage.Pending, &usage.Paused); err != nil {
		return AdminUsage{}, fmt.Errorf("count account statuses: %w", err)
	}

	start := now.AddDate(0, 0, -(days - 1)).Format("2006-01-02")
	daily, err := s.dailySeries(ctx, start, usage.Day, now, days)
	if err != nil {
		return AdminUsage{}, err
	}
	usage.Daily = daily

	const top = `
		SELECT a.email, a.name, u.tts_chars, u.tts_requests
		FROM usage u
		JOIN accounts a ON a.id = u.account_id
		WHERE u.period = $1 AND u.tts_chars > 0
		ORDER BY u.tts_chars DESC
		LIMIT $2`

	rows, err := s.pool.Query(ctx, top, period, 5)
	if err != nil {
		return AdminUsage{}, fmt.Errorf("top accounts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var account AccountUsage
		if err := rows.Scan(&account.Email, &account.Name, &account.Chars, &account.Requests); err != nil {
			return AdminUsage{}, fmt.Errorf("scan top account: %w", err)
		}
		usage.Top = append(usage.Top, account)
	}
	if err := rows.Err(); err != nil {
		return AdminUsage{}, fmt.Errorf("read top accounts: %w", err)
	}
	return usage, nil
}

func (s *Store) dailySeries(ctx context.Context, start, end string, now time.Time, days int) ([]DailyUsage, error) {
	const query = `
		SELECT day, SUM(tts_chars), SUM(tts_requests)
		FROM usage_daily
		WHERE day >= $1::date AND day <= $2::date
		GROUP BY day
		ORDER BY day`

	rows, err := s.pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("daily usage: %w", err)
	}
	defer rows.Close()

	byDay := make(map[string]DailyUsage, days)
	for rows.Next() {
		var day time.Time
		var entry DailyUsage
		if err := rows.Scan(&day, &entry.Chars, &entry.Requests); err != nil {
			return nil, fmt.Errorf("scan daily usage: %w", err)
		}
		entry.Day = day.Format("2006-01-02")
		byDay[entry.Day] = entry
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read daily usage: %w", err)
	}

	// Discord charts need one point per day, so empty days become zeros.
	series := make([]DailyUsage, 0, days)
	for offset := days - 1; offset >= 0; offset-- {
		day := now.AddDate(0, 0, -offset).Format("2006-01-02")
		entry, ok := byDay[day]
		if !ok {
			entry = DailyUsage{Day: day}
		}
		series = append(series, entry)
	}
	return series, nil
}
