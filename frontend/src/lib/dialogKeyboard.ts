import { APP_USER_MENU_OPEN_SELECTOR } from './domSelectors'

function readTagOpen(n: unknown): { tag: string; open: boolean } | null {
  if (!n || typeof n !== 'object' || !('tagName' in n)) {
    return null
  }
  const el = n as { tagName: unknown; open?: boolean }
  if (typeof el.tagName !== 'string') {
    return null
  }
  /* HTML uses uppercase `tagName`; normalize so tests and odd paths still match. */
  return { tag: el.tagName.toUpperCase(), open: el.open === true }
}

/** Uses `tagName` + `open` so tests can run without a full DOM `instanceof`. */
function pathIncludesOpenTag(path: readonly unknown[], tag: 'DIALOG' | 'DETAILS'): boolean {
  for (const n of path) {
    const r = readTagOpen(n)
    if (r && r.tag === tag && r.open) {
      return true
    }
  }
  return false
}

/**
 * True when an event path includes a native `dialog` element that is currently open
 * (Escape and other keys should not be handled by layered chrome, e.g. mobile shell).
 */
export function eventPathIncludesOpenDialog(path: readonly unknown[]): boolean {
  return pathIncludesOpenTag(path, 'DIALOG')
}

/**
 * True when the path includes an open `<details>` (e.g. focus inside the disclosure). See also {@link isAppUserMenuDetailsOpen}.
 */
export function eventPathIncludesOpenDetails(path: readonly unknown[]): boolean {
  return pathIncludesOpenTag(path, 'DETAILS')
}

export function keyEventInvolvesOpenDialog(e: KeyboardEvent): boolean {
  return eventPathIncludesOpenDialog(e.composedPath())
}

/** Topbar account menu (`layout.html`); DOM query because the mobile drawer can cover the bar while `[open]` stays true. */
export function isAppUserMenuDetailsOpen(): boolean {
  return document.querySelector(APP_USER_MENU_OPEN_SELECTOR) !== null
}

/** One pass over `composedPath()` — used by {@link shouldDeferMobileShellEscape} to avoid two scans + a query when the path already answers. */
function pathIncludesOpenDialogOrDetails(path: readonly unknown[]): boolean {
  for (const n of path) {
    const r = readTagOpen(n)
    if (r && r.open && (r.tag === 'DIALOG' || r.tag === 'DETAILS')) {
      return true
    }
  }
  return false
}

/** Single guard for `shellSidebar` Escape: open `dialog` or open `<details>` in path, else open account menu via DOM (drawer can cover the bar while `[open]` stays true). */
export function shouldDeferMobileShellEscape(e: KeyboardEvent): boolean {
  const path = e.composedPath()
  return pathIncludesOpenDialogOrDetails(path) || isAppUserMenuDetailsOpen()
}
