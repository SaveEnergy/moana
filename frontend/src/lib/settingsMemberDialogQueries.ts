import { SETTINGS_ADD_MEMBER_OPEN_SELECTOR } from './domSelectors'

/**
 * Open button for the settings add-member **`dialog`** (may be **`null`** when the template omits it — dismiss still attaches).
 * {@link initSettingsMemberDialog} runs after the idempotent **`WeakSet`** short-circuit on the **`dialog`**.
 */
export type SettingsAddMemberInitContext = {
  openBtn: HTMLElement | null
}

/**
 * **`SETTINGS_ADD_MEMBER_OPEN_SELECTOR`** on **`dialog.parentElement`** when present (`settings.html` sibling), else **`contentRoot`**
 * (same **`ParentNode`** as the single **`resolveBootContentQueryRoot()`** in **`initSettingsMemberDialog`** — equivalent to **`queryBootContent`** for the **`dialog`** **`querySelector`**).
 */
export function querySettingsAddMemberInitContext(
  contentRoot: ParentNode,
  dialog: HTMLDialogElement,
): SettingsAddMemberInitContext {
  const openRoot = dialog.parentElement ?? contentRoot
  return {
    openBtn: openRoot.querySelector<HTMLElement>(SETTINGS_ADD_MEMBER_OPEN_SELECTOR),
  }
}
