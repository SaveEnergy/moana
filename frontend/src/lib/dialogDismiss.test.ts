import { describe, expect, it, vi } from 'vitest'

import { attachNativeDialogDismiss, shouldCloseNativeDialogFromClick } from './dialogDismiss'

describe('shouldCloseNativeDialogFromClick', () => {
  it('closes on backdrop (target === dialog)', () => {
    const close = vi.fn()
    // Plain object is not `instanceof Element` in Node; `closest` satisfies `clickEventTargetElement`.
    const dialog = { close, closest: () => null } as unknown as HTMLDialogElement
    const e = { target: dialog } as unknown as MouseEvent
    expect(shouldCloseNativeDialogFromClick(e, dialog, ['#x'])).toBe(true)
  })

  it('closes when closest matches a selector', () => {
    const close = vi.fn()
    const dialog = { close } as unknown as HTMLDialogElement
    const inner = {
      closest: (sel: string) => (sel === '#close' ? inner : null),
    } as unknown as Element
    expect(
      shouldCloseNativeDialogFromClick({ target: inner } as unknown as MouseEvent, dialog, ['#close']),
    ).toBe(true)
  })

  it('does not close when click is inside dialog but not on a close control', () => {
    const dialog = {} as unknown as HTMLDialogElement
    const content = {
      closest: () => null,
    } as unknown as Element
    expect(
      shouldCloseNativeDialogFromClick({ target: content } as unknown as MouseEvent, dialog, ['#close']),
    ).toBe(false)
  })
})

describe('attachNativeDialogDismiss', () => {
  it('closes the dialog on backdrop click (target === dialog)', () => {
    const close = vi.fn()
    const addEventListener = vi.fn()
    /** `closest` so `clickEventTargetElement` resolves in Node (no `instanceof Element`). */
    const dialog = { close, addEventListener, closest: () => null } as unknown as HTMLDialogElement

    attachNativeDialogDismiss(dialog, ['#ignored-for-backdrop'])

    const onClick = addEventListener.mock.calls[0][1] as (e: MouseEvent) => void
    onClick({ target: dialog } as unknown as MouseEvent)

    expect(close).toHaveBeenCalledTimes(1)
  })

  it('does not close when inner target matches no dismiss selector', () => {
    const close = vi.fn()
    const addEventListener = vi.fn()
    const dialog = { close, addEventListener } as unknown as HTMLDialogElement
    const inner = { closest: () => null } as unknown as Element

    attachNativeDialogDismiss(dialog, ['#cat-modal-close'])

    const onClick = addEventListener.mock.calls[0][1] as (e: MouseEvent) => void
    onClick({ target: inner } as unknown as MouseEvent)

    expect(close).not.toHaveBeenCalled()
  })
})
