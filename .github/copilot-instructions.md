# Repository Guidelines

## Project Overview

- This is a small Go HTTP API module named `practice`.
- The application starts in `main.go` and listens on `localhost:8080`.
- HTTP handlers are implemented in the root package (`main`); shared request and response models are in `models/`.
- SQLite access uses `database/sql` with the `modernc.org/sqlite` driver.

## Structure and Conventions

- Keep route registration and application startup in `main.go`.
- Keep HTTP response helpers and route handlers in the existing root-level handler files.
- Keep database initialization, schema creation, and JSON seeding in `sqlite_setup.go`.
- Keep user query functions in `sqlite_handler.go`.
- Use the existing `models.User`, `models.Address`, and authentication DTOs rather than duplicating types.
- Preserve the current Go standard-library HTTP patterns, including method/path patterns such as `GET /users/{id}`.
- Use wrapped errors at database boundaries and handle `sql.ErrNoRows` at the HTTP layer.

## Data and Runtime

- The application expects `data/users.json` relative to the process working directory when the database is empty.
- `users.db` is a local runtime artifact and should not be treated as source code.
- SQLite is configured for a single open and idle connection; retain that constraint unless concurrency behavior is intentionally redesigned.

## Validation

- Format Go changes with `gofmt`.
- Run `go test ./...` from the repository root after changes.
- For endpoint changes, also exercise the affected route with the local server and verify status codes and JSON responses.