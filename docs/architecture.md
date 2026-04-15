# Moana — layout

## Go packages

| Path | Role |
|------|------|
| `cmd/moana` | `main.go` → `run.go` (dispatch); `serve.go` HTTP; `cli_db.go` opens DB for CLI; `user.go` + `user_add.go` / `user_password.go` subcommands. |
| `internal/app` | Composition root: `New` builds `handlers.App` (`tmpl.Parse` → `render.Engine`); `HTTPHandler` wires `server.NewRouter` for production. |
| `internal/assets` | Embedded HTML templates and static files (`css/`, `js/`) consumed via `go:embed`. |
| `internal/server` | `router.go` (wraps mux with `WithRequestTimeout` from `Config.RequestTimeout`); `static.go` + `health.go` + `static_health.go` (`/static/`, `/health`); `logging.go`; `graceful.go` (SIGINT/SIGTERM); delegates app routes to `handlers.RegisterRoutes`. |
| `internal/htmlview` | Template helpers: `funcs.go` / `format.go` / `format_money.go` / `display_identity.go` / `display_roles.go` / `icons.go`; `MergeFuncMaps` overlays handler-specific funcs. |
| `internal/httperr` | Generic HTTP 500 responses (`InternalMessage` + `slog`); avoids leaking DB/template errors to clients — used by `handlers` and `render`. |
| `internal/tmpl` | `funcs.go` (`TemplateFuncMap`) + `parse.go` (`Parse`); `category_icons.go` holds category/dashboard icon helpers — keeps HTTP handlers free of `ParseFS` wiring. |
| `internal/render` | HTML execution: `layout_data.go`, `engine.go`, `shell.go` (`layout.html`), `simple.go` (login, etc.) — no routing or domain logic. |
| `internal/dashboard` | Dashboard domain: `page_data*.go` (`PageData`, `BuildPageData`, `?period=` parsing in `page_data_period.go`, outflow + heatmap helpers), `metrics.go`, `heatmap.go`, `donut.go` — no `net/http`; handlers only render. |
| `internal/category` | Category color normalization, picker accents/icons/hints, legacy emoji→Lucide mapping (`NormalizeStoredIcon`), `BuildCategoriesList` (household-scoped categories from the store), form parsing — no HTML. |
| `internal/household` | `permissions_roles.go` (manage / leave) + `permissions_remove.go` (remove); `settings_page.go` (`LoadSettingsPage`). |
| `internal/historyview` | History page: `page_types.go` (`PageData`, `Nav`, `DayGroup`), `page.go` (`BuildPage`), `limit.go` (default row cap + over-fetch to detect truncation, `trimHistoryRows`), `query_parse.go` (`ParseHistoryURL`), `nav.go` (`BuildNav`), `groups.go` (`GroupByDay`) — no `net/http`. |
| `internal/icons` | Embedded Lucide SVG registry (`data_gen.go`), `SVG` / `Inner` helpers — no HTTP. |
| `internal/handlers` | `App` in `app.go`; `layout.go`; `forms.go`; `params.go` (path id parsing); `RegisterRoutes` + `routes_*.go` (auth, dashboard, ledger, settings); `current_user.go` + `middleware_auth.go`; `login.go`; transactions split (`transaction_create*.go`, `transaction_form_types.go`, `transaction_edit*.go`); `categories.go` + `categories_render.go`; settings split (`settings.go`, `settings_profile.go`, `settings_redirect.go`, `settings_household.go`, `settings_members.go`). |
| `internal/store` | SQLite persistence: `user*.go`, `transaction*.go` (get/list/sql/row + types + mutate), `category*.go`, `household*.go` (types + query/mutate), `aggregate_sum.go` (`SumIncomeExpenseCentsInRange` + per-kind sums) + `aggregate_category.go` + `aggregate_time.go`, `movement.go` (timestamps via `timeutil` sqlite helpers). Tests split into `store_*_test.go` plus `store_test_helpers_test.go`. |
| `internal/db` | SQLite open (`open.go` + `paths.go`); `migrate.go` runs ordered steps (`migrationSteps`: v1–v5 in `migrate_steps.go`, v6 in `migrate_v6.go`, v7 `idx_users_household_id` in `migrate_v7.go`). |
| `internal/auth` | `argon2id.go` (PHC hashes), `password.go` (verify + bcrypt legacy), `session.go` + `session_sign.go` + `session_read.go` (HMAC cookies). |
| `internal/config` | `config.go` (`Load`, `Config`), `env.go` (getenv, `DBPath`), `doc.go` (package overview); `config_test.go` regression-tests `MOANA_*` wiring. |
| `internal/tz` | Browser IANA zone from `moana_tz` cookie → `time.Location` (shared with `frontend/src/main.ts`). |
| `internal/dbutil` | `OpenStore(path)` = `db.Open` + `store.New` (caller closes `*sql.DB`); `MustOpenMemStore` for unit tests outside `store`. |
| `internal/money` | EUR parse/format, `AbsCents` — no I/O. |
| `internal/timeutil` | Time ranges: `calendar.go`, `trailing_days.go`, `range.go` (local date parsing); `sqlite.go` (RFC3339 TEXT round-trip for SQLite). |
| `internal/safepath` | Validates in-app `next` redirect targets (blocks open redirects). |
| `internal/txform` | `parsed.go` (`Parsed`) + `parse.go` (`Parse`) — form strings → cents, UTC time, category id. |
| `internal/testutil` | Integration tests: `app.go` (`DefaultTestConfig` includes `RepoURL` / `DefaultTestRepoURL` for footer parity, `NewApp`), `http.go` (`NewAppServer`, `NewServer`, cookie client, `MustLogin`), `user.go` (`MustCreateUser`). |

## Frontend

- Source: `frontend/src` (Vite). Production bundles are written to `internal/assets/static` and embedded into the binary.

## Tests

- Integration tests that need the full router live in `internal/handlers` as **`package handlers_test`** to avoid an import cycle (`handlers` → `server` → `handlers`). They use **`internal/testutil`** (`NewAppServer` / `NewApp`, `MustLogin`, etc.).
- Files: `integration_server_test.go` (login, health, static `/static/`, unauthenticated `/` → `/login`, dashboard period, smoke pages, history filters), `integration_login_test.go`, `integration_logout_test.go`, `integration_transactions_test.go` (create/edit validation, edit 404), `integration_categories_test.go` (create/delete/duplicate, update rename, invalid update id), `integration_settings_test.go` (profile, password change + mismatch + wrong current, household name, member add + duplicate email, remove flows, owner blocked from leave when others remain, sole owner leave, invalid remove id, member cannot remove owner).
- Unit: `internal/config/config_test.go` exercises `config.Load` and `MOANA_*` env behavior (listen, DB path, sessions, timeouts, secure cookies, repo URL).
