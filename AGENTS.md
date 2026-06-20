# Mugo / LazyFood Agent Guidelines

This repo contains a Go API/ADK server and an Expo React Native mobile app. Prefer the closest `AGENTS.md` file for stack-specific rules: `server/AGENTS.md` for backend work and `mobile/AGENTS.md` for mobile work.

## Current Commands
- Start dependencies: `make docker`
- Start API: `make server` (serves API on `http://localhost:8888`)
- Start API with hot reload + docker lifecycle: `make dev`
- Build API: `make build`
- Generate sqlc code: `make sqlc`
- Run migrations up: `make migrate-up`
- Run migrations down: `make migrate-down steps=1`
- Force migration version: `make migrate-force version=<version>`
- Start mobile dev server: `make mobile`
- Start Android emulator: `make emulator`
- Regenerate mobile API client: `make orval`
- Go tests: `make test-server`
- Mobile lint: `make lint-mobile`
- Mobile tests: `make test-mobile`

## Verification Matrix
- Go API change: `make test-server`
- Go build or wiring change: `make build`
- SQL query or migration change: `make sqlc && make test-server`
- Runtime DB behavior change: `make docker`, then `make migrate-up`, then relevant Go tests
- Mobile UI/state change: `make lint-mobile && make test-mobile`
- API contract change: start the API with `make server`, then run `make orval`, then `make lint-mobile && make test-mobile`

## Environment Notes
- Runtime commands generally use `infisical run` through the `Makefile`.
- Important env vars include `DATABASE_URL`, `GOOGLE_API_KEY`, `API_SERVER_URL`, `DB_PORT`, `WHISPER_PORT`, `TRANSCRIPTION_SERVER_URL`, and `PORT`.
- Do not commit secrets. Check existing config and command output before asking the user for env details.
- Docker Compose starts Postgres and Whisper. The Postgres service is named `db`; the Whisper service is named `whisper`.

## Generated Code Boundaries
- Do not hand-edit `server/internal/db/dbgenerated/`; edit SQL in `server/internal/db/queries/` or migrations, then run `make sqlc`.
- Do not hand-edit `mobile/lib/api/`; update Huma route schemas, run the API, then run `make orval`.
- Generated files may appear in diffs after `make sqlc` or `make orval`; keep them only when they are required by the source/schema change.

## API Contract Workflow
- Huma serves OpenAPI at `/openapi.json` when the API is running.
- Orval reads `process.env.API_SERVER_URL + "/openapi.json"` from `mobile/orval.config.ts`.
- If request/response structs, route paths, tags, or operation IDs change, regenerate the mobile client with `make orval`.

## Project Memory
- `server/TODOS/` contains known backend architecture and bug tasks. Check it before assuming AuthN, RLS, or related infrastructure exists.
- Relevant Notion pages are tracked in project/global agent memory, but local source and TODO files are the source of truth for implementation.
