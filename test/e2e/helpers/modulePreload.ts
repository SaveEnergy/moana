/**
 * Canonical URLs for the Vite-built shell bundle (see `frontend/vite.config.ts` →
 * `internal/assets/static/`). **Must** stay aligned with `login.html` / `layout.html`
 * `<link href>` values and `design.md` §2 / §13.
 */
export const STATIC_APP_CSS_PATH = '/static/css/app.css'
export const STATIC_APP_JS_PATH = '/static/js/app.js'

/**
 * Playwright locator for `<link rel="modulepreload" href="…app.js">` on shell pages.
 */
export const APP_JS_MODULE_PRELOAD = `link[rel="modulepreload"][href="${STATIC_APP_JS_PATH}"]`
