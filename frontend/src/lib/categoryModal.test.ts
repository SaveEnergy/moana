import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  CATEGORY_LIST_SECTION_SELECTOR,
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
import { initCategoryModal } from './categoryModal'

describe('initCategoryModal', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('does not stack listeners when init runs twice with the same DOM', () => {
    const dialog = { addEventListener: vi.fn(), closest: () => null }
    const formAdd = vi.fn()
    const form = {
      addEventListener: formAdd,
      querySelector: vi.fn(() => null),
      querySelectorAll: vi.fn(() => [] as unknown as NodeListOf<HTMLInputElement>),
      action: '',
      reset: vi.fn(),
    }
    const idInput = { value: '' }
    const titleEl = { textContent: '' }
    const submitBtn = { textContent: '' }
    const preview = { style: {} }
    const iconWrap = {
      innerHTML: '',
      textContent: '',
      classList: { add: vi.fn(), remove: vi.fn() },
    }
    const nameInput = { value: '', focus: vi.fn() }
    const addBtn = { addEventListener: vi.fn() }

    const bySelector: Record<string, unknown> = {
      [CATEGORY_MODAL_SELECTOR]: dialog,
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
      [CATEGORY_MODAL_OPEN_CREATE_SELECTOR]: addBtn,
      [CATEGORY_LIST_SECTION_SELECTOR]: null,
    }

    vi.stubGlobal('document', {
      querySelector: (sel: string) => bySelector[sel] ?? null,
    })

    initCategoryModal()
    initCategoryModal()

    expect(formAdd).toHaveBeenCalledTimes(2)
    expect(dialog.addEventListener).toHaveBeenCalledTimes(1)
    expect(addBtn.addEventListener).toHaveBeenCalledTimes(1)
  })
})
