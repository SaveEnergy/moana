import { type ClickTargetEvent, clickEventTargetElement } from './clickTarget'
import { APP_SIDEBAR_CLOSE_SELECTOR } from './domSelectors'

/**
 * Whether a bubbled `click` on `#app-shell` should close the mobile drawer.
 * Caller must still gate on the mobile `matchMedia` breakpoint. Uses {@link ClickTargetEvent} like `dialogDismiss`.
 */
export function shouldCloseMobileSidebarFromShellClick(
  e: ClickTargetEvent,
  backdrop: Element | null,
): boolean {
  const el = clickEventTargetElement(e)
  if (!el) {
    return false
  }
  if (backdrop !== null && el === backdrop) {
    return true
  }
  return el.closest(APP_SIDEBAR_CLOSE_SELECTOR) !== null
}
