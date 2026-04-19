import { afterEach, describe, expect, it, vi } from 'vitest'

import { initSettingsMemberDialog } from './settingsMemberDialog'
import {
  SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID,
  SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID,
} from './domSelectors'

describe('initSettingsMemberDialog', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('no-ops when the dialog is absent', () => {
    vi.stubGlobal('document', {
      getElementById: () => null,
    })
    expect(() => initSettingsMemberDialog()).not.toThrow()
  })

  it('wires open button to dialog.showModal and attaches dismiss click handler', () => {
    const showModal = vi.fn()
    const addEventListener = vi.fn()
    const dialog = { showModal, addEventListener } as unknown as HTMLDialogElement

    const openAddEventListener = vi.fn()
    const openBtn = { addEventListener: openAddEventListener } as unknown as HTMLElement

    vi.stubGlobal('document', {
      getElementById: (id: string) => {
        if (id === SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID) {
          return dialog
        }
        if (id === SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID) {
          return openBtn
        }
        return null
      },
    })

    initSettingsMemberDialog()

    expect(openAddEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    const onOpen = openAddEventListener.mock.calls[0][1] as () => void
    onOpen()
    expect(showModal).toHaveBeenCalledTimes(1)

    expect(addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
  })
})
