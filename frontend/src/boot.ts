import { applyLocalTimeElements } from './lib/localTime'
import { setBrowserTimezoneCookie } from './lib/timezoneCookie'
import { initShellSidebar } from './lib/shellSidebar'
import { initSettingsMemberDialog } from './lib/settingsMemberDialog'
import { initCategoryModal } from './lib/categoryModal'

/** Wire all client behaviors; each initializer no-ops when its DOM is missing. */
export function bootApp(): void {
  setBrowserTimezoneCookie()
  applyLocalTimeElements()
  initShellSidebar()
  initSettingsMemberDialog()
  initCategoryModal()
}
