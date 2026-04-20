import { APP_MAIN_SELECTOR } from './domSelectors'

/** One `querySelector(APP_MAIN_SELECTOR)` per `document` per load — `bootApp` runs several initializers that all need this root. */
const contentRootMemo = new WeakMap<Document, ParentNode>()

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
  const doc = document as Document
  const hit = contentRootMemo.get(doc)
  if (hit !== undefined) {
    return hit
  }
  const root = doc.querySelector(APP_MAIN_SELECTOR) ?? doc
  contentRootMemo.set(doc, root)
  return root
}
