/**
 * Category modal and similar UIs: radios are indexed by `value` (including "").
 * Avoids duplicate if/else + non-null assertions when selecting by dataset-driven values.
 */

/**
 * Build a map of `HTMLInputElement.value` → input for radios under `scope` matching `radiosSelector`.
 * Later duplicates overwrite earlier entries (same as `querySelectorAll` iteration order).
 */
export function buildRadioMapByValue(
  scope: ParentNode,
  radiosSelector: string,
): Map<string, HTMLInputElement> {
  const m = new Map<string, HTMLInputElement>()
  for (const r of scope.querySelectorAll<HTMLInputElement>(radiosSelector)) {
    m.set(r.value, r)
  }
  return m
}

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
