# WayTutor

A Chrome extension that reads study material aloud, passage by passage, and
produces a spoken summary grounded in that material.

- `apps/extension` — WXT, side panel, content script, offscreen audio
- `apps/api` — Go service on Fly.io, PostgreSQL on Neon, models through Vercel AI
  Gateway

## Requirements

- Docker, for the local database and stack
- Node 22+
- Go 1.26+
- A PostgreSQL database (Neon) for anything beyond local work

## Run

The whole stack — database, migrations, api on `:8080`:

```sh
make up
make down
```

Just the api, against the compose database:

```sh
docker compose up -d db
cp apps/api/.env.example apps/api/.env
make migrate
make dev
```

The extension:

```sh
cd apps/extension && npm install && npm run dev
```

Load `apps/extension/.output/chrome-mv3` as an unpacked extension. The panel
talks to `WXT_API_URL`, which defaults to `http://localhost:8080`; set it in
`apps/extension/.env` for a production build.

For an external database, point `DATABASE_URL` at Neon with `sslmode=require`.
Keep the local Postgres major version in step with Neon's.

## Login

The extension signs in through Google. In the Google Cloud console, create an
OAuth web client and add both redirect URIs:

```
http://localhost:8080/auth/google/callback
https://waytutor-api.marcio.run/auth/google/callback
```

Then fill `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET`, `OAUTH_STATE_SECRET` (16
characters or more) and `BASE_URL` in `apps/api/.env`.

New accounts are created with `status = 'pending'` and cannot use the API until
approved:

```sh
psql "$DATABASE_URL" -c "UPDATE accounts SET status = 'active', updated_at = now() WHERE email = 'someone@example.com'"
```

## API

| Method | Path | Auth |
| --- | --- | --- |
| GET | `/healthz` | — |
| GET | `/auth/google?redirect_uri=` | — |
| GET | `/auth/google/callback` | — |
| POST | `/v1/auth/exchange` | — |
| POST | `/v1/auth/logout` | bearer |
| GET | `/v1/me` | bearer |
| POST | `/v1/tts` | bearer |

Errors are `{"error":{"code":"..."}}` with an English code: `invalid_request`,
`invalid_redirect_uri`, `invalid_state`, `invalid_code`, `oauth_failed`,
`unauthenticated`, `approval_required`, `account_paused`, `empty_text`,
`text_too_long`, `unknown_voice`, `quota_exceeded`, `rate_limited`,
`provider_auth`, `provider_quota`, `provider_error`, `internal_error`.

`POST /v1/tts` takes `{"text": "...", "voice": "Kore"}` and returns the WAV
bytes. The same text, voice and instructions always hit the same cache object,
and a cache hit does not consume quota.

## Audio cache

Without `R2_ENDPOINT`, `R2_BUCKET`, `R2_ACCESS_KEY_ID` and `R2_SECRET_ACCESS_KEY`
the server caches audio in `CACHE_DIR`. Production sets all four and the machine
becomes stateless.

## Check

```sh
make check
```

## Deploy

```sh
cd apps/api && fly deploy
```

`fly.toml` runs `cmd/migrate` as the release command, so migrations land before
the new release serves traffic. Secrets live in `fly secrets`, never in the
image.
