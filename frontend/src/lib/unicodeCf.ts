/**
 * Unicode "Format" (Cf) — ZWSP, ZWJ, BOM-as-format, etc.
 * `String.prototype.trim()` does not remove these; strip before `trim()` when normalizing HTML / attribute text.
 */
export const UNICODE_FORMAT_CHARS = /\p{Cf}/gu

export function stripUnicodeFormatChars(s: string): string {
  return s.replace(UNICODE_FORMAT_CHARS, '')
}
