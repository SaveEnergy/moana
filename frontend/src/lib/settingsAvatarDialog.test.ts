import { afterEach, describe, expect, it, vi } from 'vitest'

import { SETTINGS_AVATAR_DIALOG_SELECTOR, SETTINGS_AVATAR_OPEN_SELECTOR } from './domSelectors'
import { stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'
import { initSettingsAvatarDialog } from './settingsAvatarDialog'

describe('initSettingsAvatarDialog', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('no-ops when the dialog is absent', () => {
    vi.stubGlobal('document', stubDocumentWithoutMainLandmark())
    expect(() => initSettingsAvatarDialog()).not.toThrow()
  })

  it('wires open button to showModal and attaches dialog click handler', () => {
    const showModal = vi.fn()
    const dialogAdd = vi.fn()
    const openAddEventListener = vi.fn()
    const openBtn = { addEventListener: openAddEventListener } as unknown as HTMLElement
    const parent = {
      querySelector: (sel: string) => (sel === SETTINGS_AVATAR_OPEN_SELECTOR ? openBtn : null),
    } as unknown as ParentNode
    const dialog = {
      showModal,
      addEventListener: dialogAdd,
      parentElement: parent,
      querySelector: () => null,
    } as unknown as HTMLDialogElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelector: (sel: string) =>
          sel === SETTINGS_AVATAR_DIALOG_SELECTOR ? dialog : null,
      }),
    )

    initSettingsAvatarDialog()

    expect(openAddEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    const onOpen = openAddEventListener.mock.calls[0][1] as () => void
    onOpen()
    expect(showModal).toHaveBeenCalledTimes(1)
    expect(dialogAdd).toHaveBeenCalledWith('click', expect.any(Function))
  })
})
