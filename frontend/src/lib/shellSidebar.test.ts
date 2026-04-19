import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  APP_SHELL_ELEMENT_ID,
  APP_SIDEBAR_BACKDROP_ELEMENT_ID,
  APP_SIDEBAR_TOGGLE_ELEMENT_ID,
} from './domSelectors'
import { initShellSidebar } from './shellSidebar'
import { MOBILE_SHELL_MEDIA_QUERY } from './shellBreakpoints'

describe('initShellSidebar', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
    vi.restoreAllMocks()
  })

  it('returns early when #app-shell is missing', () => {
    const getElementById = vi.fn(() => null)
    const addEventListener = vi.fn()
    vi.stubGlobal('document', { getElementById, addEventListener, querySelector: vi.fn() })
    vi.stubGlobal('window', { matchMedia: vi.fn() })

    initShellSidebar()

    expect(getElementById).toHaveBeenCalledWith(APP_SHELL_ELEMENT_ID)
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

    const getElementById = vi.fn((id: string) => {
      if (id === APP_SHELL_ELEMENT_ID) return appShell
      if (id === APP_SIDEBAR_TOGGLE_ELEMENT_ID) return toggle
      if (id === APP_SIDEBAR_BACKDROP_ELEMENT_ID) return backdrop
      return null
    })

    const doc = {
      getElementById,
      addEventListener: vi.fn(),
      querySelector: vi.fn(() => null),
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
})
