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
})
