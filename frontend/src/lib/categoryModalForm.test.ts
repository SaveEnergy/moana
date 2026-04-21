import { afterEach, describe, expect, it, vi } from 'vitest'

import { attachCategoryModalFormPreviewListeners } from './categoryModalForm'
import type { CategoryModalPreviewController } from './categoryModalPreview'
import {
  CATEGORY_COLOR_NATIVE_CLASS,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
} from './domSelectors'

describe('attachCategoryModalFormPreviewListeners', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  function stubBrowserElementTypes() {
    class ElementStub {}
    class HTMLInputElementStub extends ElementStub {}
    vi.stubGlobal('Element', ElementStub)
    vi.stubGlobal('HTMLInputElement', HTMLInputElementStub)
    return { ElementStub, HTMLInputElementStub }
  }

  it('native color input checks custom radio and schedules raf', () => {
    const { ElementStub } = stubBrowserElementTypes()
    const schedule = vi.fn()
    const previewCtl = {
      raf: { schedule },
      sync: vi.fn(),
      resetPaintState: vi.fn(),
    } as unknown as CategoryModalPreviewController

    const inputFns: EventListener[] = []
    const changeFns: EventListener[] = []
    const form = {
      addEventListener(type: string, fn: EventListener) {
        if (type === 'input') inputFns.push(fn)
        if (type === 'change') changeFns.push(fn)
      },
    } as unknown as HTMLFormElement

    attachCategoryModalFormPreviewListeners(form, previewCtl)

    const customRadio = { checked: false }
    const wrap = { querySelector: () => customRadio }
    const native = new ElementStub() as InstanceType<typeof ElementStub> & {
      classList: { contains: (c: string) => boolean }
      closest: () => typeof wrap
    }
    native.classList = { contains: (c: string) => c === CATEGORY_COLOR_NATIVE_CLASS }
    native.closest = () => wrap

    expect(inputFns).toHaveLength(1)
    inputFns[0]({ target: native } as unknown as Event)

    expect(customRadio.checked).toBe(true)
    expect(schedule).toHaveBeenCalledTimes(1)
  })

  it('change on color group syncs with colorRadioTarget', () => {
    const { HTMLInputElementStub } = stubBrowserElementTypes()
    const sync = vi.fn()
    const previewCtl = {
      raf: { schedule: vi.fn() },
      sync,
      resetPaintState: vi.fn(),
    } as unknown as CategoryModalPreviewController

    const changeFns: EventListener[] = []
    const form = {
      addEventListener(type: string, fn: EventListener) {
        if (type === 'change') changeFns.push(fn)
      },
    } as unknown as HTMLFormElement

    attachCategoryModalFormPreviewListeners(form, previewCtl)

    const colorRadio = new HTMLInputElementStub() as InstanceType<typeof HTMLInputElementStub> & {
      name: string
      value: string
    }
    colorRadio.name = CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME
    colorRadio.value = '#abc'

    expect(changeFns).toHaveLength(1)
    changeFns[0]({ target: colorRadio } as unknown as Event)

    expect(sync).toHaveBeenCalledWith({ colorRadioTarget: colorRadio })
  })

  it('change on icon group syncs with iconRadioTarget', () => {
    const { HTMLInputElementStub } = stubBrowserElementTypes()
    const sync = vi.fn()
    const previewCtl = {
      raf: { schedule: vi.fn() },
      sync,
      resetPaintState: vi.fn(),
    } as unknown as CategoryModalPreviewController

    const changeFns: EventListener[] = []
    const form = {
      addEventListener(type: string, fn: EventListener) {
        if (type === 'change') changeFns.push(fn)
      },
    } as unknown as HTMLFormElement

    attachCategoryModalFormPreviewListeners(form, previewCtl)

    const iconRadio = new HTMLInputElementStub() as InstanceType<typeof HTMLInputElementStub> & {
      name: string
      value: string
    }
    iconRadio.name = CATEGORY_MODAL_ICON_RADIO_GROUP_NAME
    iconRadio.value = 'wallet'

    changeFns[0]({ target: iconRadio } as unknown as Event)

    expect(sync).toHaveBeenCalledWith({ iconRadioTarget: iconRadio })
  })

  it('ignores change targets that are not HTMLInputElement', () => {
    const { ElementStub, HTMLInputElementStub } = stubBrowserElementTypes()
    const sync = vi.fn()
    const previewCtl = {
      raf: { schedule: vi.fn() },
      sync,
      resetPaintState: vi.fn(),
    } as unknown as CategoryModalPreviewController

    const changeFns: EventListener[] = []
    const form = {
      addEventListener(type: string, fn: EventListener) {
        if (type === 'change') changeFns.push(fn)
      },
    } as unknown as HTMLFormElement

    attachCategoryModalFormPreviewListeners(form, previewCtl)

    const notAnInput = new ElementStub()
    expect(notAnInput).not.toBeInstanceOf(HTMLInputElementStub)

    changeFns[0]({ target: notAnInput } as unknown as Event)

    expect(sync).not.toHaveBeenCalled()
  })

  it('does not stack listeners when attachCategoryModalFormPreviewListeners runs twice on the same form', () => {
    stubBrowserElementTypes()
    const previewCtl = {
      raf: { schedule: vi.fn() },
      sync: vi.fn(),
      resetPaintState: vi.fn(),
    } as unknown as CategoryModalPreviewController

    let inputCount = 0
    let changeCount = 0
    const form = {
      addEventListener(type: string, _fn: EventListener) {
        if (type === 'input') inputCount += 1
        if (type === 'change') changeCount += 1
      },
    } as unknown as HTMLFormElement

    attachCategoryModalFormPreviewListeners(form, previewCtl)
    attachCategoryModalFormPreviewListeners(form, previewCtl)

    expect(inputCount).toBe(1)
    expect(changeCount).toBe(1)
  })
})
