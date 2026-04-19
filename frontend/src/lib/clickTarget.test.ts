import { describe, expect, it } from 'vitest'

import { clickEventTargetElement } from './clickTarget'

describe('clickEventTargetElement', () => {
  it('returns an Element target unchanged', () => {
    const div = {
      closest: () => null,
    } as unknown as Element
    expect(clickEventTargetElement({ target: div } as unknown as MouseEvent)).toBe(div)
  })

  it('returns parentElement when target has no closest (e.g. text-like node)', () => {
    const parent = {
      closest: (sel: string) => (sel === '.btn' ? parent : null),
    } as unknown as Element
    const text = { parentElement: parent }
    expect(clickEventTargetElement({ target: text } as unknown as MouseEvent)).toBe(parent)
  })

  it('returns null for null target', () => {
    expect(clickEventTargetElement({ target: null } as unknown as MouseEvent)).toBeNull()
  })

  it('returns null when target is undefined', () => {
    expect(clickEventTargetElement({ target: undefined } as unknown as MouseEvent)).toBeNull()
  })

  it('returns null for non-object targets', () => {
    expect(clickEventTargetElement({ target: 0 } as unknown as MouseEvent)).toBeNull()
  })

  it('returns null when object has neither closest nor parentElement', () => {
    expect(clickEventTargetElement({ target: {} } as unknown as MouseEvent)).toBeNull()
  })
})
