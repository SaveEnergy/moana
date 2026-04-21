import { applyLocalTimeElements } from './lib/localTime'
import { setBrowserTimezoneCookie } from './lib/timezoneCookie'
import { initShellSidebar } from './lib/shellSidebar'
import { initSettingsMemberDialog } from './lib/settingsMemberDialog'
import { initCategoryModal } from './lib/categoryModal'
import { initHistoryControls } from './lib/historyControls'
import { initConfirmSubmitForms } from './lib/confirmSubmitForms'

/** One synchronous `bootApp` step (name must match {@link BOOT_INITIALIZER_NAMES}). */
export type BootInitializer = () => void

/**
 * Ordered client initializers (single source of truth for `bootApp` order).
 * Each no-ops when its DOM is missing. Documented in `design.md`; expected `function.name` order lives in `bootInitializerNames.ts` (asserted by `bootInitializers.test.ts`, which does not mock lib modules).
 */
export const BOOT_APP_INITIALIZERS: ReadonlyArray<BootInitializer> = [
  setBrowserTimezoneCookie,
  applyLocalTimeElements,
  initShellSidebar,
  initSettingsMemberDialog,
  initCategoryModal,
  initHistoryControls,
  initConfirmSubmitForms,
]
