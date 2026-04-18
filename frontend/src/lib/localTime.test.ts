import { describe, it, expect } from 'vitest'
import { createLocalTimeLabelMemo, formatLocalTimeLabel } from './localTime'

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
