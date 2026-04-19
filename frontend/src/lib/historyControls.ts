import { HISTORY_SORT_SELECTOR } from './domSelectors'

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
  const form = select.form
  if (!form) {
    return
  }
  select.addEventListener('change', () => {
    form.requestSubmit()
  })
}

export function initHistoryControls(): void {
  wireHistorySortAutoSubmit(queryHistorySortSelect(document))
}
