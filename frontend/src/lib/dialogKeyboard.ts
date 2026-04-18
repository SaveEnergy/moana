/**
 * True when an event path includes a native `dialog` element that is currently open
 * (Escape and other keys should not be handled by layered chrome, e.g. mobile shell).
 *
 * Uses `tagName === 'DIALOG'` + `open` so tests can run without a full DOM `instanceof`.
 */
export function eventPathIncludesOpenDialog(path: readonly unknown[]): boolean {
  for (const n of path) {
    if (!n || typeof n !== 'object' || !('tagName' in n)) {
      continue
    }
    const el = n as { tagName: string; open?: boolean }
    if (el.tagName === 'DIALOG' && el.open === true) {
      return true
    }
  }
  return false
}

export function keyEventInvolvesOpenDialog(e: KeyboardEvent): boolean {
  return eventPathIncludesOpenDialog(e.composedPath())
}
