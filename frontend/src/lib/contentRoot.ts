import { APP_MAIN_SELECTOR } from './domSelectors'

/** One `querySelector(APP_MAIN_SELECTOR)` per `document` per load — `bootApp` runs several initializers that all need this root. */
const contentRootMemo = new WeakMap<Document, ParentNode>()

/**
 * During {@link withBootContentRootScope}, {@link resolveBootContentQueryRoot} returns this node without re-entering
 * {@link resolveContentQueryRoot} (one resolved landmark per boot pass; restores previous scope on exit).
 */
let activeBootContentRoot: ParentNode | null = null

/**
 * Prefer `<main.app-main>` for boot-time queries when `parent` is `document` (`layout.html`),
 * so scans skip shell chrome. Falls back to `document` when the landmark is absent (e.g. `login.html`).
 * Returns `parent` unchanged when it is not the global `document` (tests, explicit subtrees) or when `document` is undefined (Node).
 *
 * Boot modules call {@link resolveBootContentQueryRoot} (→ `document` here). A non-empty **`runBootInitializers`** pass (**`boot.ts`**) opens {@link withBootContentRootScope} so those calls share one pinned **`ParentNode`** without re-entering {@link resolveContentQueryRoot}. **`applyLocalTimeElements`** (default **`document`** arg): one {@link resolveBootContentQueryRoot} + **`querySelectorAll`** (**`LOCAL_TIME_ELEMENTS_SELECTOR`**); explicit subtrees: {@link resolveContentQueryRoot} then **`querySelectorAll`**. Other boot features: **`initConfirmSubmitForms`** (one {@link resolveBootContentQueryRoot} + **`querySelectorAll`** for **`form[data-confirm]`** — same **`NodeList`** as {@link queryBootContentAll} with **`FORM_DATA_CONFIRM_SELECTOR`**), **`initCategoryModal`**, **`initHistoryControls`** ({@link resolveBootContentQueryRoot} + **`queryHistorySortSelect`** — same **`ParentNode`** as {@link queryBootContent} with **`HISTORY_SORT_SELECTOR`**), **`initSettingsMemberDialog`** (one {@link resolveBootContentQueryRoot} for **`dialog`** **`querySelector`** + **`querySettingsAddMemberInitContext`**). {@link queryBootContent} / {@link queryBootContentAll} remain for ad-hoc queries and **tests** (equivalence to resolve + **`querySelector`** / **`querySelectorAll`**).
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
  if (activeBootContentRoot !== null) {
    return activeBootContentRoot
  }
  return resolveContentQueryRoot(document)
}

/**
 * Runs `fn` with {@link resolveBootContentQueryRoot} pinned to one {@link resolveContentQueryRoot}(`document`) result
 * for nested calls (a full `runBootInitializers` pass in `boot.ts`). No-ops the pin when `document` is undefined (Node tests).
 * Restores the previous scope when `fn` returns or throws (supports nesting).
 */
export function withBootContentRootScope<T>(fn: () => T): T {
  const prev = activeBootContentRoot
  activeBootContentRoot =
    typeof document !== 'undefined' ? resolveContentQueryRoot(document) : null
  try {
    return fn()
  } finally {
    activeBootContentRoot = prev
  }
}

/**
 * `resolveBootContentQueryRoot().querySelector` — convenience for single-selector boot wiring; keeps **`WeakMap`** memo reuse.
 * **`initHistoryControls`** / **`initSettingsMemberDialog`** call {@link resolveBootContentQueryRoot} once explicitly, then domain **`querySelector`** helpers.
 */
export function queryBootContent<E extends Element = Element>(selector: string): E | null {
  return resolveBootContentQueryRoot().querySelector<E>(selector)
}

/**
 * `resolveBootContentQueryRoot().querySelectorAll` — boot-scoped list queries (e.g. **`applyLocalTimeElements`** default path).
 */
export function queryBootContentAll<E extends Element = Element>(selector: string): NodeListOf<E> {
  return resolveBootContentQueryRoot().querySelectorAll<E>(selector)
}
