import { resolveBootContentQueryRoot } from './contentRoot'
import { attachNativeDialogDismiss } from './dialogDismiss'
import { attachShowModalOnClick } from './dialogModal'
import { querySettingsAddMemberInitContext } from './settingsMemberDialogQueries'
import {
  SETTINGS_ADD_MEMBER_DIALOG_SELECTOR,
  SETTINGS_ADD_MEMBER_DISMISS_SELECTORS,
  SETTINGS_ADD_MEMBER_OPEN_SELECTOR,
} from './domSelectors'

/** One init pass per dialog (`bootApp` may run more than once during tests or future SPA hooks). */
const settingsMemberDialogInitialized = new WeakSet<HTMLDialogElement>()

/** Resolve the add-member `<dialog>` from a root (`initSettingsMemberDialog` uses one {@link resolveBootContentQueryRoot} for `contentRoot`, then `querySelector`). */
export function querySettingsAddMemberDialog(root: ParentNode): HTMLDialogElement | null {
  return root.querySelector<HTMLDialogElement>(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR)
}

/** Resolve the control that opens the add-member dialog (`settings.html`). */
export function querySettingsAddMemberOpenButton(root: ParentNode): HTMLElement | null {
  return root.querySelector<HTMLElement>(SETTINGS_ADD_MEMBER_OPEN_SELECTOR)
}

export function initSettingsMemberDialog(): void {
  const contentRoot = resolveBootContentQueryRoot()
  const dialog = contentRoot.querySelector<HTMLDialogElement>(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR)
  if (!dialog) {
    return
  }
  if (settingsMemberDialogInitialized.has(dialog)) {
    return
  }

  const { openBtn } = querySettingsAddMemberInitContext(contentRoot, dialog)

  attachShowModalOnClick(openBtn, dialog)

  attachNativeDialogDismiss(dialog, SETTINGS_ADD_MEMBER_DISMISS_SELECTORS)
  settingsMemberDialogInitialized.add(dialog)
}
