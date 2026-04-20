import { describe, expect, it } from 'vitest'

import { STRING_NEEDS_EDGE_TRIM, trimEdgesIfNeeded } from './trimEdges'

describe('STRING_NEEDS_EDGE_TRIM', () => {
  it('is false for empty and tight strings', () => {
    expect(STRING_NEEDS_EDGE_TRIM.test('')).toBe(false)
    expect(STRING_NEEDS_EDGE_TRIM.test('#close')).toBe(false)
    expect(STRING_NEEDS_EDGE_TRIM.test('2020-06-15T14:30:00.000Z')).toBe(false)
  })

  it('is true when leading or trailing whitespace is present', () => {
    expect(STRING_NEEDS_EDGE_TRIM.test(' a')).toBe(true)
    expect(STRING_NEEDS_EDGE_TRIM.test('a ')).toBe(true)
    expect(STRING_NEEDS_EDGE_TRIM.test('\t#x')).toBe(true)
  })
})

describe('trimEdgesIfNeeded', () => {
  it('returns the same string when edges are already tight', () => {
    const s = 'preset'
    expect(trimEdgesIfNeeded(s)).toBe(s)
  })

  it('trims when the probe matches', () => {
    expect(trimEdgesIfNeeded('  hi  ')).toBe('hi')
  })
})
