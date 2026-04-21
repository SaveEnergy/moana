import {
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_PAGE_INTRO_SECTION_SELECTOR,
  CATEGORY_MODAL_FORM_SELECTOR,
  CATEGORY_MODAL_ID_INPUT_SELECTOR,
  CATEGORY_MODAL_NAME_SELECTOR,
  CATEGORY_MODAL_OPEN_CREATE_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ICON_SELECTOR,
  CATEGORY_MODAL_PREVIEW_SELECTOR,
  CATEGORY_MODAL_SUBMIT_SELECTOR,
  CATEGORY_MODAL_TITLE_SELECTOR,
} from './domSelectors'

/**
 * Resolved DOM for wiring the categories **`dialog`** (after **`#cat-modal`** is found).
 * {@link initCategoryModal} checks **`categoryModalInitialized`** before calling this so duplicate boots skip inner **`querySelector`** work.
 */
export type CategoryModalInitContext = {
  form: HTMLFormElement
  idInput: HTMLInputElement
  titleEl: HTMLElement
  submitBtn: HTMLElement
  preview: HTMLElement
  iconWrap: HTMLElement
  nameInput: HTMLInputElement
  addCategoryBtn: HTMLElement | null
  editDelegationRoot: ParentNode | null
}

/**
 * **`dialog.querySelector`** for required **`CATEGORY_MODAL_*`** controls plus page chrome (**Add category**, list **`Edit`** root).
 * Returns **`null`** when the dialog exists but any required inner node is missing.
 */
export function queryCategoryModalInitContext(
  contentRoot: ParentNode,
  dialog: HTMLDialogElement,
): CategoryModalInitContext | null {
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
    return null
  }

  const editDelegationRoot =
    categoriesPageRoot?.querySelector(CATEGORY_LIST_SECTION_SELECTOR) ??
    contentRoot.querySelector(CATEGORY_LIST_SECTION_SELECTOR)

  return {
    form,
    idInput,
    titleEl,
    submitBtn,
    preview,
    iconWrap,
    nameInput,
    addCategoryBtn,
    editDelegationRoot,
  }
}
