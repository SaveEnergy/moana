import { afterEach, describe, expect, it, vi } from 'vitest'

import { runBootInitializers } from '../boot'
import { APP_MAIN_SELECTOR, FORM_DATA_CONFIRM_SELECTOR, HISTORY_SORT_SELECTOR } from './domSelectors'
import {
  queryBootContent,
  queryBootContentAll,
  resolveBootContentQueryRoot,
  resolveContentQueryRoot,
  withBootContentRootScope,
} from './contentRoot'
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

  it('matches resolveBootContentQueryRoot().querySelector (API equivalence)', () => {
    const select = { _: 'history-sort' } as unknown as HTMLSelectElement
    const main = {
      querySelector: (sel: string) => (sel === HISTORY_SORT_SELECTOR ? select : null),
    } as unknown as ParentNode
    vi.stubGlobal('document', stubDocumentMainLandmark(main))
    const root = resolveBootContentQueryRoot()
    expect(queryBootContent<HTMLSelectElement>(HISTORY_SORT_SELECTOR)).toBe(
      root.querySelector<HTMLSelectElement>(HISTORY_SORT_SELECTOR),
    )
  })
})

describe('queryBootContentAll', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('querySelectorAll on the resolved boot content root', () => {
    const f1 = { _: 'f1' } as unknown as HTMLFormElement
    const f2 = { _: 'f2' } as unknown as HTMLFormElement
    const list = [f1, f2] as unknown as NodeListOf<HTMLFormElement>
    const main = {
      querySelectorAll: (sel: string) => (sel === FORM_DATA_CONFIRM_SELECTOR ? list : ([] as unknown as NodeListOf<HTMLFormElement>)),
    } as unknown as ParentNode
    vi.stubGlobal('document', stubDocumentMainLandmark(main))
    const out = queryBootContentAll<HTMLFormElement>(FORM_DATA_CONFIRM_SELECTOR)
    expect(out).toBe(list)
    expect(out.length).toBe(2)
    expect(out[0]).toBe(f1)
    expect(out[1]).toBe(f2)
  })

  it('matches resolveBootContentQueryRoot().querySelectorAll (API equivalence)', () => {
    const f1 = { _: 'f1' } as unknown as HTMLFormElement
    const list = [f1] as unknown as NodeListOf<HTMLFormElement>
    const main = {
      querySelectorAll: (sel: string) => (sel === FORM_DATA_CONFIRM_SELECTOR ? list : ([] as unknown as NodeListOf<HTMLFormElement>)),
    } as unknown as ParentNode
    vi.stubGlobal('document', stubDocumentMainLandmark(main))
    const root = resolveBootContentQueryRoot()
    expect(queryBootContentAll<HTMLFormElement>(FORM_DATA_CONFIRM_SELECTOR)).toBe(
      root.querySelectorAll<HTMLFormElement>(FORM_DATA_CONFIRM_SELECTOR),
    )
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

describe('withBootContentRootScope', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('keeps a stable boot content root for nested resolveBootContentQueryRoot during the callback', () => {
    const main = { _: 'main' } as unknown as ParentNode
    vi.stubGlobal('document', stubDocumentMainLandmark(main))
    withBootContentRootScope(() => {
      const a = resolveBootContentQueryRoot()
      const b = resolveBootContentQueryRoot()
      expect(a).toBe(main)
      expect(b).toBe(main)
      expect(a).toBe(b)
    })
  })

  it('restores the outer pin after nested scopes (different document globals)', () => {
    const mainOuter = { _: 'outer' } as unknown as ParentNode
    const mainInner = { _: 'inner' } as unknown as ParentNode
    const docOuter = stubDocumentMainLandmark(mainOuter) as unknown as Document
    const docInner = stubDocumentMainLandmark(mainInner) as unknown as Document
    vi.stubGlobal('document', docOuter)
    withBootContentRootScope(() => {
      expect(resolveBootContentQueryRoot()).toBe(mainOuter)
      vi.stubGlobal('document', docInner)
      withBootContentRootScope(() => {
        expect(resolveBootContentQueryRoot()).toBe(mainInner)
      })
      vi.stubGlobal('document', docOuter)
      expect(resolveBootContentQueryRoot()).toBe(mainOuter)
    })
  })

  it('runBootInitializers exposes the same boot content root to each step', () => {
    const main = { _: 'main' } as unknown as ParentNode
    vi.stubGlobal('document', stubDocumentMainLandmark(main))
    runBootInitializers([
      () => {
        expect(resolveBootContentQueryRoot()).toBe(main)
      },
      () => {
        expect(resolveBootContentQueryRoot()).toBe(main)
      },
    ])
  })

  it('runBootInitializers([]) is a no-op', () => {
    expect(() => runBootInitializers([])).not.toThrow()
  })
})
