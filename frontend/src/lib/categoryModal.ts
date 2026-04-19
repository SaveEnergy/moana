import { CATEGORY_MODAL_DEFAULT_PREVIEW_BG, sanitizeCategoryCustomHex } from './categoryColor'
import { attachNativeDialogDismiss } from './dialogDismiss'
import { clickEventTargetElement } from './clickTarget'
import {
  CATEGORY_COLOR_NATIVE_CLASS,
  CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR,
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_MODAL_COLOR_NATIVE_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
  CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  CATEGORY_MODAL_DISMISS_SELECTORS,
  CATEGORY_MODAL_ELEMENT_ID,
  CATEGORY_MODAL_FORM_ELEMENT_ID,
  CATEGORY_MODAL_ID_INPUT_ELEMENT_ID,
  CATEGORY_MODAL_NAME_ELEMENT_ID,
  CATEGORY_MODAL_OPEN_CREATE_ELEMENT_ID,
  CATEGORY_MODAL_OPEN_EDIT_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ELEMENT_ID,
  CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS,
  CATEGORY_MODAL_PREVIEW_ICON_ELEMENT_ID,
  CATEGORY_MODAL_SUBMIT_ELEMENT_ID,
  CATEGORY_MODAL_TITLE_ELEMENT_ID,
  MOANA_ICON_CAT_PREVIEW_CLASS,
  MOANA_ICON_SVG_SELECTOR,
} from './domSelectors'
import { setRadioCheckedByValue } from './radioMap'

export function initCategoryModal(): void {
  const dialog = document.getElementById(CATEGORY_MODAL_ELEMENT_ID) as HTMLDialogElement | null
  if (!dialog) {
    return
  }

  const form = document.getElementById(CATEGORY_MODAL_FORM_ELEMENT_ID) as HTMLFormElement | null
  const idInput = document.getElementById(CATEGORY_MODAL_ID_INPUT_ELEMENT_ID) as HTMLInputElement | null
  const titleEl = document.getElementById(CATEGORY_MODAL_TITLE_ELEMENT_ID)
  const submitBtn = document.getElementById(CATEGORY_MODAL_SUBMIT_ELEMENT_ID)
  const preview = document.getElementById(CATEGORY_MODAL_PREVIEW_ELEMENT_ID)
  const iconWrap = document.getElementById(CATEGORY_MODAL_PREVIEW_ICON_ELEMENT_ID)
  const nameInput = document.getElementById(CATEGORY_MODAL_NAME_ELEMENT_ID) as HTMLInputElement | null
  const addCategoryBtn = document.getElementById(CATEGORY_MODAL_OPEN_CREATE_ELEMENT_ID)

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

  const CATEGORY_MODAL_RADIO_GROUP_SELECTORS = {
    color: CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
    icon: CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  } as const

  function radiosByValue(kind: keyof typeof CATEGORY_MODAL_RADIO_GROUP_SELECTORS): Map<string, HTMLInputElement> {
    const m = new Map<string, HTMLInputElement>()
    for (const r of catForm.querySelectorAll<HTMLInputElement>(
      CATEGORY_MODAL_RADIO_GROUP_SELECTORS[kind],
    )) {
      m.set(r.value, r)
    }
    return m
  }

  const colorRadioByValue = radiosByValue('color')
  const iconRadioByValue = radiosByValue('icon')

  function syncCatModalPreview() {
    const cr = catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR)
    let bg: string = CATEGORY_MODAL_DEFAULT_PREVIEW_BG
    if (cr?.value === CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM) {
      const nat = catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_NATIVE_SELECTOR)
      bg = nat?.value?.trim() || '#818cf8'
    } else if (cr?.value) {
      bg = cr.value
    }
    catPreview.style.background = bg

    const ir = catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR)
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
    const nat = catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_NATIVE_SELECTOR)
    if (nat) nat.value = '#818cf8'
    syncCatModalPreview()
    catName.focus()
    modal.showModal()
  }

  function openEditModal(btn: HTMLElement) {
    const id = btn.dataset.id
    if (!id) return
    catForm.action = '/categories/update'
    catId.value = id
    catTitle.textContent = 'Edit category'
    catSubmit.textContent = 'Save changes'

    catName.value = btn.dataset.name ?? ''

    const rawColor = (btn.dataset.color ?? '').trim()
    const isCustom = btn.dataset.custom === '1'
    const customHex = (btn.dataset.customHex ?? '#818cf8').trim()

    if (isCustom) {
      setRadioCheckedByValue(colorRadioByValue, CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM, '')
      const nat = catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_COLOR_NATIVE_SELECTOR)
      if (nat) nat.value = sanitizeCategoryCustomHex(customHex)
    } else {
      setRadioCheckedByValue(colorRadioByValue, rawColor, '')
    }

    const iconVal = (btn.dataset.icon ?? '').trim()
    setRadioCheckedByValue(iconRadioByValue, iconVal, '')

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
