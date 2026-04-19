import { attachNativeDialogDismiss } from './dialogDismiss'
import {
  SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID,
  SETTINGS_ADD_MEMBER_DISMISS_SELECTORS,
  SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID,
} from './domSelectors'

export function initSettingsMemberDialog(): void {
  const dialog = document.getElementById(SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID) as HTMLDialogElement | null
  if (!dialog) {
    return
  }

  const openBtn = document.getElementById(SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID)

  openBtn?.addEventListener('click', () => {
    dialog.showModal()
  })

  attachNativeDialogDismiss(dialog, SETTINGS_ADD_MEMBER_DISMISS_SELECTORS)
}
