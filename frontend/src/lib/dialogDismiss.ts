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
  return closeWithinSelectors.some((sel) => el.closest(sel) !== null)
}

/** One `click` listener: backdrop + close/cancel controls identified by `closest()`. */
export function attachNativeDialogDismiss(
  dialog: HTMLDialogElement,
  closeWithinSelectors: readonly string[],
): void {
  if (nativeDialogDismissWiredDialogs.has(dialog)) {
    return
  }
  nativeDialogDismissWiredDialogs.add(dialog)
  dialog.addEventListener('click', (e) => {
    if (shouldCloseNativeDialogFromClick(e, dialog, closeWithinSelectors)) {
      dialog.close()
    }
  })
}
