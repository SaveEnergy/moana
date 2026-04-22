import { SETTINGS_AVATAR_OPEN_SELECTOR } from './domSelectors'

export type SettingsAvatarInitContext = {
  openBtn: HTMLElement | null
}

/** Open control for the profile photo dialog — same `parentElement ?? contentRoot` rule as add-member. */
export function querySettingsAvatarInitContext(
  contentRoot: ParentNode,
  dialog: HTMLDialogElement,
): SettingsAvatarInitContext {
  const openRoot = dialog.parentElement ?? contentRoot
  return {
    openBtn: openRoot.querySelector<HTMLElement>(SETTINGS_AVATAR_OPEN_SELECTOR),
  }
}
