-- +goose Up
CREATE TABLE accounts (
    id         uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    google_sub text NOT NULL UNIQUE,
    email      text NOT NULL,
    name       text NOT NULL DEFAULT '',
    picture    text NOT NULL DEFAULT '',
    status     text NOT NULL DEFAULT 'pending'
               CHECK (status IN ('pending', 'active', 'paused')),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX accounts_email_key ON accounts (lower(email));

-- +goose Down
DROP TABLE accounts;
