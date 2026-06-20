# Mugo Server Agent Guidelines

The server is a Go 1.24 API using Huma v2, Chi, pgx/sqlc, Postgres migrations, and Google ADK/Gemini agents.

## Commands
- Run API: `make server`
- Run API with hot reload + docker lifecycle: `make dev` (uses [air](https://github.com/air-verse/air) — auto-rebuilds on save, starts/stops docker containers)
- Build API: `make build`
- Run Go tests: `make test-server`
- Start dependencies: `make docker`
- Run migrations: `make migrate-up`
- Generate sqlc code: `make sqlc`
- Tidy dependencies: `make tidy`

## Layout
- Entrypoints: `cmd/api/main.go`, `cmd/migrate/main.go`
- Routes: `internal/routes/`
- Domain models: `internal/models/`
- ADK runner wrapper: `internal/adk/`
- Agent definitions: `internal/agents/`
- DB connection/lifecycle: `internal/db/`
- Repositories: `internal/db/repository/`
- SQL queries: `internal/db/queries/`
- Migrations: `internal/db/migrations/`
- sqlc output: `internal/db/dbgenerated/`
- Integration tests: `tests/integration/`

## API Rules
- Define Huma request/response types with a nested `Body` field for JSON bodies.
- Use descriptive `OperationID`s and `Tags`; Orval groups the mobile client by tags.
- Return Huma structured errors such as `huma.Error404NotFound`, `huma.Error400BadRequest`, or `huma.Error500InternalServerError`.
- Avoid generic error returns from route handlers unless they are already structured Huma errors.

## Database Rules
- Access the database through the `db.DBProvider` interface and `GetDB(provider)` helper in `internal/routes/util.go`.
- Route handlers should use repository methods, not sqlc generated queries directly.
- Direct sqlc usage belongs in `internal/db/repository/`.
- Use `database.WithTx(ctx, ...)` for multi-step writes that must be atomic.
- Do not manually manage pgx pools in route handlers; rely on `LazyDatabase` initialized in `cmd/api/main.go`.
- Integration tests skip when `DATABASE_URL` is unset.

## SQL and Migrations
- Add or change SQL in `internal/db/queries/*.sql`; run `make sqlc` afterward.
- Add schema changes as migrations in `internal/db/migrations/` and verify with `make migrate-up` when the DB is available.
- Never hand-edit `internal/db/dbgenerated/`.

## ADK Rules
- Create agents in `internal/agents/` and runners through `internal/adk/`.
- Use `llmagent.New()` with Gemini models created via `agents.NewGeminiModel()` unless an existing pattern says otherwise.
- `GOOGLE_API_KEY` is required for Gemini-backed agents; route handlers must tolerate nil runners by returning structured Huma errors.
- Session management currently uses the ADK session service created in `cmd/api/main.go`.

## Known Constraints
- Do not assume authentication or row-level security exists until `server/TODOS/001-authentication.md` and related TODOs are completed.
- Several backend TODOs document known design and correctness issues. Check `server/TODOS/` before making broad architectural assumptions.
