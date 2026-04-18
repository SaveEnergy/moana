import { attachNativeDialogDismiss } from './dialogDismiss'

export function initSettingsMemberDialog(): void {
  const dialog = document.getElementById('settings-add-member-dialog') as HTMLDialogElement | null
  if (!dialog) {
    return
  }

  const openBtn = document.getElementById('settings-add-member-open')

  openBtn?.addEventListener('click', () => {
    dialog.showModal()
  })

  attachNativeDialogDismiss(dialog, ['#settings-add-member-close', '#settings-add-member-cancel'])
}
