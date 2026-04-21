import { APP_USER_MENU_OPEN_SELECTOR } from './domSelectors'

/**
 * When **`open === true`**, returns normalized tag ( **`DIALOG` / `DETAILS`** via cheap string compare, else **`toUpperCase`** ).
 * When not open, returns **`null`** — skips tag work on typical composed-path nodes.
 */
function readOpenTagNameIfOpen(n: unknown): string | null {
  if (!n || typeof n !== 'object' || !('tagName' in n)) {
    return null
  }
  const el = n as { tagName: unknown; open?: boolean }
  if (el.open !== true || typeof el.tagName !== 'string') {
    return null
  }
  const raw = el.tagName
  if (raw === 'DIALOG' || raw === 'dialog') {
    return 'DIALOG'
  }
  if (raw === 'DETAILS' || raw === 'details') {
    return 'DETAILS'
  }
  return raw.toUpperCase()
}

/**
 * True when an event path includes a native `dialog` element that is currently open
 * (Escape and other keys should not be handled by layered chrome, e.g. mobile shell).
 */
export function eventPathIncludesOpenDialog(path: readonly unknown[]): boolean {
  for (let i = 0, n = path.length; i < n; i++) {
    if (readOpenTagNameIfOpen(path[i]) === 'DIALOG') {
      return true
    }
  }
  return false
}

/**
 * True when the path includes an open `<details>` (e.g. focus inside the disclosure). See also {@link isAppUserMenuDetailsOpen}.
 */
export function eventPathIncludesOpenDetails(path: readonly unknown[]): boolean {
  for (let i = 0, n = path.length; i < n; i++) {
    if (readOpenTagNameIfOpen(path[i]) === 'DETAILS') {
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

/** One pass over `composedPath()` — used by {@link shouldDeferMobileShellEscape} to avoid two scans + a query when the path already answers. */
function pathIncludesOpenDialogOrDetails(path: readonly unknown[]): boolean {
  for (let i = 0, n = path.length; i < n; i++) {
    const tag = readOpenTagNameIfOpen(path[i])
    if (tag === 'DIALOG' || tag === 'DETAILS') {
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
