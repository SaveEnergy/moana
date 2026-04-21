import { resolveBootContentQueryRoot } from './contentRoot'
import { attachCategoryListEditDelegation } from './categoryModalEditDelegation'
import { attachCategoryModalFormPreviewListeners } from './categoryModalForm'
import { createCategoryModalOpenFlow } from './categoryModalOpenFlow'
import { createCategoryModalPreviewController } from './categoryModalPreview'
import { attachNativeDialogDismiss } from './dialogDismiss'
import {
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_PAGE_INTRO_SECTION_SELECTOR,
  CATEGORY_MODAL_COLOR_NATIVE_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  CATEGORY_MODAL_DISMISS_SELECTORS,
  CATEGORY_MODAL_FORM_SELECTOR,
  CATEGORY_MODAL_ID_INPUT_SELECTOR,
  CATEGORY_MODAL_NAME_SELECTOR,
  CATEGORY_MODAL_OPEN_CREATE_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ICON_SELECTOR,
  CATEGORY_MODAL_PREVIEW_SELECTOR,
  CATEGORY_MODAL_SELECTOR,
  CATEGORY_MODAL_SUBMIT_SELECTOR,
  CATEGORY_MODAL_TITLE_SELECTOR,
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

  const form = dialog.querySelector<HTMLFormElement>(CATEGORY_MODAL_FORM_SELECTOR)
  const idInput = dialog.querySelector<HTMLInputElement>(CATEGORY_MODAL_ID_INPUT_SELECTOR)
  const titleEl = dialog.querySelector<HTMLElement>(CATEGORY_MODAL_TITLE_SELECTOR)
  const submitBtn = dialog.querySelector<HTMLElement>(CATEGORY_MODAL_SUBMIT_SELECTOR)
  const preview = dialog.querySelector<HTMLElement>(CATEGORY_MODAL_PREVIEW_SELECTOR)
  const iconWrap = dialog.querySelector<HTMLElement>(CATEGORY_MODAL_PREVIEW_ICON_SELECTOR)
  const nameInput = dialog.querySelector<HTMLInputElement>(CATEGORY_MODAL_NAME_SELECTOR)
  const intro = contentRoot.querySelector(CATEGORY_PAGE_INTRO_SECTION_SELECTOR)
  const categoriesPageRoot = intro?.parentElement
  const addCategoryBtn =
    intro?.querySelector<HTMLElement>(CATEGORY_MODAL_OPEN_CREATE_SELECTOR) ??
    contentRoot.querySelector<HTMLElement>(CATEGORY_MODAL_OPEN_CREATE_SELECTOR)

  if (!form || !idInput || !titleEl || !submitBtn || !preview || !iconWrap || !nameInput) {
    return
  }

  /* Narrow once for nested functions (TS does not always narrow captured consts in closures). */
  const modal = dialog
  const catForm = form
  const catId = idInput
  const catTitle = titleEl
  const catSubmit = submitBtn
  const catPreview = preview
  const catIconWrap = iconWrap
  const catName = nameInput

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
  const editDelegationRoot =
    categoriesPageRoot?.querySelector(CATEGORY_LIST_SECTION_SELECTOR) ??
    contentRoot.querySelector(CATEGORY_LIST_SECTION_SELECTOR)
  if (editDelegationRoot) {
    attachCategoryListEditDelegation(editDelegationRoot, openEditModal)
  }

  attachNativeDialogDismiss(modal, CATEGORY_MODAL_DISMISS_SELECTORS)
  categoryModalInitialized.add(dialog)
}
