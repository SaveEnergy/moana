/**
 * Expected `function.name` for each entry in {@link BOOT_APP_INITIALIZERS} (`boot.ts`), in order.
 * Update when changing initializer order or adding a step; `bootInitializers.test.ts` asserts against this list.
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
