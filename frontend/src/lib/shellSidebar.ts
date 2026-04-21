import { setAttributeIfChanged } from './domAttribute'
import { shouldDeferMobileShellEscape } from './dialogKeyboard'
import {
  APP_SHELL_SELECTOR,
  APP_SIDEBAR_BACKDROP_SELECTOR,
  APP_SIDEBAR_TOGGLE_SELECTOR,
} from './domSelectors'
import { shouldCloseMobileSidebarFromShellClick } from './mobileShellDismiss'
import { MOBILE_SHELL_MEDIA_QUERY, onMediaQueryChange } from './shellBreakpoints'

const shellSidebarWiredShells = new WeakSet<HTMLElement>()

/** Partial **`initShellSidebar`** retries (throw before **`shellSidebarWiredShells`** records the shell) must not stack these listeners. */
const shellSidebarToggleClickWired = new WeakSet<HTMLElement>()
const shellSidebarShellClickWired = new WeakSet<HTMLElement>()

/** `#app-shell` — mobile drawer root (`layout.html`). */
export function queryAppShell(root: ParentNode): HTMLElement | null {
  return root.querySelector<HTMLElement>(APP_SHELL_SELECTOR)
}

/** `#app-sidebar-toggle` — menu affordance (`layout.html`). */
export function querySidebarToggle(root: ParentNode): HTMLElement | null {
  return root.querySelector<HTMLElement>(APP_SIDEBAR_TOGGLE_SELECTOR)
}

/** `#app-sidebar-backdrop` — dimmed overlay (`layout.html`). */
export function querySidebarBackdrop(root: ParentNode): HTMLElement | null {
  return root.querySelector<HTMLElement>(APP_SIDEBAR_BACKDROP_SELECTOR)
}

/**
 * Mobile drawer: open/close sidebar, sync aria and backdrop.
 * Uses {@link setAttributeIfChanged} — avoids redundant ARIA writes on repeated close (e.g. Escape).
 */
export function initShellSidebar(): void {
  const shell = queryAppShell(document)
  if (!shell) {
    return
  }
  if (shellSidebarWiredShells.has(shell)) {
    return
  }
  const appShell = shell
  /* Toggle + backdrop are under `#app-shell` in `layout.html` — scope queries to the shell subtree. */
  const toggle = querySidebarToggle(shell)
  const backdrop = querySidebarBackdrop(shell)

  const mqMobile = window.matchMedia(MOBILE_SHELL_MEDIA_QUERY)

  function setExpanded(open: boolean) {
    if (!toggle) {
      return
    }
    const expanded = open ? 'true' : 'false'
    const label = open ? 'Close navigation menu' : 'Open navigation menu'
    setAttributeIfChanged(toggle, 'aria-expanded', expanded)
    setAttributeIfChanged(toggle, 'aria-label', label)
  }

  function openMobileSidebar() {
    appShell.classList.add('sidebar-open')
    if (backdrop) {
      setAttributeIfChanged(backdrop, 'aria-hidden', 'false')
    }
    setExpanded(true)
  }

  function closeMobileSidebar() {
    appShell.classList.remove('sidebar-open')
    if (backdrop) {
      setAttributeIfChanged(backdrop, 'aria-hidden', 'true')
    }
    setExpanded(false)
  }

  function toggleMobileSidebar() {
    if (appShell.classList.contains('sidebar-open')) {
      closeMobileSidebar()
    } else {
      openMobileSidebar()
    }
  }

  if (toggle && !shellSidebarToggleClickWired.has(toggle)) {
    toggle.addEventListener('click', () => {
      if (!mqMobile.matches) {
        return
      }
      toggleMobileSidebar()
    })
    shellSidebarToggleClickWired.add(toggle)
  }

  /** Backdrop + in-drawer close — one bubbling listener (see `mobileShellDismiss.ts`). */
  if (!shellSidebarShellClickWired.has(appShell)) {
    appShell.addEventListener('click', (e) => {
      if (!mqMobile.matches) {
        return
      }
      if (!appShell.classList.contains('sidebar-open')) {
        return
      }
      if (shouldCloseMobileSidebarFromShellClick(e, backdrop)) {
        closeMobileSidebar()
      }
    })
    shellSidebarShellClickWired.add(appShell)
  }

  onMediaQueryChange(mqMobile, () => {
    if (!mqMobile.matches) {
      closeMobileSidebar()
    }
  })

  /* Capture phase: decide before bubble so we still see `dialog.open` during the same Escape.
   * When the drawer is already closed, skip defer work (`composedPath` + optional menu `querySelector`) — hot path on mobile. */
  document.addEventListener(
    'keydown',
    (e) => {
      if (e.key !== 'Escape' || !mqMobile.matches) {
        return
      }
      if (!appShell.classList.contains('sidebar-open')) {
        return
      }
      if (shouldDeferMobileShellEscape(e)) {
        return
      }
      closeMobileSidebar()
    },
    { capture: true },
  )

  shellSidebarWiredShells.add(shell)
}
