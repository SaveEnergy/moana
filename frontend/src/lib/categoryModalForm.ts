import type { CategoryModalPreviewController } from './categoryModalPreview'
import {
  CATEGORY_COLOR_NATIVE_CLASS,
  CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
} from './domSelectors'

/** Forms that already have preview **`input`** / **`change`** listeners (`attachCategoryModalFormPreviewListeners` is safe to call more than once). */
const categoryModalFormPreviewWired = new WeakSet<HTMLFormElement>()

/**
 * Form-level preview delegation for the category modal: one `input` + one `change` listener
 * on `#cat-modal-form` instead of per-swatch wiring. Used by {@link initCategoryModal}.
 * Idempotent per form (**`WeakSet`**, same pattern as **`attachNativeDialogDismiss`**).
 */
export function attachCategoryModalFormPreviewListeners(
  form: HTMLFormElement,
  previewCtl: CategoryModalPreviewController,
): void {
  if (categoryModalFormPreviewWired.has(form)) {
    return
  }
  form.addEventListener('input', (e) => {
    const t = e.target
    if (!(t instanceof Element) || !t.classList.contains(CATEGORY_COLOR_NATIVE_CLASS)) {
      return
    }
    const wrap = t.closest(CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR)
    const r = wrap?.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR)
    if (r) {
      r.checked = true
      previewCtl.raf.schedule()
    }
  })
  form.addEventListener('change', (e) => {
    const t = e.target
    if (!(t instanceof HTMLInputElement)) {
      return
    }
    if (t.name === CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME) {
      previewCtl.sync({ colorRadioTarget: t })
      return
    }
    if (t.name === CATEGORY_MODAL_ICON_RADIO_GROUP_NAME) {
      previewCtl.sync({ iconRadioTarget: t })
    }
  })
  categoryModalFormPreviewWired.add(form)
}
