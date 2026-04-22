import { resolveBootContentQueryRoot } from './contentRoot'
import { HISTORY_N_SELECTOR, HISTORY_SORT_SELECTOR } from './domSelectors'

/** Avoid duplicate `change` → `requestSubmit` if history wiring runs twice. */
const historySortWiredSelects = new WeakSet<HTMLSelectElement>()
const historyRowsWiredSelects = new WeakSet<HTMLSelectElement>()

/**
 * Resolve the history sort `<select>` from a root (`initHistoryControls` uses {@link resolveBootContentQueryRoot} then this helper).
 * Uses {@link HISTORY_SORT_SELECTOR} so the id string lives only in `domSelectors.ts`.
 */
export function queryHistorySortSelect(root: ParentNode): HTMLSelectElement | null {
  return root.querySelector<HTMLSelectElement>(HISTORY_SORT_SELECTOR)
}

/**
 * Max-rows `<select name="n">` — see {@link HISTORY_N_SELECTOR}.
 */
export function queryHistoryRowsSelect(root: ParentNode): HTMLSelectElement | null {
  return root.querySelector<HTMLSelectElement>(HISTORY_N_SELECTOR)
}

/**
 * History page (GET `/history`): changing sort submits the form (replaces inline `onchange`).
 */
export function wireHistorySortAutoSubmit(select: HTMLSelectElement | null): void {
  if (!select) {
    return
  }
  if (historySortWiredSelects.has(select)) {
    return
  }
  if (!select.form) {
    return
  }
  /* Read `select.form` when `change` fires — not the closure from wire time — so a reassigned `form` still submits the right target. */
  select.addEventListener('change', () => {
    select.form?.requestSubmit()
  })
  historySortWiredSelects.add(select)
}

/**
 * History page: `n=` row cap (GET): changing the select submits the form (same pattern as sort).
 */
export function wireHistoryRowsAutoSubmit(select: HTMLSelectElement | null): void {
  if (!select) {
    return
  }
  if (historyRowsWiredSelects.has(select)) {
    return
  }
  if (!select.form) {
    return
  }
  select.addEventListener('change', () => {
    select.form?.requestSubmit()
  })
  historyRowsWiredSelects.add(select)
}

/**
 * History page: wires sort and max-rows `<select>`s when **`select.form`** exists; on **`change`**, **`requestSubmit`**
 * (form read at event time). Skips when an element is already wired (duplicate `bootApp`).
 */
export function initHistoryControls(): void {
  const root = resolveBootContentQueryRoot()
  wireHistorySortAutoSubmit(queryHistorySortSelect(root))
  wireHistoryRowsAutoSubmit(queryHistoryRowsSelect(root))
}
