import { describe, expect, it } from 'vitest'

import {
  CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT,
  sanitizeCategoryCustomHex,
} from './categoryColor'

describe('category modal custom color default', () => {
  it('matches categories.html native color input initial value', () => {
    expect(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT).toBe('#818cf8')
  })
})

describe('sanitizeCategoryCustomHex', () => {
  it('accepts 6-digit hex with leading #', () => {
    expect(sanitizeCategoryCustomHex('#aABBcc')).toBe('#aABBcc')
    expect(sanitizeCategoryCustomHex('  #00ff00  ')).toBe('#00ff00')
  })

  it('rejects invalid values and returns fallback', () => {
    expect(sanitizeCategoryCustomHex('', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('#fff', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('not-a-color')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
  })
})
