import { afterEach, describe, expect, it, vi } from 'vitest'

import { createRafScheduler } from './scheduleAnimationFrame'

describe('createRafScheduler', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('coalesces multiple schedule() calls into a single rAF callback', () => {
    const run = vi.fn()
    const callbacks: FrameRequestCallback[] = []
    vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
      callbacks.push(cb)
      return callbacks.length
    })
    vi.stubGlobal('cancelAnimationFrame', vi.fn())

    const { schedule } = createRafScheduler(run)
    schedule()
    schedule()
    schedule()

    expect(callbacks).toHaveLength(1)
    callbacks[0]!(0)
    expect(run).toHaveBeenCalledTimes(1)
  })

  it('cancelPending removes the frame without running', () => {
    const run = vi.fn()
    vi.stubGlobal('requestAnimationFrame', (_cb: FrameRequestCallback) => 42)
    const cancel = vi.fn()
    vi.stubGlobal('cancelAnimationFrame', cancel)

    const { schedule, cancelPending } = createRafScheduler(run)
    schedule()
    cancelPending()

    expect(cancel).toHaveBeenCalledWith(42)
    expect(run).not.toHaveBeenCalled()
  })

  it('allows scheduling again after the frame runs', () => {
    const run = vi.fn()
    const callbacks: FrameRequestCallback[] = []
    vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
      callbacks.push(cb)
      return callbacks.length
    })
    vi.stubGlobal('cancelAnimationFrame', vi.fn())

    const { schedule } = createRafScheduler(run)
    schedule()
    callbacks[0]!(0)
    schedule()
    expect(callbacks).toHaveLength(2)
    callbacks[1]!(0)
    expect(run).toHaveBeenCalledTimes(2)
  })

  it('allows scheduling again after cancelPending (new frame, no stale callback)', () => {
    const run = vi.fn()
    const callbacks: FrameRequestCallback[] = []
    vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
      callbacks.push(cb)
      return callbacks.length
    })
    vi.stubGlobal('cancelAnimationFrame', vi.fn())

    const { schedule, cancelPending } = createRafScheduler(run)
    schedule()
    cancelPending()
    schedule()
    expect(callbacks).toHaveLength(2)
    callbacks[1]!(0)
    expect(run).toHaveBeenCalledTimes(1)
  })

  it('runs synchronously when requestAnimationFrame is not a function', () => {
    vi.stubGlobal('requestAnimationFrame', undefined as unknown as typeof requestAnimationFrame)
    const run = vi.fn()
    const { schedule } = createRafScheduler(run)
    schedule()
    expect(run).toHaveBeenCalledTimes(1)
  })
})
