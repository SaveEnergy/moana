import './app.css'

const shell = document.getElementById('app-shell')
const toggle = document.getElementById('app-sidebar-toggle')
const backdrop = document.getElementById('app-sidebar-backdrop')

const mqMobile = window.matchMedia('(max-width: 768px)')

function setExpanded(open: boolean) {
  toggle?.setAttribute('aria-expanded', open ? 'true' : 'false')
}

function openMobileSidebar() {
  shell?.classList.add('sidebar-open')
  backdrop?.setAttribute('aria-hidden', 'false')
  setExpanded(true)
}

function closeMobileSidebar() {
  shell?.classList.remove('sidebar-open')
  backdrop?.setAttribute('aria-hidden', 'true')
  setExpanded(false)
}

function toggleMobileSidebar() {
  if (shell?.classList.contains('sidebar-open')) {
    closeMobileSidebar()
  } else {
    openMobileSidebar()
  }
}

toggle?.addEventListener('click', () => {
  if (mqMobile.matches) {
    toggleMobileSidebar()
  } else {
    shell?.classList.toggle('sidebar-collapsed')
  }
})

backdrop?.addEventListener('click', () => {
  closeMobileSidebar()
})

mqMobile.addEventListener('change', () => {
  if (!mqMobile.matches) {
    closeMobileSidebar()
  }
})

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && mqMobile.matches) {
    closeMobileSidebar()
  }
})
