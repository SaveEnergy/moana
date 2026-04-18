import { describe, expect, it } from 'vitest'

import { sanitizeCategoryCustomHex } from './categoryColor'

describe('sanitizeCategoryCustomHex', () => {
  it('accepts 6-digit hex with leading #', () => {
    expect(sanitizeCategoryCustomHex('#aABBcc')).toBe('#aABBcc')
    expect(sanitizeCategoryCustomHex('  #00ff00  ')).toBe('#00ff00')
  })

  it('rejects invalid values and returns fallback', () => {
    expect(sanitizeCategoryCustomHex('', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('#fff', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('not-a-color')).toBe('#818cf8')
  })
})
