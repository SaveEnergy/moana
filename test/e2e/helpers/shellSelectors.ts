/**
 * Playwright locator strings aligned with `frontend/src/lib/domSelectors.ts`
 * (`APP_TOPBAR_NOTIF_*`, `layout.html`). Keeps E2E hooks in one place when shell markup changes.
 */
export const NOTIFICATIONS_PATH = '/notifications' as const

/** Same intent as `APP_TOPBAR_NOTIF_LINK_SELECTOR` + explicit `href` (stable for routing). */
export const TOPBAR_NOTIFICATIONS_LINK = `a.app-topbar-notif-btn[href="${NOTIFICATIONS_PATH}"]` as const

/** Child of the bell link when the shell renders a positive unread count — `APP_NOTIF_BADGE_CLASS`. */
export const TOPBAR_NOTIFICATION_BADGE = '.app-notif-badge' as const
