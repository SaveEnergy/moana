/**
 * Playwright locator string for the shell’s early `app.js` fetch
 * (`<link rel="modulepreload" href="/static/js/app.js">` on `login.html` / `layout.html`).
 * See `design.md` §2 / §14.
 */
export const APP_JS_MODULE_PRELOAD = 'link[rel="modulepreload"][href="/static/js/app.js"]'
