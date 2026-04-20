import { describe, expect, it } from 'vitest'

import { APP_MAIN_SELECTOR } from './domSelectors'
import { stubDocumentMainLandmark } from './stubDocumentMainLandmark'

describe('stubDocumentMainLandmark', () => {
  it('returns only the provided main for APP_MAIN_SELECTOR', () => {
    const main = { _: 'main' } as unknown as ParentNode
    const doc = stubDocumentMainLandmark(main)
    expect(doc.querySelector(APP_MAIN_SELECTOR)).toBe(main)
    expect(doc.querySelector('#anything-else')).toBeNull()
  })
})
