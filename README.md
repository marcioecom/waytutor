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

Errors are `{"error":{"code":"..."}}`, and `code` is one of `invalid_request`,
`invalid_redirect_uri`, `invalid_state`, `invalid_code`, `oauth_failed`,
`unauthenticated`, `approval_required`, `account_paused`, `internal_error`.

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
