import { describe, expect, it, vi } from 'vitest'

import { queryCategoryModalInitContext } from './categoryModalQueries'
import {
  CATEGORY_LIST_SECTION_SELECTOR,
  CATEGORY_MODAL_FORM_SELECTOR,
  CATEGORY_MODAL_ID_INPUT_SELECTOR,
  CATEGORY_MODAL_NAME_SELECTOR,
  CATEGORY_MODAL_OPEN_CREATE_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ICON_SELECTOR,
  CATEGORY_MODAL_PREVIEW_SELECTOR,
  CATEGORY_MODAL_SUBMIT_SELECTOR,
  CATEGORY_MODAL_TITLE_SELECTOR,
  CATEGORY_PAGE_INTRO_SECTION_SELECTOR,
} from './domSelectors'

function modalInnerStub(map: Record<string, unknown>): HTMLDialogElement {
  return {
    querySelector: vi.fn((sel: string) => (map[sel] as Element | null | undefined) ?? null),
  } as unknown as HTMLDialogElement
}

describe('queryCategoryModalInitContext', () => {
  it('returns null when any required inner control is missing', () => {
    const dialog = modalInnerStub({})
    const contentRoot = { querySelector: vi.fn(() => null) } as unknown as ParentNode
    expect(queryCategoryModalInitContext(contentRoot, dialog)).toBeNull()
    expect(dialog.querySelector).toHaveBeenCalled()
  })

  it('resolves addCategoryBtn from intro before contentRoot fallback', () => {
    const form = { _: 'form' }
    const idInput = { _: 'id' }
    const titleEl = { _: 'title' }
    const submitBtn = { _: 'submit' }
    const preview = { _: 'preview' }
    const iconWrap = { _: 'icon' }
    const nameInput = { _: 'name' }
    const addFromIntro = { _: 'add-intro' }
    const addFromRoot = { _: 'add-root' }

    const inner: Record<string, unknown> = {
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
    }
    const dialog = modalInnerStub(inner)

    const intro = {
      querySelector: vi.fn((sel: string) =>
        sel === CATEGORY_MODAL_OPEN_CREATE_SELECTOR ? addFromIntro : null,
      ),
    }
    const contentRoot = {
      querySelector: vi.fn((sel: string) => {
        if (sel === CATEGORY_PAGE_INTRO_SECTION_SELECTOR) return intro
        if (sel === CATEGORY_MODAL_OPEN_CREATE_SELECTOR) return addFromRoot
        if (sel === CATEGORY_LIST_SECTION_SELECTOR) return null
        return null
      }),
    } as unknown as ParentNode

    const out = queryCategoryModalInitContext(contentRoot, dialog)
    expect(out).not.toBeNull()
    expect(out!.addCategoryBtn).toBe(addFromIntro)
  })

  it('prefers .cat-list-section under intro parent over contentRoot', () => {
    const form = { _: 'form' }
    const idInput = { _: 'id' }
    const titleEl = { _: 'title' }
    const submitBtn = { _: 'submit' }
    const preview = { _: 'preview' }
    const iconWrap = { _: 'icon' }
    const nameInput = { _: 'name' }

    const inner: Record<string, unknown> = {
      [CATEGORY_MODAL_FORM_SELECTOR]: form,
      [CATEGORY_MODAL_ID_INPUT_SELECTOR]: idInput,
      [CATEGORY_MODAL_TITLE_SELECTOR]: titleEl,
      [CATEGORY_MODAL_SUBMIT_SELECTOR]: submitBtn,
      [CATEGORY_MODAL_PREVIEW_SELECTOR]: preview,
      [CATEGORY_MODAL_PREVIEW_ICON_SELECTOR]: iconWrap,
      [CATEGORY_MODAL_NAME_SELECTOR]: nameInput,
    }
    const dialog = modalInnerStub(inner)

    const listScoped = { _: 'list-scoped' }
    const listFallback = { _: 'list-fallback' }
    const pageRoot = {
      querySelector: vi.fn((sel: string) =>
        sel === CATEGORY_LIST_SECTION_SELECTOR ? listScoped : null,
      ),
    }
    const intro = {
      querySelector: vi.fn(() => null),
      parentElement: pageRoot,
    }
    const contentRoot = {
      querySelector: vi.fn((sel: string) => {
        if (sel === CATEGORY_PAGE_INTRO_SECTION_SELECTOR) return intro
        if (sel === CATEGORY_LIST_SECTION_SELECTOR) return listFallback
        return null
      }),
    } as unknown as ParentNode

    const out = queryCategoryModalInitContext(contentRoot, dialog)
    expect(out).not.toBeNull()
    expect(out!.editDelegationRoot).toBe(listScoped)
    expect(pageRoot.querySelector).toHaveBeenCalledWith(CATEGORY_LIST_SECTION_SELECTOR)
  })
})
