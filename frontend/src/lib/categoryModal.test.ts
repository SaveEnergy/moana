import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  APP_MAIN_SELECTOR,
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_PAGE_INTRO_SECTION_SELECTOR,
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
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'

describe('initCategoryModal', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('does not stack listeners when init runs twice with the same DOM', () => {
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

    const modalInnerBySelector: Record<string, unknown> = {
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
    }

    const dialog = {
      addEventListener: vi.fn(),
      closest: () => null,
      querySelector: vi.fn((sel: string) => modalInnerBySelector[sel] ?? null),
    }

    const introQuerySelector = vi.fn((sel: string) =>
      sel === CATEGORY_MODAL_OPEN_CREATE_SELECTOR ? addBtn : null,
    )
    const intro = { querySelector: introQuerySelector }

    vi.stubGlobal('document', {
      querySelector: (sel: string) => {
        if (sel === CATEGORY_MODAL_SELECTOR) return dialog
        if (sel === CATEGORY_PAGE_INTRO_SECTION_SELECTOR) return intro
        if (sel === CATEGORY_LIST_SECTION_SELECTOR) return null
        return null
      },
    })

    initCategoryModal()
    initCategoryModal()

    expect(introQuerySelector).toHaveBeenCalledWith(CATEGORY_MODAL_OPEN_CREATE_SELECTOR)
    expect(dialog.querySelector).toHaveBeenCalledTimes(7)
    expect(formAdd).toHaveBeenCalledTimes(2)
    expect(dialog.addEventListener).toHaveBeenCalledTimes(1)
    expect(addBtn.addEventListener).toHaveBeenCalledTimes(1)
  })

  it('resolves Add category from document when .cat-page-intro is absent', () => {
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

    const modalInnerBySelector: Record<string, unknown> = {
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
    }

    const dialog = {
      addEventListener: vi.fn(),
      closest: () => null,
      querySelector: vi.fn((sel: string) => modalInnerBySelector[sel] ?? null),
    }

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelector: (sel: string) => {
          if (sel === CATEGORY_MODAL_SELECTOR) return dialog
          if (sel === CATEGORY_PAGE_INTRO_SECTION_SELECTOR) return null
          if (sel === CATEGORY_MODAL_OPEN_CREATE_SELECTOR) return addBtn
          if (sel === CATEGORY_LIST_SECTION_SELECTOR) return null
          return null
        },
      }),
    )

    initCategoryModal()

    expect(addBtn.addEventListener).toHaveBeenCalledTimes(1)
  })

  it('resolves modal and intro from main.app-main when the landmark exists', () => {
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

    const modalInnerBySelector: Record<string, unknown> = {
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
    }

    const dialog = {
      addEventListener: vi.fn(),
      closest: () => null,
      querySelector: vi.fn((sel: string) => modalInnerBySelector[sel] ?? null),
    }

    const introQuerySelector = vi.fn((sel: string) =>
      sel === CATEGORY_MODAL_OPEN_CREATE_SELECTOR ? addBtn : null,
    )
    const intro = { querySelector: introQuerySelector }

    const mainQuerySelector = vi.fn((sel: string) => {
      if (sel === CATEGORY_MODAL_SELECTOR) return dialog
      if (sel === CATEGORY_PAGE_INTRO_SECTION_SELECTOR) return intro
      if (sel === CATEGORY_LIST_SECTION_SELECTOR) return null
      return null
    })
    const main = { querySelector: mainQuerySelector } as unknown as ParentNode

    const baseDoc = stubDocumentMainLandmark(main)
    const docQuerySelector = vi.fn((sel: string) => baseDoc.querySelector(sel))

    vi.stubGlobal('document', { querySelector: docQuerySelector })

    initCategoryModal()

    expect(docQuerySelector).toHaveBeenCalledWith(APP_MAIN_SELECTOR)
    expect(mainQuerySelector).toHaveBeenCalledWith(CATEGORY_MODAL_SELECTOR)
    expect(mainQuerySelector).toHaveBeenCalledWith(CATEGORY_PAGE_INTRO_SECTION_SELECTOR)
    expect(docQuerySelector.mock.calls.some((c) => c[0] === CATEGORY_MODAL_SELECTOR)).toBe(false)
    expect(addBtn.addEventListener).toHaveBeenCalledTimes(1)
  })

  it('attaches Edit delegation to .cat-list-section under intro parent before document fallback', () => {
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

    const modalInnerBySelector: Record<string, unknown> = {
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
    }

    const dialog = {
      addEventListener: vi.fn(),
      closest: () => null,
      querySelector: vi.fn((sel: string) => modalInnerBySelector[sel] ?? null),
    }

    const listSectionScoped = { addEventListener: vi.fn() }
    const listSectionDocumentFallback = { addEventListener: vi.fn() }
    const pageRootQuerySelector = vi.fn((sel: string) =>
      sel === CATEGORY_LIST_SECTION_SELECTOR ? listSectionScoped : null,
    )
    const pageRoot = { querySelector: pageRootQuerySelector }
    const introQuerySelector = vi.fn((sel: string) =>
      sel === CATEGORY_MODAL_OPEN_CREATE_SELECTOR ? addBtn : null,
    )
    const intro = { querySelector: introQuerySelector, parentElement: pageRoot }

    const docQuerySelector = vi.fn((sel: string) => {
      if (sel === CATEGORY_MODAL_SELECTOR) return dialog
      if (sel === CATEGORY_PAGE_INTRO_SECTION_SELECTOR) return intro
      if (sel === CATEGORY_LIST_SECTION_SELECTOR) return listSectionDocumentFallback
      return null
    })

    vi.stubGlobal('document', stubDocumentWithoutMainLandmark({ querySelector: docQuerySelector }))

    initCategoryModal()

    expect(pageRootQuerySelector).toHaveBeenCalledWith(CATEGORY_LIST_SECTION_SELECTOR)
    expect(listSectionScoped.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    expect(listSectionDocumentFallback.addEventListener).not.toHaveBeenCalled()
  })
})
