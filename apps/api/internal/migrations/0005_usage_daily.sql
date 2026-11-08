-- +goose Up
CREATE TABLE usage_daily (
    account_id   uuid NOT NULL REFERENCES accounts (id) ON DELETE CASCADE,
    day          date NOT NULL,
    tts_chars    bigint NOT NULL DEFAULT 0,
    tts_requests bigint NOT NULL DEFAULT 0,
    summaries    bigint NOT NULL DEFAULT 0,
    PRIMARY KEY (account_id, day)
);

CREATE INDEX usage_daily_day_idx ON usage_daily (day);

-- +goose Down
DROP TABLE usage_daily;
