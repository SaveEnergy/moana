/**
 * Category modal and similar UIs: radios are indexed by `value` (including "").
 * Avoids duplicate if/else + non-null assertions when selecting by dataset-driven values.
 */

/**
 * Select a radio from a map keyed by `HTMLInputElement.value`.
 * If `preferred` is not a key, selects `fallbackKey` (default `""` for auto/unset radios).
 * @returns whether a matching input existed and was checked.
 */
export function setRadioCheckedByValue(
  map: ReadonlyMap<string, HTMLInputElement>,
  preferred: string,
  fallbackKey = '',
): boolean {
  const key = map.has(preferred) ? preferred : fallbackKey
  const input = map.get(key)
  if (!input) {
    return false
  }
  input.checked = true
  return true
}
