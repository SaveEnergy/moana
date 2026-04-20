import { describe, expect, it } from 'vitest'

import { type ClickTargetEvent, clickEventTargetElement } from './clickTarget'

/** Vitest stubs are not real `EventTarget` instances — cast at the test boundary only. */
function stubEv(target: unknown): ClickTargetEvent {
  return { target } as unknown as ClickTargetEvent
}

describe('clickEventTargetElement', () => {
  it('returns an Element target unchanged', () => {
    const div = {
      closest: () => null,
    } as unknown as Element
    expect(clickEventTargetElement(stubEv(div))).toBe(div)
  })

  it('returns the target when it has closest() but is not an Element (legacy / partial DOM)', () => {
    const fake = { closest: () => null as Element | null }
    expect(clickEventTargetElement(stubEv(fake))).toBe(fake)
  })

  it('returns parentElement when target has no closest (e.g. text-like node)', () => {
    const parent = {
      closest: (sel: string) => (sel === '.btn' ? parent : null),
    } as unknown as Element
    const text = { parentElement: parent }
    expect(clickEventTargetElement(stubEv(text))).toBe(parent)
  })

  it('returns null when text-like node has no parentElement', () => {
    const text = { parentElement: null }
    expect(clickEventTargetElement(stubEv(text))).toBeNull()
  })

  it('returns null for null target', () => {
    expect(clickEventTargetElement(stubEv(null))).toBeNull()
  })

  it('returns null when target is undefined', () => {
    expect(clickEventTargetElement(stubEv(undefined))).toBeNull()
  })

  it('returns null for non-object targets', () => {
    expect(clickEventTargetElement(stubEv(0))).toBeNull()
  })

  it('returns null when object has neither closest nor parentElement', () => {
    expect(clickEventTargetElement(stubEv({}))).toBeNull()
  })
})
