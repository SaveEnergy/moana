import { describe, expect, it, vi } from 'vitest'

import { SETTINGS_ADD_MEMBER_OPEN_SELECTOR } from './domSelectors'
import { querySettingsAddMemberInitContext } from './settingsMemberDialogQueries'

describe('querySettingsAddMemberInitContext', () => {
  it('queries SETTINGS_ADD_MEMBER_OPEN_SELECTOR on dialog.parentElement when set', () => {
    const openBtn = { _: 'open' } as unknown as HTMLElement
    const parentQuerySelector = vi.fn((sel: string) =>
      sel === SETTINGS_ADD_MEMBER_OPEN_SELECTOR ? openBtn : null,
    )
    const parent = { querySelector: parentQuerySelector } as unknown as ParentNode
    const dialog = { parentElement: parent } as unknown as HTMLDialogElement
    const contentRoot = {
      querySelector: vi.fn(() => {
        throw new Error('contentRoot should not be queried when dialog.parentElement is set')
      }),
    } as unknown as ParentNode

    const out = querySettingsAddMemberInitContext(contentRoot, dialog)
    expect(out.openBtn).toBe(openBtn)
    expect(parentQuerySelector).toHaveBeenCalledWith(SETTINGS_ADD_MEMBER_OPEN_SELECTOR)
  })

  it('falls back to contentRoot when dialog.parentElement is null', () => {
    const openBtn = { _: 'open' } as unknown as HTMLElement
    const contentRoot = {
      querySelector: vi.fn((sel: string) =>
        sel === SETTINGS_ADD_MEMBER_OPEN_SELECTOR ? openBtn : null,
      ),
    } as unknown as ParentNode
    const dialog = { parentElement: null } as unknown as HTMLDialogElement

    const out = querySettingsAddMemberInitContext(contentRoot, dialog)
    expect(out.openBtn).toBe(openBtn)
    expect(contentRoot.querySelector).toHaveBeenCalledWith(SETTINGS_ADD_MEMBER_OPEN_SELECTOR)
  })
})
