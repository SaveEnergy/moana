import { afterEach, describe, expect, it, vi } from 'vitest'

import { SETTINGS_ADD_MEMBER_DIALOG_SELECTOR, SETTINGS_ADD_MEMBER_OPEN_SELECTOR } from './domSelectors'
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'
import {
  initSettingsMemberDialog,
  querySettingsAddMemberDialog,
  querySettingsAddMemberOpenButton,
} from './settingsMemberDialog'

describe('querySettingsAddMemberDialog', () => {
  it('uses SETTINGS_ADD_MEMBER_DIALOG_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelector: (sel: string) => {
        seen = sel
        return null
      },
    } as unknown as ParentNode
    querySettingsAddMemberDialog(root)
    expect(seen).toBe(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR)
  })
})

describe('querySettingsAddMemberOpenButton', () => {
  it('uses SETTINGS_ADD_MEMBER_OPEN_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelector: (sel: string) => {
        seen = sel
        return null
      },
    } as unknown as ParentNode
    querySettingsAddMemberOpenButton(root)
    expect(seen).toBe(SETTINGS_ADD_MEMBER_OPEN_SELECTOR)
  })
})

describe('initSettingsMemberDialog', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('no-ops when the dialog is absent', () => {
    vi.stubGlobal('document', stubDocumentWithoutMainLandmark())
    expect(() => initSettingsMemberDialog()).not.toThrow()
  })

  it('resolves open button from document when dialog has no parent', () => {
    const showModal = vi.fn()
    const dialogAdd = vi.fn()
    const openAddEventListener = vi.fn()
    const openBtn = { addEventListener: openAddEventListener } as unknown as HTMLElement

    const dialog = {
      showModal,
      addEventListener: dialogAdd,
      parentElement: null,
    } as unknown as HTMLDialogElement

    vi.stubGlobal('document', {
      querySelector: (sel: string) => {
        if (sel === SETTINGS_ADD_MEMBER_DIALOG_SELECTOR) return dialog
        if (sel === SETTINGS_ADD_MEMBER_OPEN_SELECTOR) return openBtn
        return null
      },
    })

    initSettingsMemberDialog()

    expect(openAddEventListener).toHaveBeenCalledTimes(1)
  })

  it('wires open button to dialog.showModal and attaches dismiss click handler', () => {
    const showModal = vi.fn()
    const addEventListener = vi.fn()
    const openAddEventListener = vi.fn()
    const openBtn = { addEventListener: openAddEventListener } as unknown as HTMLElement

    const parent = {
      querySelector: (sel: string) => (sel === SETTINGS_ADD_MEMBER_OPEN_SELECTOR ? openBtn : null),
    } as unknown as ParentNode
    const dialog = {
      showModal,
      addEventListener,
      parentElement: parent,
    } as unknown as HTMLDialogElement

    vi.stubGlobal('document', {
      querySelector: (sel: string) =>
        sel === SETTINGS_ADD_MEMBER_DIALOG_SELECTOR ? dialog : null,
    })

    initSettingsMemberDialog()

    expect(openAddEventListener).toHaveBeenCalledWith('click', expect.any(Function))
    const onOpen = openAddEventListener.mock.calls[0][1] as () => void
    onOpen()
    expect(showModal).toHaveBeenCalledTimes(1)

    expect(addEventListener).toHaveBeenCalledWith('click', expect.any(Function))
  })

  it('does not stack listeners when init runs twice', () => {
    const showModal = vi.fn()
    const dialogAdd = vi.fn()
    const openAddEventListener = vi.fn()
    const openBtn = { addEventListener: openAddEventListener } as unknown as HTMLElement

    const parent = {
      querySelector: (sel: string) => (sel === SETTINGS_ADD_MEMBER_OPEN_SELECTOR ? openBtn : null),
    } as unknown as ParentNode
    const dialog = {
      showModal,
      addEventListener: dialogAdd,
      parentElement: parent,
    } as unknown as HTMLDialogElement

    vi.stubGlobal('document', {
      querySelector: (sel: string) =>
        sel === SETTINGS_ADD_MEMBER_DIALOG_SELECTOR ? dialog : null,
    })

    initSettingsMemberDialog()
    initSettingsMemberDialog()

    expect(openAddEventListener).toHaveBeenCalledTimes(1)
    expect(dialogAdd).toHaveBeenCalledTimes(1)
  })

  it('resolves dialog under main.app-main when the landmark exists', () => {
    const showModal = vi.fn()
    const dialogAdd = vi.fn()
    const openAddEventListener = vi.fn()
    const openBtn = { addEventListener: openAddEventListener } as unknown as HTMLElement
    const parent = {
      querySelector: (sel: string) => (sel === SETTINGS_ADD_MEMBER_OPEN_SELECTOR ? openBtn : null),
    } as unknown as ParentNode
    const dialog = {
      showModal,
      addEventListener: dialogAdd,
      parentElement: parent,
    } as unknown as HTMLDialogElement
    const main = {
      querySelector: (sel: string) =>
        sel === SETTINGS_ADD_MEMBER_DIALOG_SELECTOR ? dialog : null,
    } as unknown as ParentNode

    vi.stubGlobal('document', stubDocumentMainLandmark(main))

    initSettingsMemberDialog()

    expect(openAddEventListener).toHaveBeenCalledTimes(1)
    expect(dialogAdd).toHaveBeenCalledWith('click', expect.any(Function))
  })
})
