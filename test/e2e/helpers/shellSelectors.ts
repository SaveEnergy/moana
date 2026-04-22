/**
 * Playwright locator strings aligned with `frontend/src/lib/domSelectors.ts` and
 * `internal/assets/templates/` (`login.html`, `layout.html`, `settings.html`, `history.html`, transaction forms).
 * Single place to update when template ids change.
 */

// --- Login (`login.html`; `login-email`, `login-password`, `login-forgot`)
export const LOGIN_INPUT_EMAIL = '#login-email' as const
export const LOGIN_INPUT_PASSWORD = '#login-password' as const
/** Forgot password link on the sign-in page (`a.login-forgot`). */
export const LOGIN_BUTTON_FORGOT = 'a.login-forgot' as const

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

// --- Topbar user menu (`layout.html`; open state uses `APP_USER_MENU_OPEN_SELECTOR` in TS)
export const SETTINGS_PATH = '/settings' as const
export const APP_USER_MENU = 'details.app-user-menu' as const
export const APP_USER_MENU_SUMMARY = 'details.app-user-menu summary.app-user-menu-btn' as const
export const USER_MENU_SETTINGS_LINK = `a.app-user-menu-settings[href="${SETTINGS_PATH}"]` as const

// --- Settings add-member (`settings.html`; `SETTINGS_ADD_MEMBER_*_ELEMENT_ID`)
export const SETTINGS_ADD_MEMBER_DIALOG_ID = 'settings-add-member-dialog' as const
export const SETTINGS_ADD_MEMBER_DIALOG = `#${SETTINGS_ADD_MEMBER_DIALOG_ID}` as const
export const SETTINGS_ADD_MEMBER_OPEN = '#settings-add-member-open' as const
export const SETTINGS_ADD_MEMBER_TITLE = '#settings-add-member-title' as const
export const SETTINGS_ADD_MEMBER_CANCEL = '#settings-add-member-cancel' as const

/** Dismiss hit target for add-member backdrop tests (`settings.html`). */
export const SETTINGS_ADD_MEMBER_DIALOG_HEADER = '.admin-add-dialog-header' as const

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

export const CATEGORY_MODAL_PANEL = '.cat-modal-panel' as const
export const CATEGORY_LIST_ROW = '.cat-list-row' as const
export const CATEGORY_DELETE_FORM = 'form.cat-delete' as const

// --- Transaction forms (shared `name` attrs on `transactions_new.html` / `transactions_edit.html`)
export const TX_INPUT_AMOUNT = 'input[name="amount"]' as const
export const TX_INPUT_OCCURRED_ON = 'input[name="occurred_on"]' as const

/** Kind pill + hidden radio (`kind-toggle` on transaction forms). */
export const TX_KIND_LABEL_EXPENSE = 'label.kind-option:has(input[name="kind"][value="expense"])' as const
export const TX_KIND_LABEL_INCOME = 'label.kind-option:has(input[name="kind"][value="income"])' as const
export const TX_INPUT_KIND_EXPENSE = 'input[name="kind"][value="expense"]' as const
export const TX_INPUT_KIND_INCOME = 'input[name="kind"][value="income"]' as const

// --- Transaction edit (`transactions_edit.html`)
export const TX_EDIT_NOTE = '#tx-edit-note' as const

/**
 * Row stamps after save / history — prefix of **`LOCAL_TIME_ELEMENTS_SELECTOR`** (`time.js-local-time[datetime]`).
 */
export const LOCAL_TIME_DISPLAY = 'time.js-local-time' as const
