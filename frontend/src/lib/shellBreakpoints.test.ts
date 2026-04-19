import { describe, expect, it, vi } from 'vitest'

import { MOBILE_SHELL_MAX_WIDTH_PX, MOBILE_SHELL_MEDIA_QUERY, onMediaQueryChange } from './shellBreakpoints'

describe('mobile shell breakpoint constants', () => {
  it('keeps media query aligned with max-width px (shellSidebar + CSS breakpoints)', () => {
    expect(MOBILE_SHELL_MEDIA_QUERY).toBe(`(max-width: ${MOBILE_SHELL_MAX_WIDTH_PX}px)`)
    expect(MOBILE_SHELL_MAX_WIDTH_PX).toBe(1023)
  })
})

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
