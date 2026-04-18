export function initSettingsMemberDialog(): void {
  const dialog = document.getElementById('settings-add-member-dialog') as HTMLDialogElement | null
  if (!dialog) {
    return
  }

  const openBtn = document.getElementById('settings-add-member-open')

  openBtn?.addEventListener('click', () => {
    dialog.showModal()
  })

  /** Backdrop, × close, and Cancel — one listener; use closest so clicks on inner nodes still match. */
  dialog.addEventListener('click', (e) => {
    const el = e.target
    if (!(el instanceof Element)) {
      return
    }
    if (el === dialog) {
      dialog.close()
      return
    }
    if (
      el.closest('#settings-add-member-close') !== null ||
      el.closest('#settings-add-member-cancel') !== null
    ) {
      dialog.close()
    }
  })
}
