/**
 * Unicode "Format" (Cf) — ZWSP, ZWJ, BOM-as-format, etc.
 * `String.prototype.trim()` does not remove these; strip before `trim()` when normalizing HTML / attribute text.
 */
export const UNICODE_FORMAT_CHARS = /\p{Cf}/gu

/** Non-global test — avoids `lastIndex` side effects from {@link UNICODE_FORMAT_CHARS}. */
const HAS_UNICODE_FORMAT_CHAR = /\p{Cf}/u

export function stripUnicodeFormatChars(s: string): string {
  if (!HAS_UNICODE_FORMAT_CHAR.test(s)) {
    return s
  }
  return s.replace(UNICODE_FORMAT_CHARS, '')
}
