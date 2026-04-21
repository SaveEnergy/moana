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
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIV', open: true }])).toBe(false)
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIALOG', open: false }])).toBe(false)
    expect(eventPathIncludesOpenDialog([{}])).toBe(false)
  })

  it('is true when an open DIALOG is in the path', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'DIALOG', open: true }])).toBe(true)
  })

  it('matches open dialog when tagName is lowercase', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'dialog', open: true }])).toBe(true)
  })

  it('matches open dialog when tagName is mixed case (HTMLDocument tagName normalization edge)', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'Dialog', open: true }])).toBe(true)
  })

  it('ignores nodes whose tagName is not a string', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 1, open: true }] as unknown[])).toBe(false)
  })

  it('is false when another element reports open: true (e.g. VIDEO)', () => {
    expect(eventPathIncludesOpenDialog([{ tagName: 'VIDEO', open: true }])).toBe(false)
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

  it('defers when details.app-user-menu[open] exists (DOM query after path scan)', () => {
    const qs = vi.fn((sel: string) => (sel === APP_USER_MENU_OPEN_SELECTOR ? ({} as Element) : null))
    vi.stubGlobal('document', { querySelector: qs })
    const e = { composedPath: () => [] } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(true)
    expect(qs).toHaveBeenCalledTimes(1)
    expect(qs).toHaveBeenCalledWith(APP_USER_MENU_OPEN_SELECTOR)
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

  it('is false when neither applies after one menu-selector query', () => {
    const qs = vi.fn(() => null)
    vi.stubGlobal('document', { querySelector: qs })
    const e = { composedPath: () => [] } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(false)
    expect(qs).toHaveBeenCalledTimes(1)
    expect(qs).toHaveBeenCalledWith(APP_USER_MENU_OPEN_SELECTOR)
  })

  it('queries the menu selector when path has nodes but none defer', () => {
    const qs = vi.fn(() => null)
    vi.stubGlobal('document', { querySelector: qs })
    const e = {
      composedPath: () => [{ tagName: 'DIV' }, { tagName: 'BUTTON' }],
    } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(false)
    expect(qs).toHaveBeenCalledTimes(1)
    expect(qs).toHaveBeenCalledWith(APP_USER_MENU_OPEN_SELECTOR)
  })

  it('does not defer when only a non-dialog node reports open: true (falls through to menu query)', () => {
    const qs = vi.fn(() => null)
    vi.stubGlobal('document', { querySelector: qs })
    const e = {
      composedPath: () => [{ tagName: 'DIV', open: true }],
    } as unknown as KeyboardEvent
    expect(shouldDeferMobileShellEscape(e)).toBe(false)
    expect(qs).toHaveBeenCalledTimes(1)
  })
})
