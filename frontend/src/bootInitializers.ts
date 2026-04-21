import { applyLocalTimeElements } from './lib/localTime'
import { setBrowserTimezoneCookie } from './lib/timezoneCookie'
import { initShellSidebar } from './lib/shellSidebar'
import { initSettingsMemberDialog } from './lib/settingsMemberDialog'
import { initCategoryModal } from './lib/categoryModal'
import { initHistoryControls } from './lib/historyControls'
import { initConfirmSubmitForms } from './lib/confirmSubmitForms'

/** One synchronous `bootApp` step (`function.name` must match the same index in {@link BOOT_INITIALIZER_NAMES}). */
export type BootInitializer = () => void

/**
 * Ordered client initializers (single source of truth for `bootApp` order).
 * Each no-ops when its DOM is missing. Documented in `design.md`.
 * {@link BOOT_INITIALIZER_NAMES} is derived here so reordering this array cannot drift from a parallel name table.
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

/**
 * `function.name` per {@link BOOT_APP_INITIALIZERS} (Vitest / tooling; unused by `bootApp` and typically tree-shaken from production).
 */
export const BOOT_INITIALIZER_NAMES: readonly string[] = Object.freeze(
  BOOT_APP_INITIALIZERS.map((fn) => fn.name),
)

/** Length of {@link BOOT_APP_INITIALIZERS} — use in tests that mock initializer implementations. */
export const BOOT_INITIALIZER_COUNT = BOOT_APP_INITIALIZERS.length
