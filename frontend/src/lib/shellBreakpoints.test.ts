import { describe, expect, it, vi } from 'vitest'

import { onMediaQueryChange } from './shellBreakpoints'

describe('onMediaQueryChange', () => {
  it('uses addEventListener/removeEventListener when available', () => {
    const fn = vi.fn()
    const mq = {
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    } as unknown as MediaQueryList

    const unsub = onMediaQueryChange(mq, fn as (this: MediaQueryList, ev: MediaQueryListEvent) => void)
    expect(mq.addEventListener).toHaveBeenCalledWith('change', fn)

    unsub()
    expect(mq.removeEventListener).toHaveBeenCalledWith('change', fn)
  })

  it('falls back to addListener/removeListener when addEventListener is absent', () => {
    const fn = vi.fn()
    const mq = {
      addEventListener: undefined,
      addListener: vi.fn(),
      removeListener: vi.fn(),
    } as unknown as MediaQueryList

    const unsub = onMediaQueryChange(mq, fn as (this: MediaQueryList, ev: MediaQueryListEvent) => void)
    expect(mq.addListener).toHaveBeenCalledWith(fn)

    unsub()
    expect(mq.removeListener).toHaveBeenCalledWith(fn)
  })
})
