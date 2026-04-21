import { type ClickTargetEvent, clickEventTargetElement } from './clickTarget'
import { closeDialogIfOpen } from './dialogModal'
import { trimEdgesIfNeeded } from './trimEdges'

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
  for (let i = 0, n = closeWithinSelectors.length; i < n; i++) {
    /* Invalid or empty selectors throw from `closest`; skip blanks so callers can pad lists safely. */
    const t = trimEdgesIfNeeded(closeWithinSelectors[i])
    if (!t) {
      continue
    }
    if (el.closest(t) !== null) {
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
  let selectors: readonly string[]
  if (closeWithinSelectors.length > 0) {
    const out: string[] = []
    for (let i = 0, n = closeWithinSelectors.length; i < n; i++) {
      const t = trimEdgesIfNeeded(closeWithinSelectors[i])
      if (t) {
        out.push(t)
      }
    }
    selectors = out
  } else {
    selectors = closeWithinSelectors
  }
  dialog.addEventListener('click', (e) => {
    if (dialog.open === false) {
      return
    }
    if (!shouldCloseNativeDialogFromClick(e, dialog, selectors)) {
      return
    }
    closeDialogIfOpen(dialog)
  })
  nativeDialogDismissWiredDialogs.add(dialog)
}
