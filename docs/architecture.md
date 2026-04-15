# Moana — layout

## Go packages

| Path | Role |
|------|------|
| `cmd/moana` | `main.go` → `run.go` (dispatch); `serve.go` HTTP; `cli_db.go` opens DB for CLI; `user.go` + `user_add.go` / `user_password.go` subcommands. |
| `internal/app` | Composition root: `New` builds `handlers.App` (`tmpl.Parse` → `render.Engine`); `HTTPHandler` wires `server.NewRouter` for production. |
| `internal/assets` | Embedded HTML templates and static files (`css/`, `js/`) consumed via `go:embed`. |
| `internal/server` | `router.go` + `static_health.go` (`/static/`, `/health`), `logging.go`, `graceful.go` (SIGINT/SIGTERM shutdown); delegates app routes to `handlers.RegisterRoutes`. |
| `internal/htmlview` | Template helpers: `funcs.go` / `format.go` / `format_money.go` / `display_identity.go` / `display_roles.go` / `icons.go`; `MergeFuncMaps` overlays handler-specific funcs. |
| `internal/tmpl` | `funcs.go` (`TemplateFuncMap`) + `parse.go` (`Parse`); `category_icons.go` holds category/dashboard icon helpers — keeps HTTP handlers free of `ParseFS` wiring. |
| `internal/render` | HTML execution: `layout_data.go`, `engine.go`, `shell.go` (`layout.html`), `simple.go` (login, etc.) — no routing or domain logic. |
| `internal/dashboard` | Dashboard domain: `page_data*.go` (`PageData`, `BuildPageData`, outflow + heatmap helpers), `metrics.go`, `heatmap.go`, `donut.go` — no `net/http`; handlers only render. |
| `internal/category` | Category color normalization, picker accents/icons/hints, legacy emoji→Lucide mapping (`NormalizeStoredIcon`), `BuildCategoriesList` (household-scoped categories from the store), form parsing — no HTML. |
| `internal/household` | Permission rules (manage / remove / leave) and `LoadSettingsPage` settings template payload. |
| `internal/historyview` | History page: `page_types.go` (`PageData`, `Nav`, `DayGroup`), `page.go` (`BuildPage`), `query_parse.go`, `nav.go` (`BuildNav`), `groups.go` (`GroupByDay`) — no `net/http`. |
| `internal/icons` | Embedded Lucide SVG registry (`data_gen.go`), `SVG` / `Inner` helpers — no HTTP. |
| `internal/handlers` | `App` in `app.go`; `layout.go`; `forms.go`; `params.go` (path id parsing); `RegisterRoutes` + `routes_*.go` (auth, dashboard, ledger, settings); `current_user.go` + `middleware_auth.go`; `login.go`; transactions split (`transaction_create*.go`, `transaction_form_types.go`, `transaction_edit*.go`); `categories.go` + `categories_render.go`; settings split (`settings.go`, `settings_profile.go`, `settings_redirect.go`, `settings_household.go`, `settings_members.go`). |
| `internal/store` | SQLite persistence: `user*.go`, `transaction*.go` (get/list/sql/row + types + mutate), `category*.go`, `household*.go` (types + query/mutate), `aggregate_sum.go` + `aggregate_category.go` + `aggregate_time.go`, `movement.go` (timestamps via `timeutil` sqlite helpers). |
| `internal/db` | SQLite open (`open.go` + `paths.go`), schema migrations (`migrate.go` + `migrate_steps.go`). |
| `internal/auth` | `argon2id.go` (PHC hashes), `password.go` (verify + bcrypt legacy), `session.go` + `session_sign.go` + `session_read.go` (HMAC cookies). |
| `internal/config` | `config.go` (`Load`, `Config`), `env.go` (getenv, `DBPath`). |
| `internal/tz` | Browser IANA zone from `moana_tz` cookie → `time.Location` (shared with `frontend/src/main.ts`). |
| `internal/dbutil` | `OpenStore(path)` = `db.Open` + `store.New` (caller closes `*sql.DB`). |
| `internal/money` | EUR parse/format, `AbsCents` — no I/O. |
| `internal/timeutil` | Time ranges: `calendar.go`, `trailing_days.go`, `range.go` (local date parsing); `sqlite.go` (RFC3339 TEXT round-trip for SQLite). |
| `internal/safepath` | Validates in-app `next` redirect targets (blocks open redirects). |
| `internal/txform` | `parsed.go` (`Parsed`) + `parse.go` (`Parse`) — form strings → cents, UTC time, category id. |
| `internal/testutil` | Integration tests: `app.go` (`DefaultTestConfig`, `NewApp`), `http.go` (`NewServer`, cookie client, `MustLogin`), `user.go` (`MustCreateUser`). |

## Frontend

- Source: `frontend/src` (Vite). Production bundles are written to `internal/assets/static` and embedded into the binary.

## Tests

- Integration tests that need the full router live in `internal/handlers` as **`package handlers_test`** to avoid an import cycle (`handlers` → `server` → `handlers`). They use **`internal/testutil`** for app + test server setup.
- Shared helper: `handlers_test.go` (`testApp`). HTTP smoke: `integration_server_test.go`. Transaction flows: `integration_transactions_test.go`.
