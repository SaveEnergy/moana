import { APP_MAIN_SELECTOR } from './domSelectors'

/** One `querySelector(APP_MAIN_SELECTOR)` per `document` per load — `bootApp` runs several initializers that all need this root. */
const contentRootMemo = new WeakMap<Document, ParentNode>()

/**
 * Prefer `<main.app-main>` for boot-time queries when `parent` is `document` (`layout.html`),
 * so scans skip shell chrome. Falls back to `document` when the landmark is absent (e.g. `login.html`).
 * Returns `parent` unchanged when it is not the global `document` (tests, explicit subtrees) or when `document` is undefined (Node).
 *
 * Boot modules call {@link resolveBootContentQueryRoot} (→ `document` here). **`applyLocalTimeElements`** uses the same when given the global `document`, else passes an explicit subtree. Other boot features: **`initConfirmSubmitForms`** ({@link queryBootContentAll} for **`form[data-confirm]`**), **`initCategoryModal`**, **`initHistoryControls`** ({@link queryBootContent} for **`#history-sort`**), **`initSettingsMemberDialog`**.
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

/**
 * Boot modules call this instead of inlining `resolveContentQueryRoot(document)` — same per-`document` WeakMap memo as {@link resolveContentQueryRoot}, clearer intent at call sites.
 */
export function resolveBootContentQueryRoot(): ParentNode {
  return resolveContentQueryRoot(document)
}

/**
 * `resolveBootContentQueryRoot().querySelector` — one call site for single-selector boot wiring
 * (e.g. **`initHistoryControls`**); keeps **`WeakMap`** memo reuse without repeating the resolve + query pattern.
 */
export function queryBootContent<E extends Element = Element>(selector: string): E | null {
  return resolveBootContentQueryRoot().querySelector<E>(selector)
}

/**
 * `resolveBootContentQueryRoot().querySelectorAll` — boot-scoped list queries (e.g. **`initConfirmSubmitForms`**).
 */
export function queryBootContentAll<E extends Element = Element>(selector: string): NodeListOf<E> {
  return resolveBootContentQueryRoot().querySelectorAll<E>(selector)
}
