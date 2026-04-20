/**
 * Expected `function.name` for each entry in {@link BOOT_APP_INITIALIZERS} (`boot.ts`), in order.
 * Keep in sync when changing initializer order; `bootInitializers.test.ts` asserts against `BOOT_APP_INITIALIZERS`.
 */
export const BOOT_INITIALIZER_NAMES = [
  'setBrowserTimezoneCookie',
  'applyLocalTimeElements',
  'initShellSidebar',
  'initSettingsMemberDialog',
  'initCategoryModal',
  'initHistoryControls',
  'initConfirmSubmitForms',
] as const

/** Length of {@link BOOT_INITIALIZER_NAMES} — use in tests that mock initializer implementations. */
export const BOOT_INITIALIZER_COUNT = BOOT_INITIALIZER_NAMES.length
