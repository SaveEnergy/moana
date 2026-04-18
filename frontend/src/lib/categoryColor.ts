/** Default preview strip when color = Auto; aligns with server-rendered modal gradient. */
export const CATEGORY_MODAL_DEFAULT_PREVIEW_BG =
  'color-mix(in srgb, var(--primary) 12%, #fff8f0)' as const

/** Sanitize custom hex from data attributes for the category modal native color input. */
export function sanitizeCategoryCustomHex(hex: string, fallback = '#818cf8'): string {
  const t = hex.trim()
  return /^#[0-9a-fA-F]{6}$/.test(t) ? t : fallback
}
