import {
  CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT,
  resolveCategoryModalPreviewBackground,
  sanitizeCategoryCustomHex,
} from './categoryColor'
import { readCategoryEditRowDataset } from './categoryModalDataset'
import { attachNativeDialogDismiss } from './dialogDismiss'
import { clickEventTargetElement } from './clickTarget'
import {
  CATEGORY_COLOR_NATIVE_CLASS,
  CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR,
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_MODAL_COLOR_NATIVE_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
  CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  CATEGORY_MODAL_DISMISS_SELECTORS,
  CATEGORY_MODAL_FORM_SELECTOR,
  CATEGORY_MODAL_ID_INPUT_SELECTOR,
  CATEGORY_MODAL_NAME_SELECTOR,
  CATEGORY_MODAL_OPEN_CREATE_SELECTOR,
  CATEGORY_MODAL_OPEN_EDIT_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS,
  CATEGORY_MODAL_PREVIEW_ICON_SELECTOR,
  CATEGORY_MODAL_PREVIEW_SELECTOR,
  CATEGORY_MODAL_SELECTOR,
  CATEGORY_MODAL_SUBMIT_SELECTOR,
  CATEGORY_MODAL_TITLE_SELECTOR,
  MOANA_ICON_CAT_PREVIEW_CLASS,
  MOANA_ICON_SVG_SELECTOR,
} from './domSelectors'
import { buildRadioMapByValue, getFormRadioGroupValue, setRadioCheckedByValue } from './radioMap'

/** Full modal wiring once per `<dialog>` (duplicate `bootApp` must not stack listeners). */
const categoryModalInitialized = new WeakSet<HTMLDialogElement>()

export function initCategoryModal(): void {
  const dialog = document.querySelector<HTMLDialogElement>(CATEGORY_MODAL_SELECTOR)
  if (!dialog) {
    return
  }

  const form = document.querySelector<HTMLFormElement>(CATEGORY_MODAL_FORM_SELECTOR)
  const idInput = document.querySelector<HTMLInputElement>(CATEGORY_MODAL_ID_INPUT_SELECTOR)
  const titleEl = document.querySelector<HTMLElement>(CATEGORY_MODAL_TITLE_SELECTOR)
  const submitBtn = document.querySelector<HTMLElement>(CATEGORY_MODAL_SUBMIT_SELECTOR)
  const preview = document.querySelector<HTMLElement>(CATEGORY_MODAL_PREVIEW_SELECTOR)
  const iconWrap = document.querySelector<HTMLElement>(CATEGORY_MODAL_PREVIEW_ICON_SELECTOR)
  const nameInput = document.querySelector<HTMLInputElement>(CATEGORY_MODAL_NAME_SELECTOR)
  const addCategoryBtn = document.querySelector<HTMLElement>(CATEGORY_MODAL_OPEN_CREATE_SELECTOR)

  if (!form || !idInput || !titleEl || !submitBtn || !preview || !iconWrap || !nameInput) {
    return
  }
  if (categoryModalInitialized.has(dialog)) {
    return
  }
  categoryModalInitialized.add(dialog)

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

  function syncCatModalPreview() {
    const colorVal = getFormRadioGroupValue(catForm, CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME)
    catPreview.style.background = resolveCategoryModalPreviewBackground(
      colorVal || undefined,
      colorNativeInput?.value,
    )

    const iconVal = getFormRadioGroupValue(catForm, CATEGORY_MODAL_ICON_RADIO_GROUP_NAME)
    const ir =
      iconRadioByValue.get(iconVal) ??
      catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR)
    catIconWrap.innerHTML = ''
    if (!ir?.value) {
      catIconWrap.classList.add(CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS)
      catIconWrap.textContent = 'A'
      return
    }
    catIconWrap.classList.remove(CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS)
    const label = ir.closest('label')
    const svg = label?.querySelector(MOANA_ICON_SVG_SELECTOR)
    if (svg) {
      const clone = svg.cloneNode(true) as SVGElement
      clone.classList.add(MOANA_ICON_CAT_PREVIEW_CLASS)
      catIconWrap.appendChild(clone)
    }
  }

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
        syncCatModalPreview()
      }
    })
    catForm.addEventListener('change', (e) => {
      const t = e.target
      if (!(t instanceof HTMLInputElement)) {
        return
      }
      if (t.name === CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME || t.name === CATEGORY_MODAL_ICON_RADIO_GROUP_NAME) {
        syncCatModalPreview()
      }
    })
  }

  wireCategoryFormPreview()

  function openCreateModal() {
    catForm.action = '/categories'
    catId.value = ''
    catTitle.textContent = 'New category'
    catSubmit.textContent = 'Create category'
    catForm.reset()
    setRadioCheckedByValue(colorRadioByValue, '', '')
    setRadioCheckedByValue(iconRadioByValue, '', '')
    if (colorNativeInput) colorNativeInput.value = CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT
    syncCatModalPreview()
    catName.focus()
    modal.showModal()
  }

  function openEditModal(btn: HTMLElement) {
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

    syncCatModalPreview()
    catName.focus()
    modal.showModal()
  }

  addCategoryBtn?.addEventListener('click', () => openCreateModal())

  /** Scoped to the list card so topbar/sidebar clicks do not run this handler. */
  const editDelegationRoot = document.querySelector(CATEGORY_LIST_SECTION_SELECTOR)
  if (editDelegationRoot) {
    editDelegationRoot.addEventListener('click', (e) => {
      const el = clickEventTargetElement(e as MouseEvent)
      if (!el) {
        return
      }
      const btn = el.closest(CATEGORY_MODAL_OPEN_EDIT_SELECTOR)
      if (!btn) {
        return
      }
      openEditModal(btn as HTMLElement)
    })
  }

  attachNativeDialogDismiss(modal, CATEGORY_MODAL_DISMISS_SELECTORS)
}
