# Moana

Personal finance web app: transactions, categories, budgets, and dashboards backed by SQLite. Server-rendered HTML with embedded templates; static assets are built with Vite and embedded into the Go binary.

## Stack

- **Go** (see `go.mod` for the toolchain version)
- **SQLite** via `modernc.org/sqlite` (no CGO in release builds)
- **Frontend:** Bun + Vite (`frontend/`) → output copied to `internal/assets/static/`

## Prerequisites

- Go **1.26+**
- [Bun](https://bun.sh/) (used by `Makefile` for the frontend; see `frontend/package.json` if you prefer another package manager)

## Build and run (local)

1. Build the frontend (required for `go:embed`):

   ```bash
   make build-frontend
   ```

2. Build the binary:

   ```bash
   make build
   ```

   Output: `bin/moana`.

3. Run the server (defaults: listen `:8080`, database `data/moana.db`):

   ```bash
   ./bin/moana
   ```

   Equivalent: `./bin/moana serve`. With no subcommand, `serve` is implied.

4. Create the first user (same database path as the server):

   ```bash
   ./bin/moana user add you@example.com 'your-password'
   ```

For hot reload during development, see `make dev` and `scripts/dev.sh`.

## Docker

The image is defined in `Dockerfile`. It sets `MOANA_ENV=production`, stores the database under `/data/moana.db`, and listens on `:8080`.

From the repository root:

```bash
export MOANA_SESSION_SECRET='a-long-random-secret'
docker compose up --build
```

Then open `http://localhost:8080`. The SQLite file persists in the named volume `moana-data`.

`MOANA_SESSION_SECRET` is **required** when `MOANA_ENV=production` (see [Configuration](#configuration)).

## Configuration

| Variable | Default | Notes |
|----------|---------|--------|
| `MOANA_LISTEN` | `:8080` | HTTP listen address |
| `MOANA_DB_PATH` | `data/moana.db` | SQLite file path |
| `MOANA_ENV` | `development` | `production` enables secure cookies and **requires** `MOANA_SESSION_SECRET` |
| `MOANA_SESSION_SECRET` | *(dev fallback)* | Must be set in production |
| `MOANA_SESSION_MAX_AGE_SEC` | `604800` | Session cookie max age (seconds) |
| `MOANA_REQUEST_TIMEOUT_SEC` | `60` | Per-request timeout |
| `MOANA_REPO_URL` | `https://github.com/SaveEnergy/moana` | Shown in the app footer |

## Tests

```bash
make test        # unit tests (builds frontend first)
make test-e2e    # Playwright E2E (requires frontend build)
```

## Docs

- [Architecture](docs/architecture.md)

## License

[MIT](LICENSE.md)
