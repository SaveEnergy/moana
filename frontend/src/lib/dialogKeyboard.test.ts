import { afterEach, describe, expect, it, vi } from 'vitest'

import { APP_USER_MENU_OPEN_SELECTOR } from './domSelectors'
import {
  eventPathIncludesOpenDetails,
  eventPathIncludesOpenDialog,
  isAppUserMenuDetailsOpen,
  keyEventInvolvesOpenDialog,
  shouldDeferMobileShellEscape,
} from './dialogKeyboard'

describe('eventPathIncludesOpenDialog', () => {
  it('is false for empty or irrelevant paths', () => {
    expect(eventPathIncludesOpenDialog([])).toBe(false)
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIV' }])).toBe(false)
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIALOG', open: false }])).toBe(false)
    expect(eventPathIncludesOpenDialog([{}])).toBe(false)
  })

  it('is true when an open DIALOG is in the path', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIALOG', open: true }])).toBe(true)
  })

  it('matches open dialog when tagName is lowercase', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'dialog', open: true }])).toBe(true)
  })

  it('ignores nodes whose tagName is not a string', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 1, open: true }] as unknown[])).toBe(false)
  })

  it('detects open DIALOG when it is not the first path node', () => {
    expect(
      eventPathIncludesOpenDialog([
        { tagName: 'DIV' },
        { tagName: 'DIALOG', open: true },
      ]),
    ).toBe(true)
  })
})

describe('eventPathIncludesOpenDetails', () => {
  it('detects open DETAILS', () => {
    expect(eventPathIncludesOpenDetails([{ tagName: 'DETAILS', open: true }])).toBe(true)
    expect(eventPathIncludesOpenDetails([{ tagName: 'DETAILS', open: false }])).toBe(false)
  })

  it('matches open details when tagName is lowercase', () => {
    expect(eventPathIncludesOpenDetails([{ tagName: 'details', open: true }])).toBe(true)
  })

  it('detects open DETAILS when not the first path node', () => {
    expect(
      eventPathIncludesOpenDetails([
        { tagName: 'DIV' },
        { tagName: 'DETAILS', open: true },
      ]),
    ).toBe(true)
  })

  it('is false for empty paths or nodes without a usable tagName', () => {
    expect(eventPathIncludesOpenDetails([])).toBe(false)
    expect(eventPathIncludesOpenDetails([{}])).toBe(false)
  })
})

describe('isAppUserMenuDetailsOpen', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('is true when the open account menu selector matches', () => {
    vi.stubGlobal('document', {
      querySelector: (sel: string) => (sel === APP_USER_MENU_OPEN_SELECTOR ? ({} as Element) : null),
    })
    expect(isAppUserMenuDetailsOpen()).toBe(true)
  })

  it('is false when the menu is closed', () => {
    vi.stubGlobal('document', { querySelector: () => null })
    expect(isAppUserMenuDetailsOpen()).toBe(false)
  })
})

describe('keyEventInvolvesOpenDialog', () => {
  it('delegates to composedPath (dialog only)', () => {
    const e = {
      composedPath: () => [{ tagName: 'DIALOG', open: true }],
    } as unknown as KeyboardEvent
    expect(keyEventInvolvesOpenDialog(e)).toBe(true)
  })

  it('is false for details-only path (defer is handled in shouldDeferMobileShellEscape)', () => {
    const e = {
      composedPath: () => [{ tagName: 'DETAILS', open: true }],
    } as unknown as KeyboardEvent
    expect(keyEventInvolvesOpenDialog(e)).toBe(false)
  })
})

describe('shouldDeferMobileShellEscape', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('defers for open dialog in path without document.querySelector', () => {
    const qs = vi.fn(() => null)
    vi.stubGlobal('document', { querySelector: qs })
    const e = {
      composedPath: () => [{ tagName: 'DIALOG', open: true }],
    } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(true)
    expect(qs).not.toHaveBeenCalled()
  })

  it('defers when details.app-user-menu[open] exists', () => {
    vi.stubGlobal('document', {
      querySelector: (sel: string) => (sel === APP_USER_MENU_OPEN_SELECTOR ? ({} as Element) : null),
    })
    const e = { composedPath: () => [] } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(true)
  })

  it('defers when open DETAILS is in composedPath without document.querySelector', () => {
    const qs = vi.fn(() => null)
    vi.stubGlobal('document', { querySelector: qs })
    const e = {
      composedPath: () => [{ tagName: 'DETAILS', open: true }],
    } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(true)
    expect(qs).not.toHaveBeenCalled()
  })

  it('is false when neither applies', () => {
    vi.stubGlobal('document', { querySelector: () => null })
    const e = { composedPath: () => [] } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(false)
  })
})
