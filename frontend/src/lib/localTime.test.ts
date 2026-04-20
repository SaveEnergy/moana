import { afterEach, describe, expect, it, vi } from 'vitest'

import { LOCAL_TIME_ELEMENTS_SELECTOR, TIME_DATETIME_ATTRIBUTE } from './domSelectors'
import {
  applyLocalTimeElements,
  createLocalTimeLabelMemo,
  formatLocalTimeLabel,
  normalizeIsoDatetimeAttr,
} from './localTime'
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'

describe('normalizeIsoDatetimeAttr', () => {
  it('returns null for empty or whitespace-only', () => {
    expect(normalizeIsoDatetimeAttr('')).toBeNull()
    expect(normalizeIsoDatetimeAttr('  \t  ')).toBeNull()
  })

  it('returns the same string when no leading or trailing whitespace', () => {
    const iso = '2020-06-15T14:30:00.000Z'
    expect(normalizeIsoDatetimeAttr(iso)).toBe(iso)
  })

  it('trims when edges are whitespace', () => {
    const inner = '2020-06-15T14:30:00.000Z'
    expect(normalizeIsoDatetimeAttr(`  ${inner}  `)).toBe(inner)
  })
})

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

  it('trims surrounding whitespace so Date can parse template padding', () => {
    const inner = '2020-06-15T14:30:00.000Z'
    expect(formatLocalTimeLabel(`  ${inner}  `)).toEqual(formatLocalTimeLabel(inner))
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

  it('does not assign textContent on the second pass when the label is already set', () => {
    const iso = '2020-06-15T14:30:00.000Z'
    const expected = formatLocalTimeLabel(iso)
    expect(expected).not.toBeNull()

    let sets = 0
    let text = ''
    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? iso : null),
      get textContent() {
        return text
      },
      set textContent(v: string) {
        sets++
        text = v
      },
    } as unknown as HTMLTimeElement
    const root = { querySelectorAll: () => [el] } as unknown as ParentNode

    applyLocalTimeElements(root)
    expect(sets).toBe(1)
    applyLocalTimeElements(root)
    expect(sets).toBe(1)
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

  it('does not touch global document when an explicit root is passed (Node / Vitest)', () => {
    const root = { querySelectorAll: () => [] } as unknown as ParentNode
    expect(globalThis.document).toBeUndefined()
    expect(() => applyLocalTimeElements(root)).not.toThrow()
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

  it('skips nodes when datetime attribute is empty string', () => {
    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? '' : null),
      textContent: 'keep',
    } as unknown as HTMLTimeElement
    const root = { querySelectorAll: () => [el] } as unknown as ParentNode
    applyLocalTimeElements(root)
    expect(el.textContent).toBe('keep')
  })

  it('hydrates when datetime has leading and trailing whitespace', () => {
    const inner = '2020-06-15T14:30:00.000Z'
    const padded = `  \t${inner}  `
    const expected = formatLocalTimeLabel(inner)
    expect(expected).not.toBeNull()

    const el = {
      getAttribute: (name: string) => (name === TIME_DATETIME_ATTRIBUTE ? padded : null),
      textContent: '',
    } as unknown as HTMLTimeElement
    const root = { querySelectorAll: () => [el] } as unknown as ParentNode

    applyLocalTimeElements(root)

    expect(el.textContent).toBe(expected)
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
    vi.stubGlobal('document', stubDocumentMainLandmark(main) as unknown as Document)

    applyLocalTimeElements()

    expect(qsArg).toBe(LOCAL_TIME_ELEMENTS_SELECTOR)
    expect(el.textContent).toBe(expected)
  })

  it('falls back to document when main.app-main is absent', () => {
    let qsArg = ''
    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelectorAll: ((sel: string) => {
          qsArg = sel
          return []
        }) as unknown as Document['querySelectorAll'],
      }) as unknown as Document,
    )

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

  it('uses the same cache entry for the same instant with surrounding whitespace', () => {
    const labelFor = createLocalTimeLabelMemo()
    const iso = '2020-06-15T14:30:00.000Z'
    expect(labelFor(`  ${iso}`)).toBe(labelFor(`${iso}  `))
  })
})
