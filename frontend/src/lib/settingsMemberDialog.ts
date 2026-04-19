import { attachNativeDialogDismiss } from './dialogDismiss'
import {
  SETTINGS_ADD_MEMBER_DIALOG_SELECTOR,
  SETTINGS_ADD_MEMBER_DISMISS_SELECTORS,
  SETTINGS_ADD_MEMBER_OPEN_SELECTOR,
} from './domSelectors'

/** One init pass per dialog (`bootApp` may run more than once during tests or future SPA hooks). */
const settingsMemberDialogInitialized = new WeakSet<HTMLDialogElement>()

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
  if (settingsMemberDialogInitialized.has(dialog)) {
    return
  }

  /* Open control sits in `#settings-top` next to the dialog (`settings.html`); scope off parent when attached. */
  const openRoot = dialog.parentElement
  const openBtn = openRoot
    ? querySettingsAddMemberOpenButton(openRoot)
    : querySettingsAddMemberOpenButton(document)

  openBtn?.addEventListener('click', () => {
    dialog.showModal()
  })

  attachNativeDialogDismiss(dialog, SETTINGS_ADD_MEMBER_DISMISS_SELECTORS)
  settingsMemberDialogInitialized.add(dialog)
}
