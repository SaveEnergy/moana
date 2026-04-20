/**
 * Category modal and similar UIs: radios are indexed by `value` (including "").
 * Avoids duplicate if/else + non-null assertions when selecting by dataset-driven values.
 */

/**
 * Build a map of trimmed `HTMLInputElement.value` → input for radios under `scope` matching `radiosSelector`.
 * Keys use **`value.trim()`** so {@link getFormRadioGroupValue} and {@link setRadioCheckedByValue} stay consistent.
 * Later duplicates overwrite earlier entries (same as `querySelectorAll` iteration order).
 */
export function buildRadioMapByValue(
  scope: ParentNode,
  radiosSelector: string,
): Map<string, HTMLInputElement> {
  const m = new Map<string, HTMLInputElement>()
  for (const r of scope.querySelectorAll<HTMLInputElement>(radiosSelector)) {
    m.set(typeof r.value === 'string' ? r.value.trim() : '', r)
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
  const input = map.get(preferred.trim()) ?? map.get(fallbackKey.trim())
  if (!input) {
    return false
  }
  input.checked = true
  return true
}

/**
 * Checked `value` for a radio group `name` inside `form` (empty string if none or non-string).
 * Uses `HTMLFormElement.elements` / `RadioNodeList.value` instead of `:checked` `querySelector` on hot paths.
 * String values are **trimmed** so `buildRadioMapByValue` lookups stay aligned with browser quirks / odd markup.
 */
export function getFormRadioGroupValue(form: HTMLFormElement, name: string): string {
  const el = form.elements.namedItem(name)
  if (!el) {
    return ''
  }
  const v = (el as unknown as { value?: unknown }).value
  if (typeof v !== 'string') {
    return ''
  }
  return v.trim()
}
