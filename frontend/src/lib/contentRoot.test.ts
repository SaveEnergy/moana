import { afterEach, describe, expect, it, vi } from 'vitest'

import { APP_MAIN_SELECTOR } from './domSelectors'
import { resolveContentQueryRoot } from './contentRoot'

describe('resolveContentQueryRoot', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns parent when it is not the document', () => {
    const parent = { _: 'subtree' } as unknown as ParentNode
    vi.stubGlobal('document', {})
    expect(resolveContentQueryRoot(parent)).toBe(parent)
  })

  it('returns main when parent is document and main exists', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const doc = {
      querySelector: (sel: string) => (sel === APP_MAIN_SELECTOR ? main : null),
    } as unknown as Document
    vi.stubGlobal('document', doc)
    expect(resolveContentQueryRoot(doc as unknown as ParentNode)).toBe(main)
  })

  it('returns document when main is absent', () => {
    const doc = {
      querySelector: () => null,
    } as unknown as Document
    vi.stubGlobal('document', doc)
    expect(resolveContentQueryRoot(doc as unknown as ParentNode)).toBe(doc)
  })

  it('queries the main landmark at most once per document instance', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const qs = vi.fn((sel: string) => (sel === APP_MAIN_SELECTOR ? main : null))
    const doc = { querySelector: qs } as unknown as Document
    vi.stubGlobal('document', doc)
    const p = doc as unknown as ParentNode
    expect(resolveContentQueryRoot(p)).toBe(main)
    expect(resolveContentQueryRoot(p)).toBe(main)
    expect(qs).toHaveBeenCalledTimes(1)
    expect(qs).toHaveBeenCalledWith(APP_MAIN_SELECTOR)
  })
})
