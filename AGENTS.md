# WayTutor engineering guide

Conventions for this repository. Read this before adding or changing code.

## Language

English everywhere: identifiers, comments, docs, commit messages, branch names,
error codes, log fields. Portuguese appears in two places only:

- extension UI copy (`apps/extension/modules/*/ui`)
- the map from backend error code to user message (`modules/*/messages.ts`)

An English string that reaches a user is a bug. A Portuguese string inside Go is
a bug.

## Layout

```
apps/api/         Go service, deployed on Fly.io
apps/extension/   Chrome extension, built by WXT
docs/adr/         decisions worth remembering
```

The repository root holds orchestration only: `Makefile`, docs, glossary.
`apps/bot` and `apps/site` are planned; they are created when they have a job.

## The bar

Review criteria, not suggestions.

- A comment explains a why the code cannot show. Never what the code does.
- No `pkg/`, no barrel files, no `index.ts`, no `utils.ts`.
- No abstraction with one implementation. The exceptions are real boundaries:
  speech provider, audio cache, database.
- A file over 200 lines is two responsibilities in one file.
- No configuration for something that never varies. The voice list is a
  catalogue, so it lives in code.
- `ui/` never imports `chrome`, `fetch`, `zod` or `react-hook-form`. That is
  `hooks/`.
- Tests cover rules: extraction, reading state, quota, auth, cache. Nothing else.
- An ADR only when the decision is hard to reverse, surprising without context,
  and the result of a real trade-off.

## Go

- Standard library first: `net/http` with ServeMux patterns. No framework, no ORM.
- Everything lives in `internal/`, binaries in `cmd/<name>/main.go`.
- `main` wires and returns; `os.Exit` appears only in `main`.
- Logging is `log/slog`, global. Call sites use `slog.Info` and `slog.Error`.
- Config uses struct tags (`env`, `envDefault`, `validate`). A value enters the
  config when a feature reads it, not before.
- PostgreSQL through `pgx`, SQL written by hand in `internal/store`, one file per
  area. Migrations live in `internal/migrations` and are applied by `cmd/migrate`.
- Wrap errors with context: `fmt.Errorf("open store: %w", err)`.
- Reserve quota before the model call and refund when the provider fails.
- One wide event per request: method, path, status, duration, bytes, request id.

## HTTP API

- Versioned paths: `/v1/...`, stable error codes in English.
- An error body is `{"error":{"code":"quota_exceeded"}}`. There is no `message`.
- Auth is an opaque bearer token, never a cookie.
- Every body goes through `http.MaxBytesReader`.
- A Passage's audio is cached under the hash of everything that changes it:
  model, voice, instructions, format and text. A cache hit costs no quota.
- Audio responses carry `ETag` and `Cache-Control: immutable`.

## TypeScript and the extension

- WXT builds it. There is no build script of our own.
- `modules/<feature>/hooks` owns chrome messages, react-query and form state.
  `modules/<feature>/ui` renders and owns nothing else.
- react-query is the only cache. No store library, no Cache Storage layer.
- Validate at the boundary: `chrome.storage`, messages and API responses use the
  zod schema of the owning module.
- The service worker registers the side panel behaviour and holds no state.
- The offscreen document fetches and plays audio. Nothing sends base64.
- Colours and spacing come from Tailwind tokens in `style.css`.

## Before you commit

```sh
make check   # gofmt, go vet, go test
```
