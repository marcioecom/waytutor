-- +goose Up
CREATE TABLE access_sessions (
    token_hash   text PRIMARY KEY,
    account_id   uuid NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    user_agent   text NOT NULL DEFAULT '',
    created_at   timestamptz NOT NULL DEFAULT now(),
    last_seen_at timestamptz NOT NULL DEFAULT now(),
    expires_at   timestamptz NOT NULL
);

CREATE INDEX access_sessions_account_id_idx ON access_sessions (account_id);

CREATE TABLE login_codes (
    code_hash  text PRIMARY KEY,
    account_id uuid NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL,
    used_at    timestamptz
);

-- +goose Down
DROP TABLE login_codes;
DROP TABLE access_sessions;
