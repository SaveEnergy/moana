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
 * Wire a single “open” control to {@link showModalIfClosed} (settings add-member, other one-button dialogs).
 */
export function attachShowModalOnClick(
  openBtn: HTMLElement | null | undefined,
  dialog: HTMLDialogElement,
): void {
  openBtn?.addEventListener('click', () => {
    showModalIfClosed(dialog)
  })
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
