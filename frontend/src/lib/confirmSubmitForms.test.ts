import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  attachConfirmBeforeSubmit,
  findDataConfirmForms,
  initConfirmSubmitForms,
  readDataConfirmMessage,
} from './confirmSubmitForms'
import { DATA_CONFIRM_ATTRIBUTE, FORM_DATA_CONFIRM_SELECTOR } from './domSelectors'
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'

function elWithAttr(value: string | null): Element {
  return {
    getAttribute: (name: string) => (name === DATA_CONFIRM_ATTRIBUTE ? value : null),
  } as unknown as Element
}

describe('readDataConfirmMessage', () => {
  it('returns null when attribute is missing', () => {
    expect(readDataConfirmMessage(elWithAttr(null))).toBeNull()
  })

  it('returns null when attribute is blank or whitespace-only', () => {
    expect(readDataConfirmMessage(elWithAttr(''))).toBeNull()
    expect(readDataConfirmMessage(elWithAttr('   \t '))).toBeNull()
    expect(readDataConfirmMessage(elWithAttr('\n\r'))).toBeNull()
    expect(readDataConfirmMessage(elWithAttr('\u2028\u2029'))).toBeNull()
  })

  it('returns trimmed text', () => {
    expect(readDataConfirmMessage(elWithAttr('  Remove this?  '))).toBe('Remove this?')
  })

  it('returns null when the message is only Unicode format chars (e.g. ZWSP) that trim() keeps', () => {
    expect(readDataConfirmMessage(elWithAttr('\u200b'))).toBeNull()
    expect(readDataConfirmMessage(elWithAttr('\u200b \u200c\t'))).toBeNull()
    expect(readDataConfirmMessage(elWithAttr('\ufeff'))).toBeNull()
    expect(readDataConfirmMessage(elWithAttr('\u2060'))).toBeNull()
  })

  it('strips Cf characters so confirm text is visible', () => {
    expect(readDataConfirmMessage(elWithAttr('\u200bReally?\u200b'))).toBe('Really?')
    expect(readDataConfirmMessage(elWithAttr('\ufeffProceed?\ufeff'))).toBe('Proceed?')
  })
})

describe('findDataConfirmForms', () => {
  function formWithConfirm(value: string | null): HTMLFormElement {
    return {
      getAttribute: (name: string) => (name === DATA_CONFIRM_ATTRIBUTE ? value : null),
    } as unknown as HTMLFormElement
  }

  it('returns only forms with a non-blank trimmed message, with message attached', () => {
    const a = formWithConfirm('  Really?  ')
    const b = formWithConfirm(' \t ')
    const c = formWithConfirm(null)
    const root = {
      querySelectorAll: (sel: string) => {
        expect(sel).toBe(FORM_DATA_CONFIRM_SELECTOR)
        return [a, b, c]
      },
    } as unknown as ParentNode

    expect(findDataConfirmForms(root)).toEqual([{ form: a, message: 'Really?' }])
  })

  it('returns empty array when selector matches nothing', () => {
    const root = {
      querySelectorAll: () => [],
    } as unknown as ParentNode
    expect(findDataConfirmForms(root)).toEqual([])
  })

  it('preserves querySelectorAll document order (stable vs accidental reordering)', () => {
    const first = formWithConfirm('One')
    const second = formWithConfirm('Two')
    const root = {
      querySelectorAll: (sel: string) => {
        expect(sel).toBe(FORM_DATA_CONFIRM_SELECTOR)
        return [first, second]
      },
    } as unknown as ParentNode
    const found = findDataConfirmForms(root)
    expect(found.map((x) => x.form)).toEqual([first, second])
    expect(found.map((x) => x.message)).toEqual(['One', 'Two'])
  })
})

describe('attachConfirmBeforeSubmit', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('calls preventDefault when confirm returns false', () => {
    vi.stubGlobal('confirm', vi.fn(() => false))
    let onSubmit: ((e: SubmitEvent) => void) | null = null
    const form = {
      addEventListener: vi.fn((ev: string, fn: (e: SubmitEvent) => void) => {
        if (ev === 'submit') {
          onSubmit = fn
        }
      }),
    } as unknown as HTMLFormElement
    attachConfirmBeforeSubmit(form, 'Remove?')
    expect(onSubmit).not.toBeNull()
    const e = { preventDefault: vi.fn() } as unknown as SubmitEvent
    onSubmit!(e)
    expect(e.preventDefault).toHaveBeenCalledTimes(1)
  })

  it('does not prevent default when confirm returns true', () => {
    vi.stubGlobal('confirm', vi.fn(() => true))
    let onSubmit: ((e: SubmitEvent) => void) | null = null
    const form = {
      addEventListener: vi.fn((ev: string, fn: (e: SubmitEvent) => void) => {
        if (ev === 'submit') {
          onSubmit = fn
        }
      }),
    } as unknown as HTMLFormElement
    attachConfirmBeforeSubmit(form, 'Remove?')
    const e = { preventDefault: vi.fn() } as unknown as SubmitEvent
    onSubmit!(e)
    expect(e.preventDefault).not.toHaveBeenCalled()
  })

  it('does not stack submit listeners when attach runs twice on the same form', () => {
    const addEventListener = vi.fn()
    const form = { addEventListener } as unknown as HTMLFormElement
    attachConfirmBeforeSubmit(form, 'Once')
    attachConfirmBeforeSubmit(form, 'Twice')
    expect(addEventListener).toHaveBeenCalledTimes(1)
  })

  it('keeps the first confirm message when attach runs twice', () => {
    const confirmFn = vi.fn(() => true)
    vi.stubGlobal('confirm', confirmFn)
    let onSubmit: ((e: SubmitEvent) => void) | null = null
    const form = {
      addEventListener: vi.fn((ev: string, fn: (e: SubmitEvent) => void) => {
        if (ev === 'submit') {
          onSubmit = fn
        }
      }),
    } as unknown as HTMLFormElement
    attachConfirmBeforeSubmit(form, 'First message')
    attachConfirmBeforeSubmit(form, 'Second message')
    expect(onSubmit).not.toBeNull()
    onSubmit!({ preventDefault: vi.fn() } as unknown as SubmitEvent)
    expect(confirmFn).toHaveBeenCalledTimes(1)
    expect(confirmFn).toHaveBeenCalledWith('First message')
  })
})

describe('initConfirmSubmitForms', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('wires submit only for forms whose data-confirm is non-blank after trim', () => {
    const wired = {
      getAttribute: (name: string) =>
        name === DATA_CONFIRM_ATTRIBUTE ? '  Really delete?  ' : null,
      addEventListener: vi.fn(),
    } as unknown as HTMLFormElement

    const skipped = {
      getAttribute: (name: string) => (name === DATA_CONFIRM_ATTRIBUTE ? ' \n\t ' : null),
      addEventListener: vi.fn(),
    } as unknown as HTMLFormElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelectorAll: ((_sel: string) => [wired, skipped]) as unknown as Document['querySelectorAll'],
      }),
    )

    initConfirmSubmitForms()

    expect(wired.addEventListener).toHaveBeenCalledWith('submit', expect.any(Function))
    expect(skipped.addEventListener).not.toHaveBeenCalled()
  })

  it('no-ops when no forms match the selector', () => {
    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelectorAll: ((_sel: string) => []) as unknown as Document['querySelectorAll'],
      }),
    )
    expect(() => initConfirmSubmitForms()).not.toThrow()
  })

  it('does not wire submit when every matching form has a blank trimmed data-confirm', () => {
    const a = {
      getAttribute: (name: string) => (name === DATA_CONFIRM_ATTRIBUTE ? ' \n\t ' : null),
      addEventListener: vi.fn(),
    } as unknown as HTMLFormElement
    const b = {
      getAttribute: (name: string) => (name === DATA_CONFIRM_ATTRIBUTE ? '' : null),
      addEventListener: vi.fn(),
    } as unknown as HTMLFormElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelectorAll: ((_sel: string) => [a, b]) as unknown as Document['querySelectorAll'],
      }),
    )

    initConfirmSubmitForms()

    expect(a.addEventListener).not.toHaveBeenCalled()
    expect(b.addEventListener).not.toHaveBeenCalled()
  })

  it('does not attach a second submit listener when init runs twice', () => {
    const form = {
      getAttribute: (name: string) =>
        name === DATA_CONFIRM_ATTRIBUTE ? 'Really delete?  ' : null,
      addEventListener: vi.fn(),
    } as unknown as HTMLFormElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelectorAll: ((_sel: string) => [form]) as unknown as Document['querySelectorAll'],
      }),
    )

    initConfirmSubmitForms()
    initConfirmSubmitForms()
    expect(form.addEventListener).toHaveBeenCalledTimes(1)
  })

  it('queries forms under main.app-main when that landmark exists', () => {
    const form = {
      getAttribute: (name: string) =>
        name === DATA_CONFIRM_ATTRIBUTE ? 'Really delete?  ' : null,
      addEventListener: vi.fn(),
    } as unknown as HTMLFormElement

    const main = {
      querySelectorAll: (sel: string) => {
        expect(sel).toBe(FORM_DATA_CONFIRM_SELECTOR)
        return [form]
      },
    } as unknown as ParentNode

    vi.stubGlobal('document', stubDocumentMainLandmark(main))

    initConfirmSubmitForms()

    expect(form.addEventListener).toHaveBeenCalledWith('submit', expect.any(Function))
  })
})
