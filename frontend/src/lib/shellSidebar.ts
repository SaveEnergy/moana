import { isAppUserMenuDetailsOpen, keyEventInvolvesOpenDialog } from './dialogKeyboard'
import { MOBILE_SHELL_MEDIA_QUERY, onMediaQueryChange } from './shellBreakpoints'

/** Mobile drawer: open/close sidebar, sync aria and backdrop. */
export function initShellSidebar(): void {
  const shell = document.getElementById('app-shell')
  if (!shell) {
    return
  }
  const appShell = shell
  const toggle = document.getElementById('app-sidebar-toggle')
  const closeBtn = document.getElementById('app-sidebar-close')
  const backdrop = document.getElementById('app-sidebar-backdrop')

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

  backdrop?.addEventListener('click', () => {
    closeMobileSidebar()
  })

  closeBtn?.addEventListener('click', () => {
    closeMobileSidebar()
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
      if (keyEventInvolvesOpenDialog(e) || isAppUserMenuDetailsOpen()) {
        return
      }
      closeMobileSidebar()
    },
    { capture: true },
  )
}
