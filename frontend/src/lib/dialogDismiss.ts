import { clickEventTargetElement } from './clickTarget'

/** One dismiss `click` listener per dialog (safe if `attachNativeDialogDismiss` runs twice). */
const nativeDialogDismissWiredDialogs = new WeakSet<HTMLDialogElement>()

/** Backdrop (`target === dialog`) or `closest()` any of the selectors (inner nodes included). */
export function shouldCloseNativeDialogFromClick(
  e: MouseEvent,
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
    if (el.closest(sel) !== null) {
      return true
    }
  }
  return false
}

/**
 * One `click` listener: backdrop + close/cancel controls identified by `closest()`.
 * The dialog is recorded **after** `addEventListener` so a failed registration does not block a later retry.
 */
export function attachNativeDialogDismiss(
  dialog: HTMLDialogElement,
  closeWithinSelectors: readonly string[],
): void {
  if (nativeDialogDismissWiredDialogs.has(dialog)) {
    return
  }
  dialog.addEventListener('click', (e) => {
    if (shouldCloseNativeDialogFromClick(e, dialog, closeWithinSelectors)) {
      dialog.close()
    }
  })
  nativeDialogDismissWiredDialogs.add(dialog)
}
