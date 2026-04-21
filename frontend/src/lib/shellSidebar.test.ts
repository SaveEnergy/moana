import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  APP_SHELL_SELECTOR,
  APP_SIDEBAR_BACKDROP_SELECTOR,
  APP_SIDEBAR_TOGGLE_SELECTOR,
} from './domSelectors'
import * as dialogKeyboard from './dialogKeyboard'
import * as mobileShellDismiss from './mobileShellDismiss'
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

    const toggle = { addEventListener: vi.fn(), setAttribute: vi.fn() }
    const backdrop = { setAttribute: vi.fn() }

    const shellQuerySelector = vi.fn((sel: string) => {
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return toggle
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })
    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains: vi.fn(() => false),
      },
      addEventListener: vi.fn(),
      querySelector: shellQuerySelector,
    }

    const querySelector = vi.fn((sel: string) => {
      if (sel === APP_SHELL_SELECTOR) return appShell
      return null
    })

    const doc = {
      querySelector,
      addEventListener: vi.fn(),
    }

    vi.stubGlobal('document', doc)
    vi.stubGlobal('window', { matchMedia })

    initShellSidebar()

    expect(shellQuerySelector).toHaveBeenCalledWith(APP_SIDEBAR_TOGGLE_SELECTOR)
    expect(shellQuerySelector).toHaveBeenCalledWith(APP_SIDEBAR_BACKDROP_SELECTOR)
    expect(shellQuerySelector).toHaveBeenCalledTimes(2)
    expect(matchMedia).toHaveBeenCalledWith(MOBILE_SHELL_MEDIA_QUERY)
    expect(toggle.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    expect(appShell.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    expect(mq.addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    expect(doc.addEventListener).toHaveBeenCalledWith('keydown', expect.any(Function), { capture: true })
  })

  it('wires shell click + capture keydown when toggle is absent but backdrop exists', () => {
    const matchMedia = vi.fn()
    const mq = {
      matches: true,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }
    matchMedia.mockReturnValue(mq)

    const backdrop = { setAttribute: vi.fn() }

    const shellQuerySelector = vi.fn((sel: string) => {
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return null
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })
    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains: vi.fn(() => false),
      },
      addEventListener: vi.fn(),
      querySelector: shellQuerySelector,
    }

    const querySelector = vi.fn((sel: string) => {
      if (sel === APP_SHELL_SELECTOR) return appShell
      return null
    })

    const doc = {
      querySelector,
      addEventListener: vi.fn(),
    }

    vi.stubGlobal('document', doc)
    vi.stubGlobal('window', { matchMedia })

    expect(() => initShellSidebar()).not.toThrow()

    expect(shellQuerySelector).toHaveBeenCalledWith(APP_SIDEBAR_TOGGLE_SELECTOR)
    expect(shellQuerySelector).toHaveBeenCalledWith(APP_SIDEBAR_BACKDROP_SELECTOR)
    expect(matchMedia).toHaveBeenCalledWith(MOBILE_SHELL_MEDIA_QUERY)
    expect(appShell.addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    expect(mq.addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    expect(doc.addEventListener).toHaveBeenCalledWith('keydown', expect.any(Function), { capture: true })
  })

  it('skips redundant backdrop and toggle ARIA writes when already closed', () => {
    vi.spyOn(dialogKeyboard, 'shouldDeferMobileShellEscape').mockReturnValue(false)

    const matchMedia = vi.fn()
    const mq = {
      matches: true,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }
    matchMedia.mockReturnValue(mq)

    const toggleAttr: Record<string, string> = {
      'aria-expanded': 'false',
      'aria-label': 'Open navigation menu',
    }
    const toggle = {
      addEventListener: vi.fn(),
      getAttribute: vi.fn((n: string) => toggleAttr[n] ?? null),
      setAttribute: vi.fn((n: string, v: string) => {
        toggleAttr[n] = v
      }),
    }

    let ariaHidden = 'true'
    const backdrop = {
      getAttribute: vi.fn((n: string) => (n === 'aria-hidden' ? ariaHidden : null)),
      setAttribute: vi.fn((n: string, v: string) => {
        if (n === 'aria-hidden') ariaHidden = v
      }),
    }

    const shellQuerySelector = vi.fn((sel: string) => {
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return toggle
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })
    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains: vi.fn(() => false),
      },
      addEventListener: vi.fn(),
      querySelector: shellQuerySelector,
    }

    const doc = {
      querySelector: vi.fn((sel: string) => (sel === APP_SHELL_SELECTOR ? appShell : null)),
      addEventListener: vi.fn(),
    }

    vi.stubGlobal('document', doc)
    vi.stubGlobal('window', { matchMedia })

    initShellSidebar()

    const keydown = doc.addEventListener.mock.calls.find((c) => c[0] === 'keydown')?.[1] as (
      e: KeyboardEvent,
    ) => void
    expect(keydown).toBeDefined()

    backdrop.setAttribute.mockClear()
    toggle.setAttribute.mockClear()

    keydown!({ key: 'Escape', composedPath: () => [] } as unknown as KeyboardEvent)

    expect(backdrop.setAttribute).not.toHaveBeenCalled()
    expect(toggle.setAttribute).not.toHaveBeenCalled()
  })

  it('shell click does not run dismiss predicate when drawer is closed', () => {
    const dismissSpy = vi.spyOn(mobileShellDismiss, 'shouldCloseMobileSidebarFromShellClick')

    const matchMedia = vi.fn()
    const mq = {
      matches: true,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }
    matchMedia.mockReturnValue(mq)

    const toggle = { addEventListener: vi.fn(), setAttribute: vi.fn() }
    const backdrop = { setAttribute: vi.fn() }

    const shellQuerySelector = vi.fn((sel: string) => {
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return toggle
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })

    const contains = vi.fn().mockReturnValue(false)
    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains,
      },
      addEventListener: vi.fn(),
      querySelector: shellQuerySelector,
    }

    const doc = {
      querySelector: vi.fn((sel: string) => (sel === APP_SHELL_SELECTOR ? appShell : null)),
      addEventListener: vi.fn(),
    }

    vi.stubGlobal('document', doc)
    vi.stubGlobal('window', { matchMedia })

    initShellSidebar()

    const shellClick = appShell.addEventListener.mock.calls.find((c) => c[0] === 'click')?.[1] as (
      e: unknown,
    ) => void
    expect(shellClick).toBeDefined()

    dismissSpy.mockClear()
    shellClick!({} as MouseEvent)
    expect(dismissSpy).not.toHaveBeenCalled()

    contains.mockReturnValue(true)
    shellClick!({} as MouseEvent)
    expect(dismissSpy).toHaveBeenCalledTimes(1)

    dismissSpy.mockRestore()
  })

  it('does not stack listeners when initShellSidebar runs twice', () => {
    const matchMedia = vi.fn()
    const mq = {
      matches: true,
      addEventListener: vi.fn(),
      removeEventListener: vi.fn(),
    }
    matchMedia.mockReturnValue(mq)

    const toggle = { addEventListener: vi.fn(), setAttribute: vi.fn() }
    const backdrop = { setAttribute: vi.fn() }

    const shellQuerySelector = vi.fn((sel: string) => {
      if (sel === APP_SIDEBAR_TOGGLE_SELECTOR) return toggle
      if (sel === APP_SIDEBAR_BACKDROP_SELECTOR) return backdrop
      return null
    })
    const appShell = {
      classList: {
        add: vi.fn(),
        remove: vi.fn(),
        contains: vi.fn(() => false),
      },
      addEventListener: vi.fn(),
      querySelector: shellQuerySelector,
    }

    const querySelector = vi.fn((sel: string) => {
      if (sel === APP_SHELL_SELECTOR) return appShell
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

    expect(shellQuerySelector).toHaveBeenCalledTimes(2)
    expect(matchMedia).toHaveBeenCalledTimes(1)
    expect(toggle.addEventListener).toHaveBeenCalledTimes(1)
    expect(appShell.addEventListener).toHaveBeenCalledTimes(1)
    expect(mq.addEventListener).toHaveBeenCalledTimes(1)
    expect(doc.addEventListener).toHaveBeenCalledTimes(1)
  })
})
