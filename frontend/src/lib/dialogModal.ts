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
