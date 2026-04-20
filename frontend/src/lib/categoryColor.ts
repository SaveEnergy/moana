import { CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM } from './domSelectors'

/** Default preview strip when color = Auto; aligns with server-rendered modal gradient. */
export const CATEGORY_MODAL_DEFAULT_PREVIEW_BG =
  'color-mix(in srgb, var(--primary) 12%, #fff8f0)' as const

/**
 * Default for native `<input type="color">` and invalid `data-custom-hex` — matches
 * `categories.html` initial `value` and `sanitizeCategoryCustomHex` fallback.
 */
export const CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT = '#818cf8' as const

/**
 * CSS `background` for the category modal preview strip from the checked color `value` and native hex.
 * Pure helper for `categoryModal.ts` (`syncCatModalPreview`).
 */
export function resolveCategoryModalPreviewBackground(
  checkedColorValue: string | undefined,
  nativeColorInputValue: string | undefined,
): string {
  const checked = checkedColorValue?.trim()
  if (checked === CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM) {
    return nativeColorInputValue?.trim() || CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT
  }
  if (checked) {
    return checked
  }
  return CATEGORY_MODAL_DEFAULT_PREVIEW_BG
}

/**
 * Whether {@link resolveCategoryModalPreviewBackground} produced a different string than last paint —
 * `categoryModal.ts` uses this to skip redundant `style.background` writes on the preview strip.
 */
export function shouldUpdateCategoryModalPreviewBackground(
  lastResolved: string | undefined,
  nextResolved: string,
): boolean {
  return lastResolved !== nextResolved
}

const VALID_CATEGORY_CUSTOM_HEX = /^#[0-9a-fA-F]{6}$/

/** Sanitize custom hex from data attributes for the category modal native color input. */
export function sanitizeCategoryCustomHex(
  hex: string,
  fallback: string = CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT,
): string {
  if (VALID_CATEGORY_CUSTOM_HEX.test(hex)) {
    return hex
  }
  const t = hex.trim()
  return VALID_CATEGORY_CUSTOM_HEX.test(t) ? t : fallback
}
