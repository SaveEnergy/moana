/**
 * Playwright locator strings aligned with `frontend/src/lib/domSelectors.ts` and
 * `internal/assets/templates/` (`layout.html`, `settings.html`, `history.html`).
 * Single place to update when template ids change.
 */

// --- Notifications (`layout.html` topbar; `APP_TOPBAR_NOTIF_*`, `NOTIFICATIONS_PATH`)
export const NOTIFICATIONS_PATH = '/notifications' as const

/** Same intent as `APP_TOPBAR_NOTIF_LINK_SELECTOR` + explicit `href`. */
export const TOPBAR_NOTIFICATIONS_LINK = `a.app-topbar-notif-btn[href="${NOTIFICATIONS_PATH}"]` as const

/** `APP_NOTIF_BADGE_CLASS` — child of the bell link when unread count is positive. */
export const TOPBAR_NOTIFICATION_BADGE = '.app-notif-badge' as const

// --- Mobile shell (`layout.html`; mirror `APP_SHELL_SELECTOR`, etc.)
export const APP_SHELL = '#app-shell' as const
export const APP_SIDEBAR_BACKDROP = '#app-sidebar-backdrop' as const
export const APP_SIDEBAR_CLOSE = '#app-sidebar-close' as const
export const APP_SIDEBAR_NAV = '#app-sidebar-nav' as const
export const APP_GLOBAL_SEARCH = '#app-global-search' as const

// --- Settings add-member (`settings.html`; `SETTINGS_ADD_MEMBER_*_ELEMENT_ID`)
export const SETTINGS_ADD_MEMBER_DIALOG_ID = 'settings-add-member-dialog' as const
export const SETTINGS_ADD_MEMBER_DIALOG = `#${SETTINGS_ADD_MEMBER_DIALOG_ID}` as const
export const SETTINGS_ADD_MEMBER_OPEN = '#settings-add-member-open' as const
export const SETTINGS_ADD_MEMBER_TITLE = '#settings-add-member-title' as const
export const SETTINGS_ADD_MEMBER_CANCEL = '#settings-add-member-cancel' as const

/** Readonly profile email (`settings.html`). */
export const SETTINGS_EMAIL = '#settings-email' as const

// --- History (`history.html`)
export const HISTORY_Q = '#history-q' as const
export const HISTORY_FROM = '#history-from' as const
export const HISTORY_TO = '#history-to' as const
/** `HISTORY_SORT_ELEMENT_ID` */
export const HISTORY_SORT = '#history-sort' as const

// --- Category modal (`categories.html`; `CATEGORY_MODAL_*_ELEMENT_ID`)
export const CATEGORY_MODAL_ID = 'cat-modal' as const
export const CATEGORY_MODAL = `#${CATEGORY_MODAL_ID}` as const
export const CATEGORY_MODAL_OPEN_CREATE = '#cat-modal-open-create' as const
export const CATEGORY_MODAL_TITLE = '#cat-modal-title' as const
export const CATEGORY_MODAL_CLOSE = '#cat-modal-close' as const
export const CATEGORY_MODAL_NAME = '#cat-modal-name' as const
export const CATEGORY_MODAL_SUBMIT = '#cat-modal-submit' as const
