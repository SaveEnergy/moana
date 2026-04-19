import { attachNativeDialogDismiss } from './dialogDismiss'
import {
  SETTINGS_ADD_MEMBER_DIALOG_SELECTOR,
  SETTINGS_ADD_MEMBER_DISMISS_SELECTORS,
  SETTINGS_ADD_MEMBER_OPEN_SELECTOR,
} from './domSelectors'

/** Resolve the add-member `<dialog>` from a root (usually `document`). */
export function querySettingsAddMemberDialog(root: ParentNode): HTMLDialogElement | null {
  return root.querySelector<HTMLDialogElement>(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR)
}

/** Resolve the control that opens the add-member dialog (`settings.html`). */
export function querySettingsAddMemberOpenButton(root: ParentNode): HTMLElement | null {
  return root.querySelector<HTMLElement>(SETTINGS_ADD_MEMBER_OPEN_SELECTOR)
}

export function initSettingsMemberDialog(): void {
  const dialog = querySettingsAddMemberDialog(document)
  if (!dialog) {
    return
  }

  const openBtn = querySettingsAddMemberOpenButton(document)

  openBtn?.addEventListener('click', () => {
    dialog.showModal()
  })

  attachNativeDialogDismiss(dialog, SETTINGS_ADD_MEMBER_DISMISS_SELECTORS)
}
