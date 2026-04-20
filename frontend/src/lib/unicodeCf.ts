import { trimEdgesIfNeeded } from './trimEdges'

/**
 * Unicode "Format" (Cf) — ZWSP, ZWJ, BOM-as-format, etc.
 * Strip **Cf** before edge normalization — `String.prototype.trim()` does not remove format chars.
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

/** **`stripUnicodeFormatChars`** then **`trimEdgesIfNeeded`** — shared by `data-*` and `data-confirm` normalization. */
export function stripCfTrimEdges(s: string): string {
  return trimEdgesIfNeeded(stripUnicodeFormatChars(s))
}
