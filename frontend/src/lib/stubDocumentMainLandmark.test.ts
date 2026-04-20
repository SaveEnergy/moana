import { describe, expect, it, vi } from 'vitest'

import { APP_MAIN_SELECTOR } from './domSelectors'
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'

describe('stubDocumentMainLandmark', () => {
  it('returns only the provided main for APP_MAIN_SELECTOR', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const doc = stubDocumentMainLandmark(main)
    expect(doc.querySelector(APP_MAIN_SELECTOR)).toBe(main)
    expect(doc.querySelector('#anything-else')).toBeNull()
  })
})

describe('stubDocumentWithoutMainLandmark', () => {
  it('defaults querySelector to null and can attach querySelectorAll', () => {
    const qsa = vi.fn((_: string) => [] as unknown as ReturnType<Document['querySelectorAll']>)
    const doc = stubDocumentWithoutMainLandmark({
      querySelectorAll: qsa as unknown as Document['querySelectorAll'],
    })
    expect(doc.querySelector(APP_MAIN_SELECTOR)).toBeNull()
    doc.querySelectorAll!('form[x]')
    expect(qsa).toHaveBeenCalledWith('form[x]')
  })

  it('uses a custom querySelector when only that override is provided', () => {
    const hit = { _: 'hit' } as unknown as Element
    const doc = stubDocumentWithoutMainLandmark({
      querySelector: (sel: string) => (sel === '#stub-hit' ? hit : null),
    })
    expect(doc.querySelector('#stub-hit')).toBe(hit)
    expect(doc.querySelector('#other')).toBeNull()
  })
})
