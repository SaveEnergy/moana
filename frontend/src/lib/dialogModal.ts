/** Open controls that already have a **`click` → {@link showModalIfClosed}** listener (`attachShowModalOnClick` is safe to call more than once). */
const showModalOnClickWiredOpenButtons = new WeakSet<HTMLElement>()

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
 * Idempotent per button (same **`WeakSet`** pattern as **`attachNativeDialogDismiss`**).
 */
export function attachShowModalOnClick(
  openBtn: HTMLElement | null | undefined,
  dialog: HTMLDialogElement,
): void {
  if (!openBtn) {
    return
  }
  if (showModalOnClickWiredOpenButtons.has(openBtn)) {
    return
  }
  openBtn.addEventListener('click', () => {
    showModalIfClosed(dialog)
  })
  showModalOnClickWiredOpenButtons.add(openBtn)
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
