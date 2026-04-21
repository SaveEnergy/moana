import { resolveBootContentQueryRoot } from './contentRoot'
import { attachCategoryListEditDelegation } from './categoryModalEditDelegation'
import { attachCategoryModalFormPreviewListeners } from './categoryModalForm'
import { queryCategoryModalInitContext } from './categoryModalQueries'
import { createCategoryModalOpenFlow } from './categoryModalOpenFlow'
import { createCategoryModalPreviewController } from './categoryModalPreview'
import { attachNativeDialogDismiss } from './dialogDismiss'
import {
  CATEGORY_MODAL_COLOR_NATIVE_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  CATEGORY_MODAL_DISMISS_SELECTORS,
  CATEGORY_MODAL_SELECTOR,
} from './domSelectors'
import { buildRadioMapByValue } from './radioMap'

/** Full modal wiring once per `<dialog>` (duplicate `bootApp` must not stack listeners). */
const categoryModalInitialized = new WeakSet<HTMLDialogElement>()

export function initCategoryModal(): void {
  const contentRoot = resolveBootContentQueryRoot()
  const dialog = contentRoot.querySelector<HTMLDialogElement>(CATEGORY_MODAL_SELECTOR)
  if (!dialog) {
    return
  }
  if (categoryModalInitialized.has(dialog)) {
    return
  }

  const ctx = queryCategoryModalInitContext(contentRoot, dialog)
  if (!ctx) {
    return
  }

  /* Narrow once for nested functions (TS does not always narrow captured consts in closures). */
  const modal = dialog
  const catForm = ctx.form
  const catId = ctx.idInput
  const catTitle = ctx.titleEl
  const catSubmit = ctx.submitBtn
  const catPreview = ctx.preview
  const catIconWrap = ctx.iconWrap
  const catName = ctx.nameInput
  const addCategoryBtn = ctx.addCategoryBtn
  const editDelegationRoot = ctx.editDelegationRoot

  /** Resolved once — avoids repeated `querySelector` on every preview sync (`input` / `change`). */
  const colorNativeInput = catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_NATIVE_SELECTOR)

  const colorRadioByValue = buildRadioMapByValue(catForm, CATEGORY_MODAL_COLOR_RADIOS_SELECTOR)
  const iconRadioByValue = buildRadioMapByValue(catForm, CATEGORY_MODAL_ICON_RADIOS_SELECTOR)

  const catPreviewCtl = createCategoryModalPreviewController({
    form: catForm,
    colorNativeInput,
    iconRadioByValue,
    preview: catPreview,
    iconWrap: catIconWrap,
  })

  attachCategoryModalFormPreviewListeners(catForm, catPreviewCtl)

  const { openCreateModal, openEditModal } = createCategoryModalOpenFlow({
    modal,
    form: catForm,
    idInput: catId,
    titleEl: catTitle,
    submitBtn: catSubmit,
    nameInput: catName,
    colorNativeInput,
    colorRadioByValue,
    iconRadioByValue,
    previewCtl: catPreviewCtl,
  })

  addCategoryBtn?.addEventListener('click', () => openCreateModal())

  /** Scoped to the list card so topbar/sidebar clicks do not run this handler (`categories.html`: list is sibling of intro under the same parent). */
  if (editDelegationRoot) {
    attachCategoryListEditDelegation(editDelegationRoot, openEditModal)
  }

  attachNativeDialogDismiss(modal, CATEGORY_MODAL_DISMISS_SELECTORS)
  categoryModalInitialized.add(dialog)
}
