/**
 * Leading or trailing Unicode whitespace (`\s` in Unicode mode).
 * Use to skip `String#trim` when the string is already tight (fewer allocations on hot paths).
 */
export const STRING_NEEDS_EDGE_TRIM = /^\s|\s$/u

/** `trim()` only when {@link STRING_NEEDS_EDGE_TRIM} matches. */
export function trimEdgesIfNeeded(s: string): string {
  return STRING_NEEDS_EDGE_TRIM.test(s) ? s.trim() : s
}
