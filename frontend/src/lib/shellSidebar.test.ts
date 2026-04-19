import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  APP_SHELL_SELECTOR,
  APP_SIDEBAR_BACKDROP_SELECTOR,
  APP_SIDEBAR_TOGGLE_SELECTOR,
} from './domSelectors'
import {
  initShellSidebar,
  queryAppShell,
  querySidebarBackdrop,
  querySidebarToggle,
} from './shellSidebar'
import { MOBILE_SHELL_MEDIA_QUERY } from './shellBreakpoints'

describe('queryAppShell', () => {
  it('uses APP_SHELL_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelector: (sel: string) => {
        seen = sel
        return null
      },
    } as unknown as ParentNode
    queryAppShell(root)
    expect(seen).toBe(APP_SHELL_SELECTOR)
  })
})

describe('querySidebarToggle', () => {
  it('uses APP_SIDEBAR_TOGGLE_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelector: (sel: string) => {
        seen = sel
        return null
      },
    } as unknown as ParentNode
    querySidebarToggle(root)
    expect(seen).toBe(APP_SIDEBAR_TOGGLE_SELECTOR)
  })
})

describe('querySidebarBackdrop', () => {
  it('uses APP_SIDEBAR_BACKDROP_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelector: (sel: string) => {
        seen = sel
        return null
      },
    } as unknown as ParentNode
    querySidebarBackdrop(root)
    expect(seen).toBe(APP_SIDEBAR_BACKDROP_SELECTOR)
  })
})

describe('initShellSidebar', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('returns early when #app-shell is missing', () => {
    const querySelector = vi.fn(() => null)
    const addEventListener = vi.fn()
    vi.stubGlobal('document', { querySelector, addEventListener })
    vi.stubGlobal('window', { matchMedia: vi.fn() })

    initShellSidebar()

    expect(querySelector).toHaveBeenCalledWith(APP_SHELL_SELECTOR)
    expect(addEventListener).not.toHaveBeenCalled()
  })

  it('wires matchMedia, toggle + shell click, media change, and capture keydown', () => {
    const matchMedia = vi.fn()
    const mq = {
      matches: true,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }
    matchMedia.mockReturnValue(mq)

    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains: vi.fn(() => false),
      },
      addEventListener: vi.fn(),
    }
    const toggle = { addEventListener: vi.fn(), setAttribute: vi.fn() }
    const backdrop = { setAttribute: vi.fn() }

    const querySelector = vi.fn((sel: string) => {
      if (sel === APP_SHELL_SELECTOR) return appShell
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return toggle
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })

    const doc = {
      querySelector,
      addEventListener: vi.fn(),
    }

    vi.stubGlobal('document', doc)
    vi.stubGlobal('window', { matchMedia })

    initShellSidebar()

    expect(matchMedia).toHaveBeenCalledWith(MOBILE_SHELL_MEDIA_QUERY)
    expect(toggle.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    expect(appShell.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    expect(mq.addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    expect(doc.addEventListener).toHaveBeenCalledWith('keydown', expect.any(Function), { capture: true })
  })

  it('does not stack listeners when initShellSidebar runs twice', () => {
    const matchMedia = vi.fn()
    const mq = {
      matches: true,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }
    matchMedia.mockReturnValue(mq)

    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains: vi.fn(() => false),
      },
      addEventListener: vi.fn(),
    }
    const toggle = { addEventListener: vi.fn(), setAttribute: vi.fn() }
    const backdrop = { setAttribute: vi.fn() }

    const querySelector = vi.fn((sel: string) => {
      if (sel === APP_SHELL_SELECTOR) return appShell
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return toggle
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })

    const doc = {
      querySelector,
      addEventListener: vi.fn(),
    }

    vi.stubGlobal('document', doc)
    vi.stubGlobal('window', { matchMedia })

    initShellSidebar()
    initShellSidebar()

    expect(matchMedia).toHaveBeenCalledTimes(1)
    expect(toggle.addEventListener).toHaveBeenCalledTimes(1)
    expect(appShell.addEventListener).toHaveBeenCalledTimes(1)
    expect(mq.addEventListener).toHaveBeenCalledTimes(1)
    expect(doc.addEventListener).toHaveBeenCalledTimes(1)
  })
})
