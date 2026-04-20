/**
 * Category modal and similar UIs: radios are indexed by `value` (including "").
 * Avoids duplicate if/else + non-null assertions when selecting by dataset-driven values.
 */

import { trimEdgesIfNeeded } from './trimEdges'

/**
 * Build a map of trimmed `HTMLInputElement.value` → input for radios under `scope` matching `radiosSelector`.
 * Keys use **`trimEdgesIfNeeded(value)`** so {@link getFormRadioGroupValue} and {@link setRadioCheckedByValue} stay consistent.
 * Later duplicates overwrite earlier entries (same as `querySelectorAll` iteration order).
 */
export function buildRadioMapByValue(
  scope: ParentNode,
  radiosSelector: string,
): Map<string, HTMLInputElement> {
  const m = new Map<string, HTMLInputElement>()
  for (const r of scope.querySelectorAll<HTMLInputElement>(radiosSelector)) {
    m.set(typeof r.value === 'string' ? trimEdgesIfNeeded(r.value) : '', r)
  }
  return m
}

/**
 * Select a radio from a map keyed by `HTMLInputElement.value`.
 * If `preferred` is not a key, selects `fallbackKey` (default `""` for auto/unset radios).
 * Skips assigning **`checked`** when that input is already selected (no redundant DOM writes / events).
 * @returns whether a matching input existed and was checked.
 */
export function setRadioCheckedByValue(
  map: ReadonlyMap<string, HTMLInputElement>,
  preferred: string,
  fallbackKey = '',
): boolean {
  const pref = trimEdgesIfNeeded(preferred)
  let input = map.get(pref)
  if (input === undefined) {
    const fb = trimEdgesIfNeeded(fallbackKey)
    if (fb !== pref) {
      input = map.get(fb)
    }
  }
  if (!input) {
    return false
  }
  if (!input.checked) {
    input.checked = true
  }
  return true
}

/**
 * Checked `value` for a radio group `name` inside `form` (empty string if none or non-string).
 * Uses `HTMLFormElement.elements` / `RadioNodeList.value` instead of `:checked` `querySelector` on hot paths.
 * String values use **`trimEdgesIfNeeded`** so `buildRadioMapByValue` lookups stay aligned with browser quirks / odd markup.
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
  return trimEdgesIfNeeded(v)
}
