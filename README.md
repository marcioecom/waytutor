# WayTutor

A Chrome extension that reads study material aloud, passage by passage, and
produces a spoken summary grounded in that material.

- `apps/extension` — WXT, side panel, content script, offscreen audio
- `apps/api` — Go service on Fly.io, PostgreSQL on Neon, models through Vercel AI
  Gateway

## Requirements

- Node 22+
- Go 1.26+
- A PostgreSQL database (Neon)

## Run

```sh
cp apps/api/.env.example apps/api/.env   # fill DATABASE_URL
make migrate
make dev                                 # api on :8080

cd apps/extension && npm install && npm run dev
```

Load `apps/extension/.output/chrome-mv3` as an unpacked extension.

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
