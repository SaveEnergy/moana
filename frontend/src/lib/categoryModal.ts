import { CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT, sanitizeCategoryCustomHex } from './categoryColor'
import { resolveBootContentQueryRoot } from './contentRoot'
import { readCategoryEditRowDataset } from './categoryModalDataset'
import { createCategoryModalPreviewController } from './categoryModalPreview'
import { attachNativeDialogDismiss } from './dialogDismiss'
import { showModalIfClosed } from './dialogModal'
import { clickEventTargetElement } from './clickTarget'
import {
  CATEGORY_COLOR_NATIVE_CLASS,
  CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR,
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_PAGE_INTRO_SECTION_SELECTOR,
  CATEGORY_MODAL_COLOR_NATIVE_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  CATEGORY_MODAL_DISMISS_SELECTORS,
  CATEGORY_MODAL_FORM_SELECTOR,
  CATEGORY_MODAL_ID_INPUT_SELECTOR,
  CATEGORY_MODAL_NAME_SELECTOR,
  CATEGORY_MODAL_OPEN_CREATE_SELECTOR,
  CATEGORY_MODAL_OPEN_EDIT_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ICON_SELECTOR,
  CATEGORY_MODAL_PREVIEW_SELECTOR,
  CATEGORY_MODAL_SELECTOR,
  CATEGORY_MODAL_SUBMIT_SELECTOR,
  CATEGORY_MODAL_TITLE_SELECTOR,
} from './domSelectors'
import { buildRadioMapByValue, setRadioCheckedByValue } from './radioMap'

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

  function wireCategoryFormPreview() {
    catForm.addEventListener('input', (e) => {
      const t = e.target
      if (!(t instanceof Element) || !t.classList.contains(CATEGORY_COLOR_NATIVE_CLASS)) {
        return
      }
      const wrap = t.closest(CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR)
      const r = wrap?.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR)
      if (r) {
        r.checked = true
        catPreviewCtl.raf.schedule()
      }
    })
    catForm.addEventListener('change', (e) => {
      const t = e.target
      if (!(t instanceof HTMLInputElement)) {
        return
      }
      if (t.name === CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME) {
        catPreviewCtl.sync({ colorRadioTarget: t })
        return
      }
      if (t.name === CATEGORY_MODAL_ICON_RADIO_GROUP_NAME) {
        catPreviewCtl.sync({ iconRadioTarget: t })
      }
    })
  }

  wireCategoryFormPreview()

  function openCreateModal() {
    catPreviewCtl.raf.cancelPending()
    catPreviewCtl.resetPaintState()
    catForm.action = '/categories'
    catId.value = ''
    catTitle.textContent = 'New category'
    catSubmit.textContent = 'Create category'
    catForm.reset()
    setRadioCheckedByValue(colorRadioByValue, '', '')
    setRadioCheckedByValue(iconRadioByValue, '', '')
    if (colorNativeInput) colorNativeInput.value = CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT
    catPreviewCtl.sync()
    catName.focus()
    showModalIfClosed(modal)
  }

  function openEditModal(btn: HTMLElement) {
    catPreviewCtl.raf.cancelPending()
    catPreviewCtl.resetPaintState()
    const row = readCategoryEditRowDataset(btn.dataset)
    if (!row) {
      return
    }
    catForm.action = '/categories/update'
    catId.value = row.id
    catTitle.textContent = 'Edit category'
    catSubmit.textContent = 'Save changes'

    catName.value = row.name

    if (row.isCustom) {
      setRadioCheckedByValue(colorRadioByValue, CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM, '')
      if (colorNativeInput) colorNativeInput.value = sanitizeCategoryCustomHex(row.customHex)
    } else {
      setRadioCheckedByValue(colorRadioByValue, row.rawColor, '')
    }

    setRadioCheckedByValue(iconRadioByValue, row.iconVal, '')

    catPreviewCtl.sync()
    catName.focus()
    showModalIfClosed(modal)
  }

  addCategoryBtn?.addEventListener('click', () => openCreateModal())

  /** Scoped to the list card so topbar/sidebar clicks do not run this handler (`categories.html`: list is sibling of intro under the same parent). */
  const editDelegationRoot =
    categoriesPageRoot?.querySelector(CATEGORY_LIST_SECTION_SELECTOR) ??
    contentRoot.querySelector(CATEGORY_LIST_SECTION_SELECTOR)
  if (editDelegationRoot) {
    editDelegationRoot.addEventListener('click', (e) => {
      const el = clickEventTargetElement(e)
      if (!el) {
        return
      }
      const btn = el.closest(CATEGORY_MODAL_OPEN_EDIT_SELECTOR)
      if (!btn || !(btn instanceof HTMLElement)) {
        return
      }
      openEditModal(btn)
    })
  }

  attachNativeDialogDismiss(modal, CATEGORY_MODAL_DISMISS_SELECTORS)
  categoryModalInitialized.add(dialog)
}
