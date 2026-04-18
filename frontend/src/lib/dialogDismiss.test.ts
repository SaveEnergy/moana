import { describe, expect, it, vi } from 'vitest'

import { shouldCloseNativeDialogFromClick } from './dialogDismiss'

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
