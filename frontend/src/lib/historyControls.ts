import { resolveBootContentQueryRoot } from './contentRoot'
import { HISTORY_SORT_SELECTOR } from './domSelectors'

/** Avoid duplicate `change` → `requestSubmit` if history wiring runs twice. */
const historySortWiredSelects = new WeakSet<HTMLSelectElement>()

/**
 * Resolve the history sort `<select>` from a root (`initHistoryControls` uses `resolveBootContentQueryRoot()`).
 * Uses {@link HISTORY_SORT_SELECTOR} so the id string lives only in `domSelectors.ts`.
 */
export function queryHistorySortSelect(root: ParentNode): HTMLSelectElement | null {
  return root.querySelector<HTMLSelectElement>(HISTORY_SORT_SELECTOR)
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
 * History page: wire sort `<select>` → `form.requestSubmit()` when present.
 * Skips when that element is already wired (duplicate `bootApp`); see {@link wireHistorySortAutoSubmit}.
 */
export function initHistoryControls(): void {
  wireHistorySortAutoSubmit(queryHistorySortSelect(resolveBootContentQueryRoot()))
}
