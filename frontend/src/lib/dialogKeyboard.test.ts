import { describe, expect, it } from 'vitest'

import {
  eventPathIncludesOpenDetails,
  eventPathIncludesOpenDialog,
  keyEventInvolvesOpenDialog,
} from './dialogKeyboard'

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

describe('eventPathIncludesOpenDetails', () => {
  it('detects open DETAILS', () => {
    expect(eventPathIncludesOpenDetails([{ tagName: 'DETAILS', open: true }])).toBe(true)
    expect(eventPathIncludesOpenDetails([{ tagName: 'DETAILS', open: false }])).toBe(false)
  })
})

describe('keyEventInvolvesOpenDialog', () => {
  it('delegates to composedPath (dialog only)', () => {
    const e = {
      composedPath: () => [{ tagName: 'DIALOG', open: true }],
    } as unknown as KeyboardEvent
    expect(keyEventInvolvesOpenDialog(e)).toBe(true)
  })

  it('is false for details-only path (use isAppUserMenuDetailsOpen in shell)', () => {
    const e = {
      composedPath: () => [{ tagName: 'DETAILS', open: true }],
    } as unknown as KeyboardEvent
    expect(keyEventInvolvesOpenDialog(e)).toBe(false)
  })
})
