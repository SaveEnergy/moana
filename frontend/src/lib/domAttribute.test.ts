import { describe, expect, it, vi } from 'vitest'

import { setAttributeIfChanged } from './domAttribute'

describe('setAttributeIfChanged', () => {
  it('sets the attribute when missing or different', () => {
    const setAttribute = vi.fn()
    const el = {
      getAttribute: vi.fn(() => null),
      setAttribute,
    } as unknown as Element
    setAttributeIfChanged(el, 'aria-hidden', 'true')
    expect(setAttribute).toHaveBeenCalledWith('aria-hidden', 'true')
  })

  it('skips setAttribute when the value already matches', () => {
    const setAttribute = vi.fn()
    const el = {
      getAttribute: vi.fn(() => 'true'),
      setAttribute,
    } as unknown as Element
    setAttributeIfChanged(el, 'aria-hidden', 'true')
    expect(setAttribute).not.toHaveBeenCalled()
  })
})
