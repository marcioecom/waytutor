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

Load `apps/extension/.output/chrome-mv3` as an unpacked extension.

For an external database, point `DATABASE_URL` at Neon with `sslmode=require`.
Keep the local Postgres major version in step with Neon's.

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
