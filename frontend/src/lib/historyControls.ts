import { HISTORY_SORT_SELECTOR } from './domSelectors'

/** Avoid duplicate `change` → `requestSubmit` if history wiring runs twice. */
const historySortWiredSelects = new WeakSet<HTMLSelectElement>()

/**
 * Resolve the history sort `<select>` from a root (usually `document`).
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
  const form = select.form
  if (!form) {
    return
  }
  historySortWiredSelects.add(select)
  select.addEventListener('change', () => {
    form.requestSubmit()
  })
}

export function initHistoryControls(): void {
  wireHistorySortAutoSubmit(queryHistorySortSelect(document))
}
