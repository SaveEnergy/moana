import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import { attachCategoryListEditDelegation } from './categoryModalEditDelegation'
import { stubClickTargetEvent } from './clickTarget'
import { CATEGORY_MODAL_OPEN_EDIT_SELECTOR } from './domSelectors'

describe('attachCategoryListEditDelegation', () => {
  beforeEach(() => {
    vi.stubGlobal('HTMLElement', class HTMLElement {})
  })
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('calls onEditClick when target resolves an edit control', () => {
    const onEditClick = vi.fn()
    const matched = new globalThis.HTMLElement() as HTMLElement
    const el = {
      closest: vi.fn((sel: string) => (sel === CATEGORY_MODAL_OPEN_EDIT_SELECTOR ? matched : null)),
    } as unknown as Element

    let clickHandler: ((e: unknown) => void) | undefined
    const listRoot = {
      addEventListener: vi.fn((type: string, fn: (e: unknown) => void) => {
        if (type === 'click') clickHandler = fn
      }),
    } as unknown as ParentNode

    attachCategoryListEditDelegation(listRoot, onEditClick)

    expect(listRoot.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    clickHandler!(stubClickTargetEvent(el))

    expect(el.closest).toHaveBeenCalledWith(CATEGORY_MODAL_OPEN_EDIT_SELECTOR)
    expect(onEditClick).toHaveBeenCalledOnce()
    expect(onEditClick.mock.calls[0]![0]).toBe(matched)
  })

  it('ignores clicks that do not resolve an edit button', () => {
    const onEditClick = vi.fn()
    const el = {
      closest: vi.fn(() => null),
    } as unknown as Element

    let clickHandler: ((e: unknown) => void) | undefined
    const listRoot = {
      addEventListener: vi.fn((_type: string, fn: (e: unknown) => void) => {
        clickHandler = fn
      }),
    } as unknown as ParentNode

    attachCategoryListEditDelegation(listRoot, onEditClick)
    clickHandler!(stubClickTargetEvent(el))

    expect(onEditClick).not.toHaveBeenCalled()
  })

  it('ignores closest match that is not an HTMLElement', () => {
    const onEditClick = vi.fn()
    const plain = { closest: () => ({}) }
    let clickHandler: ((e: unknown) => void) | undefined
    const listRoot = {
      addEventListener: vi.fn((_type: string, fn: (e: unknown) => void) => {
        clickHandler = fn
      }),
    } as unknown as ParentNode

    attachCategoryListEditDelegation(listRoot, onEditClick)
    clickHandler!(stubClickTargetEvent(plain))

    expect(onEditClick).not.toHaveBeenCalled()
  })

  it('ignores null target', () => {
    const onEditClick = vi.fn()
    let clickHandler: ((e: unknown) => void) | undefined
    const listRoot = {
      addEventListener: vi.fn((_type: string, fn: (e: unknown) => void) => {
        clickHandler = fn
      }),
    } as unknown as ParentNode

    attachCategoryListEditDelegation(listRoot, onEditClick)
    clickHandler!(stubClickTargetEvent(null))

    expect(onEditClick).not.toHaveBeenCalled()
  })
})
