import { shouldDeferMobileShellEscape } from './dialogKeyboard'
import {
  APP_SHELL_SELECTOR,
  APP_SIDEBAR_BACKDROP_SELECTOR,
  APP_SIDEBAR_TOGGLE_SELECTOR,
} from './domSelectors'
import { shouldCloseMobileSidebarFromShellClick } from './mobileShellDismiss'
import { MOBILE_SHELL_MEDIA_QUERY, onMediaQueryChange } from './shellBreakpoints'

const shellSidebarWiredShells = new WeakSet<HTMLElement>()

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

/** Mobile drawer: open/close sidebar, sync aria and backdrop. */
export function initShellSidebar(): void {
  const shell = queryAppShell(document)
  if (!shell) {
    return
  }
  if (shellSidebarWiredShells.has(shell)) {
    return
  }
  shellSidebarWiredShells.add(shell)
  const appShell = shell
  const toggle = querySidebarToggle(document)
  const backdrop = querySidebarBackdrop(document)

  const mqMobile = window.matchMedia(MOBILE_SHELL_MEDIA_QUERY)

  function setExpanded(open: boolean) {
    toggle?.setAttribute('aria-expanded', open ? 'true' : 'false')
    toggle?.setAttribute('aria-label', open ? 'Close navigation menu' : 'Open navigation menu')
  }

  function openMobileSidebar() {
    appShell.classList.add('sidebar-open')
    backdrop?.setAttribute('aria-hidden', 'false')
    setExpanded(true)
  }

  function closeMobileSidebar() {
    appShell.classList.remove('sidebar-open')
    backdrop?.setAttribute('aria-hidden', 'true')
    setExpanded(false)
  }

  function toggleMobileSidebar() {
    if (appShell.classList.contains('sidebar-open')) {
      closeMobileSidebar()
    } else {
      openMobileSidebar()
    }
  }

  toggle?.addEventListener('click', () => {
    if (!mqMobile.matches) {
      return
    }
    toggleMobileSidebar()
  })

  /** Backdrop + in-drawer close — one bubbling listener (see `mobileShellDismiss.ts`). */
  appShell.addEventListener('click', (e) => {
    if (!mqMobile.matches) {
      return
    }
    if (shouldCloseMobileSidebarFromShellClick(e, backdrop)) {
      closeMobileSidebar()
    }
  })

  onMediaQueryChange(mqMobile, () => {
    if (!mqMobile.matches) {
      closeMobileSidebar()
    }
  })

  /* Capture phase: decide before bubble so we still see `dialog.open` during the same Escape. */
  document.addEventListener(
    'keydown',
    (e) => {
      if (e.key !== 'Escape' || !mqMobile.matches) {
        return
      }
      if (shouldDeferMobileShellEscape(e)) {
        return
      }
      closeMobileSidebar()
    },
    { capture: true },
  )
}
