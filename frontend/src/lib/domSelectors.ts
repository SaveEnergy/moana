/**
 * Shared selector strings for DOM that is also authored in HTML templates (`internal/assets/templates/`).
 * Centralizes class hooks so refactors don’t miss a string copy in `dialogKeyboard` / related logic.
 */

/** Topbar account menu when expanded (`layout.html` `details.app-user-menu[open]`). */
export const APP_USER_MENU_OPEN_SELECTOR = 'details.app-user-menu[open]' as const
