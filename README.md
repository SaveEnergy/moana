# Moana

Personal finance web app: transactions, categories, budgets, and dashboards backed by SQLite. Server-rendered HTML with embedded templates; static assets are built with Vite and embedded into the Go binary.

## Stack

- **Go** (see `go.mod` for the toolchain version)
- **SQLite** via `modernc.org/sqlite` (no CGO in release builds)
- **Frontend:** Bun + Vite (`frontend/`) → output copied to `internal/assets/static/`

## Prerequisites

- Go **1.26+**
- [Bun](https://bun.sh/) (install + `bun run` for `package.json` scripts — same commands as `Makefile`; npm is not required for the JS toolchain)

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
4. Create the first user (same database path as the server; uses `MOANA_DB_PATH` if set):
  ```bash
   ./bin/moana user add -email you@example.com -password 'your-password'
  ```
   Optional: `-role admin` (default is `user`). With Docker: `docker compose exec moana /bin/moana user add -email you@example.com -password '…'`.

For hot reload during development, see `make dev` and `scripts/dev.sh`.

## Docker

The image is defined in `Dockerfile`. It sets `MOANA_ENV=production`, stores the database under `/data/moana.db`, and listens on `:8080`.

From the repository root:

```bash
cp .env.example .env
# set MOANA_SESSION_SECRET (and MOANA_PUBLIC_BASE_URL / SendGrid if you use mail)
docker compose up --build
```

Then open `http://localhost:8080`. The SQLite file persists in the named volume `moana-data`. See `.env.example` and `docker-compose.yml` for all supported variables.

`MOANA_SESSION_SECRET` is **required** when `MOANA_ENV=production` (see [Configuration](#configuration)).

**Build notes (VPS/CI):** The first `docker compose up --build` can sit on the **assets** step (`bun run build`) or **apk** for several minutes while layers download and Vite compiles. That is expected on a small machine or a slow path to Docker Hub. If the **Go** compile step is killed, the host is likely OOM: add swap, build on a machine with more RAM, or run `go build` with less parallelism (for example `GOMAXPROCS=1` when invoking `go build` in a custom image). Rebuilds are faster thanks to layer and module caches when BuildKit is enabled.

## Configuration


| Variable                                    | Default                               | Notes                                                                                                                                                                                                                     |
| ------------------------------------------- | ------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `MOANA_LISTEN`                              | `:8080`                               | HTTP listen address                                                                                                                                                                                                       |
| `MOANA_DB_PATH`                             | `data/moana.db`                       | SQLite file path                                                                                                                                                                                                          |
| `MOANA_ENV`                                 | `development`                         | `production` enables secure cookies and **requires** `MOANA_SESSION_SECRET`                                                                                                                                               |
| `MOANA_SESSION_SECRET`                      | *(dev fallback)*                      | Must be set in production                                                                                                                                                                                                 |
| `MOANA_SESSION_MAX_AGE_SEC`                 | `604800`                              | Session cookie max age (seconds)                                                                                                                                                                                          |
| `MOANA_REQUEST_TIMEOUT_SEC`                 | `60`                                  | Per-request timeout                                                                                                                                                                                                       |
| `MOANA_REPO_URL`                            | `https://github.com/SaveEnergy/moana` | Shown in the app footer                                                                                                                                                                                                   |
| `MOANA_PUBLIC_BASE_URL`                     | —                                     | Public site URL (`https://…`), **required** with `MOANA_SENDGRID_API_KEY` so password reset links are absolute                                                                                                            |
| `MOANA_PASSWORD_RESET_TTL_MIN`              | `60`                                  | How long a reset link stays valid (minutes)                                                                                                                                                                               |
| `MOANA_SENDGRID_API_KEY`                    | —                                     | [SendGrid](https://sendgrid.com/) API key; with `MOANA_MAIL_FROM`, `MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID`, and `MOANA_PUBLIC_BASE_URL` enables “forgot password” (dynamic template API)                              |
| `MOANA_MAIL_FROM`                           | —                                     | Sender address, verified in SendGrid; **required** if `MOANA_SENDGRID_API_KEY` is set                                                                                                                                     |
| `MOANA_SENDGRID_PASSWORD_RESET_TEMPLATE_ID` | —                                     | Dynamic template id (`d-…`); design must use Handlebars `{{reset_url}}` (see `internal/mail/template_password_reset.html` as a starting point)                                                                            |
| `MOANA_TRUST_X_FORWARDED_FOR`               | *(off)*                               | `1`, `true`, or `yes` (case-insensitive): use the **first** `X-Forwarded-For` value as the client IP for **rate limits**. Enable only when Moana sits behind a **trusted** reverse proxy that sets or strips this header. |
| `MOANA_RATE_LIMIT_LOGIN_PER_MIN`            | `20`                                  | Rolling per-minute cap on `POST` **login** per client IP; `0` disables. On limit, responses use **429** and `Retry-After: 60`.                                                                                            |
| `MOANA_RATE_LIMIT_FORGOT_PASSWORD_PER_MIN`  | `10`                                  | Same for `POST` **forgot password**; `0` disables.                                                                                                                                                                        |


**HTTP hardening (no extra env):** all responses get baseline security headers (`X-Content-Type-Options: nosniff`, `X-Frame-Options: DENY`, `Referrer-Policy`, and a tight `Permissions-Policy`). A full Content-Security-Policy is not set (templates use inline styles).

**SendGrid: API counted but no inbox?** A `202` from `/v3/mail/send` only means the message was *accepted*; check **Email API → Email activity** (or Activity) for the recipient. The server logs a successful queue as `x_message_id=…` in `sendgrid` — search that id in activity. If *Delivered* but nothing in the mailbox, check **spam** and **Domain Authentication** + a verified **From** address. If *Dropped* or *Blocked*, check **Suppression** for that address. The dynamic template must be **active** and include `{{reset_url}}` (same spelling as the API field).

**Activity shows only “Processed” (no “Delivered”), “Sending IP: N/A”, end-to-end time N/A?** That usually means the message is **stuck in SendGrid’s queue** for an **account** reason (compliance / review, new-account limits, **dormant** account, **billing** freeze, or **dedicated IP** not enabled after reactivation). Moana has already done its job. Fix it in the Twilio/SendGrid console: look for a **red banner** (“Under review” etc.), [Settings → IP Addresses](https://app.sendgrid.com/settings/ip_addresses) and ensure a sending IP is present and “Allow my account to send mail using this IP” is on, and read [Emails stuck in Processing status](https://support.sendgrid.com/hc/en-us/articles/31236228514587-Emails-stuck-in-Processing-status) — often you must reply to **Compliance** or open a **support ticket**; mail can sit up to **72 hours** then bounce if unresolved.

## Tests

```bash
make test        # unit tests (builds frontend first)
make test-e2e    # Playwright E2E (requires frontend build)
```

## Docs

- [Architecture](docs/architecture.md)

## License

[MIT](LICENSE.md)