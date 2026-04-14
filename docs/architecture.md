# Moana — layout

## Go packages

| Path | Role |
|------|------|
| `cmd/moana` | `main.go` → `run.go` (dispatch); `serve.go` HTTP; `cli_db.go` opens DB for CLI; `user.go` + `user_add.go` / `user_password.go` subcommands. |
| `internal/app` | Composition root: `New` builds `handlers.App` (`tmpl.Parse` → `render.Engine`); `HTTPHandler` wires `server.NewRouter` for production. |
| `internal/assets` | Embedded HTML templates and static files (`css/`, `js/`) consumed via `go:embed`. |
| `internal/server` | `router.go` + `static_health.go` (`/static/`, `/health`), `logging.go`; delegates app routes to `handlers.RegisterRoutes`. |
| `internal/htmlview` | Template helpers: `funcs.go` / `format.go` / `display.go` + `icons.go`; `MergeFuncMaps` overlays handler-specific funcs. |
| `internal/tmpl` | Parses embedded `*.html` from `assets` and merges `htmlview` funcs with category/household template helpers — keeps HTTP handlers free of `ParseFS` wiring. |
| `internal/render` | `Engine` + `LayoutData`: executes the app shell (`layout.html`) and standalone pages — no routing or domain logic. |
| `internal/dashboard` | Dashboard domain: `page_data.go` (`PageData`, `BuildPageData`), `metrics.go`, `heatmap.go`, `donut.go` — no `net/http`; handlers only render. |
| `internal/category` | Category color normalization, picker accents/icons/hints, legacy emoji→Lucide mapping (`NormalizeStoredIcon`), `BuildCategoriesList` for the categories page payload, form parsing — no HTML. |
| `internal/household` | Permission rules (manage / remove / leave) and `LoadSettingsPage` settings template payload. |
| `internal/historyview` | History page: `page.go` (`BuildPage`), `nav.go` (`BuildNav`), `groups.go` (`GroupByDay`) — no `net/http`. |
| `internal/icons` | Embedded Lucide SVG registry (`data_gen.go`), `SVG` / `Inner` helpers — no HTTP. |
| `internal/handlers` | `App` in `app.go`; `layout.go` (shell render); `forms.go` (parse helpers); `RegisterRoutes` + `routes_*.go`; `auth.go`; transactions; `categories.go`; `settings*.go`. |
| `internal/store` | SQLite persistence: `user*.go`, `transaction*.go`, `category*.go` (types / query / mutate), `aggregate.go`, `movement.go`, `household.go`. |
| `internal/db` | SQLite open/migrate. |
| `internal/auth` | Password hashing and session cookies. |
| `internal/config` | Environment config. |
| `internal/tz` | Browser IANA zone from `moana_tz` cookie → `time.Location` (shared with `frontend/src/main.ts`). |
| `internal/dbutil` | `OpenStore(path)` = `db.Open` + `store.New` (caller closes `*sql.DB`). |
| `internal/money` | EUR parse/format, `AbsCents` — no I/O. |
| `internal/timeutil` | Calendar/month ranges and trailing local-day windows for stats. |
| `internal/safepath` | Validates in-app `next` redirect targets (blocks open redirects). |
| `internal/txform` | Parses transaction create/edit POST fields into cents + UTC time + category (`parse.go`). |
| `internal/testutil` | Integration tests: `app.go` (`DefaultTestConfig`, `NewApp`), `http.go` (`NewServer`, cookie client, `MustLogin`), `user.go` (`MustCreateUser`). |

## Frontend

- Source: `frontend/src` (Vite). Production bundles are written to `internal/assets/static` and embedded into the binary.

## Tests

- Integration tests that need the full router live in `internal/handlers` as **`package handlers_test`** to avoid an import cycle (`handlers` → `server` → `handlers`). They use **`internal/testutil`** for app + test server setup.
- Shared helper: `handlers_test.go` (`testApp`). HTTP smoke: `integration_server_test.go`. Transaction flows: `integration_transactions_test.go`.
