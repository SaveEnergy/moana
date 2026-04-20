import { APP_MAIN_SELECTOR } from './domSelectors'

/**
 * Prefer `<main.app-main>` for boot-time queries when `parent` is `document` (`layout.html`),
 * so scans skip shell chrome. Falls back to `document` when the landmark is absent (e.g. `login.html`).
 * Returns `parent` unchanged when it is not the global `document` (tests, explicit subtrees) or when `document` is undefined (Node).
 *
 * Used by `applyLocalTimeElements`, `initConfirmSubmitForms`, `initCategoryModal`, `initHistoryControls`, and `initSettingsMemberDialog`.
 */
export function resolveContentQueryRoot(parent: ParentNode): ParentNode {
  if (typeof document === 'undefined') {
    return parent
  }
  if (parent !== document) {
    return parent
  }
  return document.querySelector(APP_MAIN_SELECTOR) ?? document
}
