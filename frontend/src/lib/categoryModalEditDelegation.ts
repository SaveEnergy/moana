import { clickEventTargetElement } from './clickTarget'
import { CATEGORY_MODAL_OPEN_EDIT_SELECTOR } from './domSelectors'

/**
 * Event delegation for **Edit** on the categories list card (`categories.html` **`.cat-list-section`**).
 * Keeps **`initCategoryModal`** from inlining **`closest`** / target resolution.
 */
export function attachCategoryListEditDelegation(
  listRoot: ParentNode,
  onEditClick: (btn: HTMLElement) => void,
): void {
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
}
