import { afterEach, describe, expect, it, vi } from 'vitest'

import { APP_MAIN_SELECTOR, LOCAL_TIME_ELEMENTS_SELECTOR, TIME_DATETIME_ATTRIBUTE } from './domSelectors'
import { applyLocalTimeElements, createLocalTimeLabelMemo, formatLocalTimeLabel } from './localTime'

describe('formatLocalTimeLabel', () => {
  it('returns null for invalid ISO strings', () => {
    expect(formatLocalTimeLabel('')).toBeNull()
    expect(formatLocalTimeLabel('bogus')).toBeNull()
  })

  it('returns a non-empty label for valid ISO datetimes', () => {
    const s = formatLocalTimeLabel('2020-06-15T14:30:00.000Z')
    expect(s).not.toBeNull()
    expect(s!.length).toBeGreaterThan(0)
  })
})

describe('applyLocalTimeElements', () => {
  it('sets textContent from datetime on matching time nodes', () => {
    const iso = '2020-06-15T14:30:00.000Z'
    const expected = formatLocalTimeLabel(iso)
    expect(expected).not.toBeNull()

    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? iso : null),
      textContent: '',
    } as unknown as HTMLTimeElement
    const root = {
      querySelectorAll: () => [el],
    } as unknown as ParentNode

    applyLocalTimeElements(root)

    expect(el.textContent).toBe(expected)
  })

  it('is stable when apply runs twice on the same nodes (double bootApp)', () => {
    const iso = '2020-06-15T14:30:00.000Z'
    const expected = formatLocalTimeLabel(iso)
    expect(expected).not.toBeNull()

    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? iso : null),
      textContent: '',
    } as unknown as HTMLTimeElement
    const root = {
      querySelectorAll: () => [el],
    } as unknown as ParentNode

    applyLocalTimeElements(root)
    applyLocalTimeElements(root)

    expect(el.textContent).toBe(expected)
  })

  it('no-ops when the selector matches nothing', () => {
    const root = { querySelectorAll: () => [] } as unknown as ParentNode
    expect(() => applyLocalTimeElements(root)).not.toThrow()
  })

  it('queries LOCAL_TIME_ELEMENTS_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelectorAll: (sel: string) => {
        seen = sel
        return []
      },
    } as unknown as ParentNode
    applyLocalTimeElements(root)
    expect(seen).toBe(LOCAL_TIME_ELEMENTS_SELECTOR)
  })

  it('does not overwrite text when datetime is invalid', () => {
    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? 'not-a-date' : null),
      textContent: 'keep',
    } as unknown as HTMLTimeElement
    const root = { querySelectorAll: () => [el] } as unknown as ParentNode
    applyLocalTimeElements(root)
    expect(el.textContent).toBe('keep')
  })

  it('skips nodes with a missing datetime attribute', () => {
    const el = {
      getAttribute: () => null,
      textContent: 'keep',
    } as unknown as HTMLTimeElement
    const root = { querySelectorAll: () => [el] } as unknown as ParentNode
    applyLocalTimeElements(root)
    expect(el.textContent).toBe('keep')
  })
})

describe('applyLocalTimeElements default document root', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('scopes to main.app-main when present', () => {
    const iso = '2020-06-15T14:30:00.000Z'
    const expected = formatLocalTimeLabel(iso)
    expect(expected).not.toBeNull()

    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? iso : null),
      textContent: '',
    } as unknown as HTMLTimeElement
    let qsArg = ''
    const main = {
      querySelectorAll: (sel: string) => {
        qsArg = sel
        return sel === LOCAL_TIME_ELEMENTS_SELECTOR ? [el] : []
      },
    } as unknown as ParentNode
    const doc = {
      querySelector: (sel: string) => (sel === APP_MAIN_SELECTOR ? main : null),
    } as unknown as Document
    vi.stubGlobal('document', doc)

    applyLocalTimeElements()

    expect(qsArg).toBe(LOCAL_TIME_ELEMENTS_SELECTOR)
    expect(el.textContent).toBe(expected)
  })

  it('falls back to document when main.app-main is absent', () => {
    let qsArg = ''
    const doc = {
      querySelector: () => null,
      querySelectorAll: (sel: string) => {
        qsArg = sel
        return []
      },
    } as unknown as Document
    vi.stubGlobal('document', doc)

    applyLocalTimeElements()

    expect(qsArg).toBe(LOCAL_TIME_ELEMENTS_SELECTOR)
  })
})

describe('createLocalTimeLabelMemo', () => {
  it('returns the same label for repeated identical ISO strings', () => {
    const labelFor = createLocalTimeLabelMemo()
    const iso = '2020-06-15T14:30:00.000Z'
    expect(labelFor(iso)).toBe(labelFor(iso))
  })

  it('returns undefined for empty input without caching as valid', () => {
    const labelFor = createLocalTimeLabelMemo()
    expect(labelFor('')).toBeUndefined()
  })

  it('returns undefined for invalid ISO and stays undefined on repeat', () => {
    const labelFor = createLocalTimeLabelMemo()
    expect(labelFor('not-a-date')).toBeUndefined()
    expect(labelFor('not-a-date')).toBeUndefined()
  })
})
