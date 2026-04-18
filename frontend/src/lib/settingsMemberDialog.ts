export function initSettingsMemberDialog(): void {
  const dialog = document.getElementById('settings-add-member-dialog') as HTMLDialogElement | null
  const openBtn = document.getElementById('settings-add-member-open')
  const closeBtn = document.getElementById('settings-add-member-close')
  const cancelBtn = document.getElementById('settings-add-member-cancel')

  openBtn?.addEventListener('click', () => {
    dialog?.showModal()
  })

  closeBtn?.addEventListener('click', () => {
    dialog?.close()
  })

  cancelBtn?.addEventListener('click', () => {
    dialog?.close()
  })

  dialog?.addEventListener('click', (e) => {
    if (e.target === dialog) {
      dialog.close()
    }
  })
}
