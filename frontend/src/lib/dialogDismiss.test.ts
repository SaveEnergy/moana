import { describe, expect, it, vi } from 'vitest'

import { type ClickTargetEvent, stubClickTargetEvent } from './clickTarget'
import { attachNativeDialogDismiss, shouldCloseNativeDialogFromClick } from './dialogDismiss'

describe('shouldCloseNativeDialogFromClick', () => {
  it('closes on backdrop (target === dialog)', () => {
    const close = vi.fn()
    // Plain object is not `instanceof Element` in Node; `closest` satisfies `clickEventTargetElement`.
    const dialog = { close, closest: () => null } as unknown as HTMLDialogElement
    expect(shouldCloseNativeDialogFromClick(stubClickTargetEvent(dialog), dialog, ['#x'])).toBe(true)
  })

  it('closes when closest matches a selector', () => {
    const close = vi.fn()
    const dialog = { close } as unknown as HTMLDialogElement
    const inner = {
      closest: (sel: string) => (sel === '#close' ? inner : null),
    } as unknown as Element
    expect(shouldCloseNativeDialogFromClick(stubClickTargetEvent(inner), dialog, ['#close'])).toBe(true)
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

  it('still closes on backdrop when dismiss selector list is empty', () => {
    const dialog = { closest: () => null } as unknown as HTMLDialogElement
    expect(shouldCloseNativeDialogFromClick(stubClickTargetEvent(dialog), dialog, [])).toBe(true)
  })

  it('does not close on inner content when dismiss selector list is empty', () => {
    const dialog = {} as unknown as HTMLDialogElement
    const inner = {
      closest: () => null,
    } as unknown as Element
    expect(shouldCloseNativeDialogFromClick(stubClickTargetEvent(inner), dialog, [])).toBe(false)
  })

  it('ignores blank dismiss selectors and still matches valid ones', () => {
    const dialog = {} as unknown as HTMLDialogElement
    const inner = {
      closest: (sel: string) => (sel === '#close' ? inner : null),
    } as unknown as Element
    expect(
      shouldCloseNativeDialogFromClick(stubClickTargetEvent(inner), dialog, ['', '   ', '#close']),
    ).toBe(true)
  })
})

describe('attachNativeDialogDismiss', () => {
  it('closes the dialog on backdrop click (target === dialog)', () => {
    const close = vi.fn()
    const addEventListener = vi.fn()
    /** `closest` so `clickEventTargetElement` resolves in Node (no `instanceof Element`). */
    const dialog = { close, addEventListener, closest: () => null } as unknown as HTMLDialogElement

    attachNativeDialogDismiss(dialog, ['#ignored-for-backdrop'])

    const onClick = addEventListener.mock.calls[0][1] as (e: ClickTargetEvent) => void
    onClick(stubClickTargetEvent(dialog))

    expect(close).toHaveBeenCalledTimes(1)
  })

  it('does not close when inner target matches no dismiss selector', () => {
    const close = vi.fn()
    const addEventListener = vi.fn()
    const dialog = { close, addEventListener } as unknown as HTMLDialogElement
    const inner = { closest: () => null } as unknown as Element

    attachNativeDialogDismiss(dialog, ['#cat-modal-close'])

    const onClick = addEventListener.mock.calls[0][1] as (e: ClickTargetEvent) => void
    onClick(stubClickTargetEvent(inner))

    expect(close).not.toHaveBeenCalled()
  })

  it('does not register a second click listener when attach runs twice', () => {
    const addEventListener = vi.fn()
    const dialog = { close: vi.fn(), addEventListener, closest: () => null } as unknown as HTMLDialogElement

    attachNativeDialogDismiss(dialog, ['#x'])
    attachNativeDialogDismiss(dialog, ['#x'])

    expect(addEventListener).toHaveBeenCalledTimes(1)
  })
})
