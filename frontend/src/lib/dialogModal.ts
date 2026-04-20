/**
 * Open a native `<dialog>` only when it is not already **`open`**
 * (double-invoke safe: repeat clicks, duplicate handlers, tests).
 */
export function showModalIfClosed(dialog: HTMLDialogElement): void {
  if (dialog.open) {
    return
  }
  dialog.showModal()
}

/**
 * Close a native `<dialog>` only when **`open !== false`**
 * (stubs without **`open`** still call **`close()`** — same rule as {@link attachNativeDialogDismiss}).
 */
export function closeDialogIfOpen(dialog: HTMLDialogElement): void {
  if (dialog.open === false) {
    return
  }
  dialog.close()
}
