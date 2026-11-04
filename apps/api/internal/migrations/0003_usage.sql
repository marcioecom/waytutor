-- +goose Up
CREATE TABLE usage (
    account_id   uuid NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    period       text NOT NULL,
    tts_chars    bigint NOT NULL DEFAULT 0,
    tts_requests bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (account_id, period)
);

CREATE INDEX usage_period_idx ON usage (period);

-- +goose Down
DROP TABLE usage;
