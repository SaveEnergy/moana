import { describe, expect, it } from 'vitest'

import { APP_SIDEBAR_CLOSE_SELECTOR } from './domSelectors'
import { shouldCloseMobileSidebarFromShellClick } from './mobileShellDismiss'

describe('shouldCloseMobileSidebarFromShellClick', () => {
  it('returns true when target is the backdrop element', () => {
    const backdrop = { closest: () => null } as unknown as Element
    const e = { target: backdrop } as unknown as MouseEvent
    expect(shouldCloseMobileSidebarFromShellClick(e, backdrop)).toBe(true)
  })

  it('returns true when closest finds app-sidebar-close', () => {
    const closeBtn = {
      closest: (sel: string) => (sel === APP_SIDEBAR_CLOSE_SELECTOR ? closeBtn : null),
    } as unknown as Element
    const e = { target: closeBtn } as unknown as MouseEvent
    expect(shouldCloseMobileSidebarFromShellClick(e, null)).toBe(true)
  })

  it('returns false for unrelated in-shell clicks (e.g. nav link)', () => {
    const link = {
      closest: () => null,
    } as unknown as Element
    const backdrop = {} as unknown as Element
    const e = { target: link } as unknown as MouseEvent
    expect(shouldCloseMobileSidebarFromShellClick(e, backdrop)).toBe(false)
  })
})
