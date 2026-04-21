import { queryBootContent, resolveBootContentQueryRoot } from './contentRoot'
import { attachNativeDialogDismiss } from './dialogDismiss'
import { showModalIfClosed } from './dialogModal'
import {
  SETTINGS_ADD_MEMBER_DIALOG_SELECTOR,
  SETTINGS_ADD_MEMBER_DISMISS_SELECTORS,
  SETTINGS_ADD_MEMBER_OPEN_SELECTOR,
} from './domSelectors'

/** One init pass per dialog (`bootApp` may run more than once during tests or future SPA hooks). */
const settingsMemberDialogInitialized = new WeakSet<HTMLDialogElement>()

/** Resolve the add-member `<dialog>` from a root (`initSettingsMemberDialog` uses {@link queryBootContent}). */
export function querySettingsAddMemberDialog(root: ParentNode): HTMLDialogElement | null {
  return root.querySelector<HTMLDialogElement>(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR)
}

/** Resolve the control that opens the add-member dialog (`settings.html`). */
export function querySettingsAddMemberOpenButton(root: ParentNode): HTMLElement | null {
  return root.querySelector<HTMLElement>(SETTINGS_ADD_MEMBER_OPEN_SELECTOR)
}

export function initSettingsMemberDialog(): void {
  const dialog = queryBootContent<HTMLDialogElement>(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR)
  if (!dialog) {
    return
  }
  if (settingsMemberDialogInitialized.has(dialog)) {
    return
  }

  /* Open control sits next to the dialog under the same parent (`settings.html`); fall back to boot content root. */
  const openBtn = querySettingsAddMemberOpenButton(dialog.parentElement ?? resolveBootContentQueryRoot())

  openBtn?.addEventListener('click', () => {
    showModalIfClosed(dialog)
  })

  attachNativeDialogDismiss(dialog, SETTINGS_ADD_MEMBER_DISMISS_SELECTORS)
  settingsMemberDialogInitialized.add(dialog)
}
