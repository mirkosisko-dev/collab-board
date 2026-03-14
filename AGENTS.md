# AGENTS.md

This file contains guidelines for agentic coding agents working on this repository.

## Build, Lint, and Test Commands

### Backend (Go - apps/api/)

- `make setup` - Download Go dependencies + generate sqlc code
- `make build` - Build the Go binary
- `make run` - Build and run the application
- `make test` - Run all Go tests
- `make lint` - Run golint on Go code
- `make fmt` - Format Go code with `go fmt`
- `make check` - Run fmt + lint + test
- **Single test:** `go test -v ./path/to/package -run TestFunctionName`

### Frontend (React/TypeScript - apps/web/)

- `pnpm dev` - Start Vite dev server
- `pnpm build` - Build for production
- `pnpm lint` - Run ESLint
- `pnpm preview` - Preview production build

### Monorepo (root)

- `pnpm dev` - Run all apps in dev mode
- `pnpm build` - Build all apps
- `pnpm lint` - Lint all apps
- `pnpm format` - Format all code with Prettier

## Go Code Style Guidelines

### Package Structure & Imports

- Handlers: `internal/handlers/{feature}/route.go`
- Services: `internal/services/{service}/service.go`
- Import order: stdlib → third-party → internal (blank line separators)

### Naming Conventions

- Handler packages: lowercase (e.g., `board`, `user`)
- Handlers: `packageHandler` struct with `NewHandler()` constructor
- Handler methods: `handleActionName` (e.g., `handleCreateBoard`)
- Services: `NewService()` constructor, methods use receiver pattern
- Tests: `TestFunctionName`, helpers: `createRandomType()`
- Routes: `RegisterRoutes()` method that registers handlers

### Types & Data

- Use sqlc-generated types from `db/sqlc` package
- UUIDs: `pgtype.UUID{Bytes: id, Valid: true}` pattern
- Use `r.Context()` as first parameter to database queries

### Error Handling

- Always check and return errors
- Use `utils.WriteError(w, status, err)` for HTTP responses
- Use `utils.ParseJSON(r, &payload)` for request parsing
- Use `utils.WriteJSON(w, status, data)` for success responses
- Avoid generic error messages in auth endpoints (prevents user enumeration)

### Handlers Pattern

```go
type Handler struct {
    storage *pool.Database
}

func NewHandler(storage *pool.Database) *Handler {
    return &Handler{storage: storage}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
    router.HandleFunc("/resource", h.handleAction).Methods(http.MethodPost)
}
```

### Testing

- Use `testify/require` for assertions
- `TestMain` setup with `testQueries` global variable
- Use `context.Background()` in tests
- Create deterministic test data with helper functions
- Test database: `postgresql://root:password@localhost:5433/collab_board`

## Frontend Code Style Guidelines

### Imports & Structure

- External libraries first, then `@/` internal imports
- Components: Function components with TypeScript, PascalCase, default export
- File organization: `@/components/ui/` (primitives), `@/components/` (app), `@/screens/` (pages)

### Type Safety

- Use `z.infer<typeof schema>` for form types
- Define interfaces for complex API types
- Use proper TypeScript types throughout

### State & Data

- `useState` for local component state
- `useMutation` for mutations, `useQuery` for queries (TanStack Query)
- React Query config in `@/config/query-client.ts`

### Forms

- react-hook-form + zod + hookform resolvers pattern
- Use existing UI components from `@/components/ui/`

### Async Error Handling

- Try-catch blocks for async operations
- Don't throw on user actions, log errors in dev
- Don't use `as string` type assertions on errors

### Package References & Formatting

- Use `workspace:*` for internal package references
- Prettier: single quotes, tabs, trailing commas (es5)
- Tailwind classes via plugin

## Database & sqlc Guidelines

### Migrations

- Migrations in `db/migrations/` with `.up.sql` and `.down.sql`
- Run `make migrateup` to apply migrations
- Run `make migratedown` to rollback

### SQL Queries

- All queries in `db/queries/{entity}.sql`
- After modifying queries, regenerate: `make generate` (or manually run sqlc)
- Generated code in `db/sqlc/`

### sqlc Config

- Engine: PostgreSQL
- Uses pgx/v5 driver
- Generates Go code with types

## General Guidelines

- **NO COMMENTS** in code
- Follow existing patterns over introducing new conventions
- Reuse existing utilities: `utils/`, `internal/handlers/auth/`, `internal/middleware/`
- Use gorilla/mux for routing
- Use pgxpool for database connection pooling
- JWT-based authentication with access/refresh tokens
