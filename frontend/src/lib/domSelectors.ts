/**
 * Shared selector strings for DOM that is also authored in HTML templates (`internal/assets/templates/`).
 * Centralizes hooks so refactors don’t miss a string copy across keyboard, shell, and dismiss logic.
 */

/** Topbar account menu when expanded (`layout.html` `details.app-user-menu[open]`). */
export const APP_USER_MENU_OPEN_SELECTOR = 'details.app-user-menu[open]' as const

/** Mobile shell root (`layout.html` `id="app-shell"`). */
export const APP_SHELL_ELEMENT_ID = 'app-shell' as const

/** `aria-controls` target; mobile menu affordance (`layout.html`). */
export const APP_SIDEBAR_TOGGLE_ELEMENT_ID = 'app-sidebar-toggle' as const

/** Dimmed overlay behind the off-canvas drawer (`layout.html`). */
export const APP_SIDEBAR_BACKDROP_ELEMENT_ID = 'app-sidebar-backdrop' as const

/** In-drawer close control (`layout.html`); use with `Element.closest()` / `#…` form below. */
export const APP_SIDEBAR_CLOSE_ELEMENT_ID = 'app-sidebar-close' as const

/** `#` + {@link APP_SIDEBAR_CLOSE_ELEMENT_ID} — `mobileShellDismiss` uses `closest()`, not `getElementById`. */
export const APP_SIDEBAR_CLOSE_SELECTOR = `#${APP_SIDEBAR_CLOSE_ELEMENT_ID}` as const
