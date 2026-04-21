import { afterEach, describe, expect, it, vi } from 'vitest'

import { APP_MAIN_SELECTOR, HISTORY_SORT_SELECTOR } from './domSelectors'
import { queryBootContent, resolveBootContentQueryRoot, resolveContentQueryRoot } from './contentRoot'
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'

describe('resolveContentQueryRoot', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('returns parent when it is not the document', () => {
    const parent = { _: 'subtree' } as unknown as ParentNode
    vi.stubGlobal('document', {})
    expect(resolveContentQueryRoot(parent)).toBe(parent)
  })

  it('returns parent when the document global is undefined (Node / non-DOM)', () => {
    const parent = { _: 'no-dom' } as unknown as ParentNode
    vi.stubGlobal('document', undefined)
    expect(resolveContentQueryRoot(parent)).toBe(parent)
  })

  it('returns main when parent is document and main exists', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const doc = stubDocumentMainLandmark(main) as unknown as Document
    vi.stubGlobal('document', doc)
    expect(resolveContentQueryRoot(doc as unknown as ParentNode)).toBe(main)
  })

  it('returns document when main is absent', () => {
    const doc = stubDocumentWithoutMainLandmark() as unknown as Document
    vi.stubGlobal('document', doc)
    expect(resolveContentQueryRoot(doc as unknown as ParentNode)).toBe(doc)
  })

  it('queries the main landmark at most once per document instance', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const base = stubDocumentMainLandmark(main)
    const qs = vi.fn((sel: string) => base.querySelector(sel))
    const doc = { querySelector: qs } as unknown as Document
    vi.stubGlobal('document', doc)
    const p = doc as unknown as ParentNode
    expect(resolveContentQueryRoot(p)).toBe(main)
    expect(resolveContentQueryRoot(p)).toBe(main)
    expect(qs).toHaveBeenCalledTimes(1)
    expect(qs).toHaveBeenCalledWith(APP_MAIN_SELECTOR)
  })
})

describe('queryBootContent', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('querySelectors on the resolved boot content root', () => {
    const select = { _: 'history-sort' } as unknown as HTMLSelectElement
    const main = {
      querySelector: (sel: string) => (sel === HISTORY_SORT_SELECTOR ? select : null),
    } as unknown as ParentNode
    vi.stubGlobal('document', stubDocumentMainLandmark(main))
    expect(queryBootContent<HTMLSelectElement>(HISTORY_SORT_SELECTOR)).toBe(select)
  })
})

describe('resolveBootContentQueryRoot', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('delegates to resolveContentQueryRoot(document) with the same memoization', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const doc = stubDocumentMainLandmark(main) as unknown as Document
    vi.stubGlobal('document', doc)
    expect(resolveBootContentQueryRoot()).toBe(main)
    expect(resolveContentQueryRoot(doc as unknown as ParentNode)).toBe(main)
  })
})
