import { afterEach, describe, expect, it, vi } from 'vitest'

import { createCategoryModalPreviewController } from './categoryModalPreview'
import { CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS } from './domSelectors'

function stubFormWithRadios(colorVal: string, iconVal: string): HTMLFormElement {
  const colorEl = { name: 'color', value: colorVal, checked: true } as HTMLInputElement
  const iconEl = { name: 'icon', value: iconVal, checked: true } as HTMLInputElement
  return {
    elements: {
      namedItem: (n: string) => {
        if (n === 'color') return colorEl
        if (n === 'icon') return iconEl
        return null
      },
    },
    querySelector: vi.fn(() => null),
  } as unknown as HTMLFormElement
}

describe('createCategoryModalPreviewController', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('paints auto icon placeholder when checked icon value is empty', () => {
    const form = stubFormWithRadios('red', '')
    const preview = { style: {} as CSSStyleDeclaration } as unknown as HTMLElement
    const toggle = vi.fn()
    const iconWrap = {
      innerHTML: '',
      textContent: '',
      classList: { toggle },
      appendChild: vi.fn(),
    } as unknown as HTMLElement
    const autoIconRadio = { name: 'icon', value: '', checked: true } as HTMLInputElement

    const ctl = createCategoryModalPreviewController({
      form,
      colorNativeInput: null,
      iconRadioByValue: new Map([['', autoIconRadio]]),
      preview,
      iconWrap,
    })
    ctl.sync()

    expect(iconWrap.textContent).toBe('A')
    expect(toggle).toHaveBeenCalledWith(CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS, true)
  })

  it('coalesces rapid raf.schedule to a single animation frame callback', () => {
    const pending: FrameRequestCallback[] = []
    vi.stubGlobal('requestAnimationFrame', (cb: FrameRequestCallback) => {
      pending.push(cb)
      return pending.length
    })
    vi.stubGlobal('cancelAnimationFrame', vi.fn())

    const form = stubFormWithRadios('red', '')
    const preview = { style: {} as CSSStyleDeclaration } as unknown as HTMLElement
    const iconWrap = {
      innerHTML: '',
      textContent: '',
      classList: { toggle: vi.fn() },
      appendChild: vi.fn(),
    } as unknown as HTMLElement
    const autoIconRadio = { name: 'icon', value: '', checked: true } as HTMLInputElement

    const ctl = createCategoryModalPreviewController({
      form,
      colorNativeInput: null,
      iconRadioByValue: new Map([['', autoIconRadio]]),
      preview,
      iconWrap,
    })

    ctl.raf.schedule()
    ctl.raf.schedule()
    expect(pending).toHaveLength(1)
    pending[0](0)
    expect(iconWrap.textContent).toBe('A')
  })
})
