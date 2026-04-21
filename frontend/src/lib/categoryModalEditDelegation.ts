import { clickEventTargetElement } from './clickTarget'
import { CATEGORY_MODAL_OPEN_EDIT_SELECTOR } from './domSelectors'

/** List roots that already have an **Edit** delegation **`click`** listener. */
const categoryListEditDelegationWired = new WeakSet<ParentNode>()

/**
 * Event delegation for **Edit** on the categories list card (`categories.html` **`.cat-list-section`**).
 * Keeps **`initCategoryModal`** from inlining **`closest`** / target resolution.
 * Idempotent per **`listRoot`** (**`WeakSet`**).
 */
export function attachCategoryListEditDelegation(
  listRoot: ParentNode,
  onEditClick: (btn: HTMLElement) => void,
): void {
  if (categoryListEditDelegationWired.has(listRoot)) {
    return
  }
  listRoot.addEventListener('click', (e) => {
    const el = clickEventTargetElement(e)
    if (!el) {
      return
    }
    const btn = el.closest(CATEGORY_MODAL_OPEN_EDIT_SELECTOR)
    if (!btn || !(btn instanceof HTMLElement)) {
      return
    }
    onEditClick(btn)
  })
  categoryListEditDelegationWired.add(listRoot)
}
