import { describe, it, expect } from 'vitest'
import { formatLocalTimeLabel } from './localTime'

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
