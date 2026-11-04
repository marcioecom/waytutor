package store

import (
	"context"
	"fmt"
)

const (
	StatusPending = "pending"
	StatusActive  = "active"
	StatusPaused  = "paused"
)

type Account struct {
	ID     string
	Email  string
	Name   string
	Status string
}

type Identity struct {
	GoogleSub string
	Email     string
	Name      string
	Picture   string
}

func (s *Store) UpsertAccount(ctx context.Context, id Identity) (Account, error) {
	const query = `
		INSERT INTO accounts (google_sub, email, name, picture)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (google_sub) DO UPDATE
		SET email = excluded.email,
		    name = excluded.name,
		    picture = excluded.picture,
		    updated_at = now()
		RETURNING id, email, name, status`

	var a Account
	err := s.pool.QueryRow(ctx, query, id.GoogleSub, id.Email, id.Name, id.Picture).
		Scan(&a.ID, &a.Email, &a.Name, &a.Status)
	if err != nil {
		return Account{}, fmt.Errorf("upsert account: %w", err)
	}
	return a, nil
}
