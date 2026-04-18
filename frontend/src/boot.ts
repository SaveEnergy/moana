import { applyLocalTimeElements } from './lib/localTime'
import { setBrowserTimezoneCookie } from './lib/timezoneCookie'
import { initShellSidebar } from './lib/shellSidebar'
import { initSettingsMemberDialog } from './lib/settingsMemberDialog'
import { initCategoryModal } from './lib/categoryModal'

/**
 * Wire all client behaviors; each initializer no-ops when its DOM is missing.
 * Order: timezone cookie → local time labels → shell (listeners) → settings dialog → category modal.
 * Cookie and time text run before interactive modules attach handlers.
 */
export function bootApp(): void {
  setBrowserTimezoneCookie()
  applyLocalTimeElements()
  initShellSidebar()
  initSettingsMemberDialog()
  initCategoryModal()
}
