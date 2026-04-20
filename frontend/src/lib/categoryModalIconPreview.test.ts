import { describe, expect, it } from 'vitest'

import { shouldRepaintCategoryModalIconPreview } from './categoryModalIconPreview'

describe('shouldRepaintCategoryModalIconPreview', () => {
  it('requires paint on first run (undefined vs any current value)', () => {
    expect(shouldRepaintCategoryModalIconPreview(undefined, '')).toBe(true)
    expect(shouldRepaintCategoryModalIconPreview(undefined, 'wallet')).toBe(true)
  })

  it('skips paint when the icon group value is unchanged', () => {
    expect(shouldRepaintCategoryModalIconPreview('', '')).toBe(false)
    expect(shouldRepaintCategoryModalIconPreview('piggy-bank', 'piggy-bank')).toBe(false)
  })

  it('requires paint when the icon group value changes', () => {
    expect(shouldRepaintCategoryModalIconPreview('', 'wallet')).toBe(true)
    expect(shouldRepaintCategoryModalIconPreview('wallet', '')).toBe(true)
    expect(shouldRepaintCategoryModalIconPreview('wallet', 'piggy-bank')).toBe(true)
  })
})
