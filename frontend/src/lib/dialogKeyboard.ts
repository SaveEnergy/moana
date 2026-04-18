import { APP_USER_MENU_OPEN_SELECTOR } from './domSelectors'

function readTagOpen(n: unknown): { tag: string; open: boolean } | null {
  if (!n || typeof n !== 'object' || !('tagName' in n)) {
    return null
  }
  const el = n as { tagName: string; open?: boolean }
  return { tag: el.tagName, open: el.open === true }
}

/**
 * True when an event path includes a native `dialog` element that is currently open
 * (Escape and other keys should not be handled by layered chrome, e.g. mobile shell).
 *
 * Uses `tagName === 'DIALOG'` + `open` so tests can run without a full DOM `instanceof`.
 */
export function eventPathIncludesOpenDialog(path: readonly unknown[]): boolean {
  for (const n of path) {
    const r = readTagOpen(n)
    if (r?.tag === 'DIALOG' && r.open) {
      return true
    }
  }
  return false
}

/**
 * True when the path includes an open `<details>` (e.g. focus inside the disclosure). See also {@link isAppUserMenuDetailsOpen}.
 */
export function eventPathIncludesOpenDetails(path: readonly unknown[]): boolean {
  for (const n of path) {
    const r = readTagOpen(n)
    if (r?.tag === 'DETAILS' && r.open) {
      return true
    }
  }
  return false
}

export function keyEventInvolvesOpenDialog(e: KeyboardEvent): boolean {
  return eventPathIncludesOpenDialog(e.composedPath())
}

/** Topbar account menu (`layout.html`); DOM query because the mobile drawer can cover the bar while `[open]` stays true. */
export function isAppUserMenuDetailsOpen(): boolean {
  return document.querySelector(APP_USER_MENU_OPEN_SELECTOR) !== null
}

/** Single guard for `shellSidebar` Escape: open `dialog` in path or open account `<details>`. */
export function shouldDeferMobileShellEscape(e: KeyboardEvent): boolean {
  return keyEventInvolvesOpenDialog(e) || isAppUserMenuDetailsOpen()
}
