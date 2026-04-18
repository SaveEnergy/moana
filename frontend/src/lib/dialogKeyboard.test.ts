import { describe, expect, it } from 'vitest'

import { eventPathIncludesOpenDialog, keyEventInvolvesOpenDialog } from './dialogKeyboard'

describe('eventPathIncludesOpenDialog', () => {
  it('is false for empty or irrelevant paths', () => {
    expect(eventPathIncludesOpenDialog([])).toBe(false)
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIV' }])).toBe(false)
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIALOG', open: false }])).toBe(false)
  })

  it('is true when an open DIALOG is in the path', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIALOG', open: true }])).toBe(true)
  })
})

describe('keyEventInvolvesOpenDialog', () => {
  it('delegates to composedPath', () => {
    const e = {
      composedPath: () => [{ tagName: 'DIALOG', open: true }],
    } as unknown as KeyboardEvent
    expect(keyEventInvolvesOpenDialog(e)).toBe(true)
  })
})
