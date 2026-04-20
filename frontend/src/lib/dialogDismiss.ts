import { type ClickTargetEvent, clickEventTargetElement } from './clickTarget'
import { closeDialogIfOpen } from './dialogModal'

/** One dismiss `click` listener per dialog (safe if `attachNativeDialogDismiss` runs twice). */
const nativeDialogDismissWiredDialogs = new WeakSet<HTMLDialogElement>()

/** Backdrop (`target === dialog`) or `closest()` any of the selectors (inner nodes included). */
export function shouldCloseNativeDialogFromClick(
  e: ClickTargetEvent,
  dialog: HTMLDialogElement,
  closeWithinSelectors: readonly string[],
): boolean {
  const el = clickEventTargetElement(e)
  if (!el) {
    return false
  }
  if (el === dialog) {
    return true
  }
  for (const sel of closeWithinSelectors) {
    /* Invalid or empty selectors throw from `closest`; skip blanks so callers can pad lists safely. */
    if (!sel.trim()) {
      continue
    }
    if (el.closest(sel) !== null) {
      return true
    }
  }
  return false
}

/**
 * One `click` listener: backdrop + close/cancel controls identified by `closest()`.
 * Uses {@link closeDialogIfOpen} — skips redundant **`close()`** when **`dialog.open === false`**.
 * The dialog is recorded **after** `addEventListener` so a failed registration does not block a later retry.
 * Whitespace-only dismiss selectors are filtered once here so hot `click` paths skip them.
 */
export function attachNativeDialogDismiss(
  dialog: HTMLDialogElement,
  closeWithinSelectors: readonly string[],
): void {
  if (nativeDialogDismissWiredDialogs.has(dialog)) {
    return
  }
  const selectors: readonly string[] =
    closeWithinSelectors.length > 0
      ? closeWithinSelectors.filter((s) => s.trim() !== '')
      : closeWithinSelectors
  dialog.addEventListener('click', (e) => {
    if (!shouldCloseNativeDialogFromClick(e, dialog, selectors)) {
      return
    }
    closeDialogIfOpen(dialog)
  })
  nativeDialogDismissWiredDialogs.add(dialog)
}
