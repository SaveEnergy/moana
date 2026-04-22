import { resolveBootContentQueryRoot } from './contentRoot'
import { attachNativeDialogDismiss } from './dialogDismiss'
import { attachShowModalOnClick } from './dialogModal'
import {
  SETTINGS_AVATAR_DIALOG_SELECTOR,
  SETTINGS_AVATAR_DISMISS_SELECTORS,
} from './domSelectors'
import { querySettingsAvatarInitContext } from './settingsAvatarDialogQueries'
import { initSettingsAvatarDropzone } from './settingsAvatarDropzone'

const settingsAvatarDialogInitialized = new WeakSet<HTMLDialogElement>()

export function querySettingsAvatarDialog(root: ParentNode): HTMLDialogElement | null {
  return root.querySelector<HTMLDialogElement>(SETTINGS_AVATAR_DIALOG_SELECTOR)
}

export function initSettingsAvatarDialog(): void {
  const contentRoot = resolveBootContentQueryRoot()
  const dialog = contentRoot.querySelector<HTMLDialogElement>(SETTINGS_AVATAR_DIALOG_SELECTOR)
  if (!dialog) {
    return
  }
  if (settingsAvatarDialogInitialized.has(dialog)) {
    return
  }

  const { openBtn } = querySettingsAvatarInitContext(contentRoot, dialog)

  attachShowModalOnClick(openBtn, dialog)
  attachNativeDialogDismiss(dialog, SETTINGS_AVATAR_DISMISS_SELECTORS)
  initSettingsAvatarDropzone(dialog)
  settingsAvatarDialogInitialized.add(dialog)
}
