import { applyLocalTimeElements } from './lib/localTime'
import { setBrowserTimezoneCookie } from './lib/timezoneCookie'
import { initShellSidebar } from './lib/shellSidebar'
import { initSettingsMemberDialog } from './lib/settingsMemberDialog'
import { initCategoryModal } from './lib/categoryModal'
import { initHistoryControls } from './lib/historyControls'
import { initConfirmSubmitForms } from './lib/confirmSubmitForms'

/**
 * Ordered client initializers (single source of truth for `bootApp` order).
 * Each no-ops when its DOM is missing. Documented in `design.md`; see `boot.test.ts`.
 */
export const BOOT_APP_INITIALIZERS: ReadonlyArray<() => void> = [
  setBrowserTimezoneCookie,
  applyLocalTimeElements,
  initShellSidebar,
  initSettingsMemberDialog,
  initCategoryModal,
  initHistoryControls,
  initConfirmSubmitForms,
]

/**
 * Wire all client behaviors; each initializer no-ops when its DOM is missing.
 * Order: timezone cookie → local time labels → shell (listeners) → settings dialog → category modal → history controls → confirm-before-submit forms.
 * Cookie and time text run before interactive modules attach handlers.
 */
export function bootApp(): void {
  for (const init of BOOT_APP_INITIALIZERS) {
    init()
  }
}
