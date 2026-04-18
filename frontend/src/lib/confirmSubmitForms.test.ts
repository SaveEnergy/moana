import { afterEach, describe, expect, it, vi } from 'vitest'

import { attachConfirmBeforeSubmit } from './confirmSubmitForms'

describe('attachConfirmBeforeSubmit', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('calls preventDefault when confirm returns false', () => {
    vi.stubGlobal('confirm', vi.fn(() => false))
    let onSubmit: ((e: SubmitEvent) => void) | null = null
    const form = {
      addEventListener: vi.fn((ev: string, fn: (e: SubmitEvent) => void) => {
        if (ev === 'submit') {
          onSubmit = fn
        }
      }),
    } as unknown as HTMLFormElement
    attachConfirmBeforeSubmit(form, 'Remove?')
    expect(onSubmit).not.toBeNull()
    const e = { preventDefault: vi.fn() } as unknown as SubmitEvent
    onSubmit!(e)
    expect(e.preventDefault).toHaveBeenCalledTimes(1)
  })

  it('does not prevent default when confirm returns true', () => {
    vi.stubGlobal('confirm', vi.fn(() => true))
    let onSubmit: ((e: SubmitEvent) => void) | null = null
    const form = {
      addEventListener: vi.fn((ev: string, fn: (e: SubmitEvent) => void) => {
        if (ev === 'submit') {
          onSubmit = fn
        }
      }),
    } as unknown as HTMLFormElement
    attachConfirmBeforeSubmit(form, 'Remove?')
    const e = { preventDefault: vi.fn() } as unknown as SubmitEvent
    onSubmit!(e)
    expect(e.preventDefault).not.toHaveBeenCalled()
  })
})
