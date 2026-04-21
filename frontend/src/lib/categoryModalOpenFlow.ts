import { CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT, sanitizeCategoryCustomHex } from './categoryColor'
import { readCategoryEditRowDataset } from './categoryModalDataset'
import type { CategoryModalPreviewController } from './categoryModalPreview'
import { showModalIfClosed } from './dialogModal'
import { CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM } from './domSelectors'
import { setRadioCheckedByValue } from './radioMap'

/** Narrow modal open/edit dependencies — used by {@link initCategoryModal} only. */
export type CategoryModalOpenFlowContext = {
  modal: HTMLDialogElement
  form: HTMLFormElement
  idInput: HTMLInputElement
  titleEl: HTMLElement
  submitBtn: HTMLElement
  nameInput: HTMLInputElement
  colorNativeInput: HTMLInputElement | null | undefined
  colorRadioByValue: ReadonlyMap<string, HTMLInputElement>
  iconRadioByValue: ReadonlyMap<string, HTMLInputElement>
  previewCtl: CategoryModalPreviewController
}

/**
 * Create / edit category modal open handlers (reset + dataset-driven radios, then {@link showModalIfClosed}).
 * Keeps {@link initCategoryModal} focused on wiring and delegation.
 */
export function createCategoryModalOpenFlow(ctx: CategoryModalOpenFlowContext) {
  const {
    modal,
    form,
    idInput,
    titleEl,
    submitBtn,
    nameInput,
    colorNativeInput,
    colorRadioByValue,
    iconRadioByValue,
    previewCtl,
  } = ctx

  function openCreateModal() {
    previewCtl.raf.cancelPending()
    previewCtl.resetPaintState()
    form.action = '/categories'
    idInput.value = ''
    titleEl.textContent = 'New category'
    submitBtn.textContent = 'Create category'
    form.reset()
    setRadioCheckedByValue(colorRadioByValue, '', '')
    setRadioCheckedByValue(iconRadioByValue, '', '')
    if (colorNativeInput) colorNativeInput.value = CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT
    previewCtl.sync()
    nameInput.focus()
    showModalIfClosed(modal)
  }

  function openEditModal(btn: HTMLElement) {
    previewCtl.raf.cancelPending()
    previewCtl.resetPaintState()
    const row = readCategoryEditRowDataset(btn.dataset)
    if (!row) {
      return
    }
    form.action = '/categories/update'
    idInput.value = row.id
    titleEl.textContent = 'Edit category'
    submitBtn.textContent = 'Save changes'

    nameInput.value = row.name

    if (row.isCustom) {
      setRadioCheckedByValue(colorRadioByValue, CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM, '')
      if (colorNativeInput) colorNativeInput.value = sanitizeCategoryCustomHex(row.customHex)
    } else {
      setRadioCheckedByValue(colorRadioByValue, row.rawColor, '')
    }

    setRadioCheckedByValue(iconRadioByValue, row.iconVal, '')

    previewCtl.sync()
    nameInput.focus()
    showModalIfClosed(modal)
  }

  return { openCreateModal, openEditModal }
}
