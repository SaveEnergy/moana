/**
 * History page (GET `/history`): changing sort submits the form (replaces inline `onchange`).
 * Exported for unit tests; `initHistoryControls` is the boot entry.
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
  wireHistorySortAutoSubmit(document.getElementById('history-sort') as HTMLSelectElement | null)
}
