import { describe, expect, it } from 'vitest'

import {
  CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT,
  CATEGORY_MODAL_DEFAULT_PREVIEW_BG,
  resolveCategoryModalPreviewBackground,
  sanitizeCategoryCustomHex,
  shouldUpdateCategoryModalPreviewBackground,
} from './categoryColor'
import { CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM } from './domSelectors'

describe('category modal custom color default', () => {
  it('matches categories.html native color input initial value', () => {
    expect(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT).toBe('#818cf8')
  })
})

describe('shouldUpdateCategoryModalPreviewBackground', () => {
  it('updates when nothing was painted yet', () => {
    expect(shouldUpdateCategoryModalPreviewBackground(undefined, '#ffffff')).toBe(true)
  })

  it('skips when the resolved background is unchanged', () => {
    expect(shouldUpdateCategoryModalPreviewBackground('#abc123', '#abc123')).toBe(false)
    expect(
      shouldUpdateCategoryModalPreviewBackground(
        CATEGORY_MODAL_DEFAULT_PREVIEW_BG,
        CATEGORY_MODAL_DEFAULT_PREVIEW_BG,
      ),
    ).toBe(false)
  })

  it('updates when the resolved background differs', () => {
    expect(shouldUpdateCategoryModalPreviewBackground('#111111', '#222222')).toBe(true)
  })
})

describe('resolveCategoryModalPreviewBackground', () => {
  it('uses auto gradient when no preset is checked', () => {
    expect(resolveCategoryModalPreviewBackground(undefined, undefined)).toBe(
      CATEGORY_MODAL_DEFAULT_PREVIEW_BG,
    )
    expect(resolveCategoryModalPreviewBackground('', undefined)).toBe(CATEGORY_MODAL_DEFAULT_PREVIEW_BG)
  })

  it('uses preset hex when a non-custom radio is selected', () => {
    expect(resolveCategoryModalPreviewBackground('#abc123', undefined)).toBe('#abc123')
  })

  it('ignores native color input when a preset swatch is selected', () => {
    expect(
      resolveCategoryModalPreviewBackground('#abc123', '#ff0000'),
    ).toBe('#abc123')
  })

  it('returns the same native color string when custom is selected and hex has no edge whitespace', () => {
    const native = '#abcdef'
    expect(
      resolveCategoryModalPreviewBackground(CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM, native),
    ).toBe(native)
  })

  it('trims checked radio value for preset and custom branches', () => {
    expect(resolveCategoryModalPreviewBackground('  #abc123  ', undefined)).toBe('#abc123')
    expect(
      resolveCategoryModalPreviewBackground(`  ${CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM}  `, '#abcdef'),
    ).toBe('#abcdef')
  })

  it('uses auto gradient when checked value is whitespace-only', () => {
    expect(resolveCategoryModalPreviewBackground('  \t  ', undefined)).toBe(CATEGORY_MODAL_DEFAULT_PREVIEW_BG)
  })

  it('uses native color when custom radio is selected', () => {
    expect(
      resolveCategoryModalPreviewBackground(
        CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
        '#abcdef',
      ),
    ).toBe('#abcdef')
  })

  it('falls back when custom is selected but native is empty or whitespace', () => {
    expect(
      resolveCategoryModalPreviewBackground(CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM, undefined),
    ).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
    expect(
      resolveCategoryModalPreviewBackground(CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM, '  \t '),
    ).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
  })
})

describe('sanitizeCategoryCustomHex', () => {
  it('returns the same string when input is already #RRGGBB (no trim)', () => {
    const hex = '#aAbBcC'
    expect(sanitizeCategoryCustomHex(hex)).toBe(hex)
  })

  it('accepts 6-digit hex with leading #', () => {
    expect(sanitizeCategoryCustomHex('#aABBcc')).toBe('#aABBcc')
    expect(sanitizeCategoryCustomHex('  #00ff00  ')).toBe('#00ff00')
  })

  it('rejects invalid values and returns fallback', () => {
    expect(sanitizeCategoryCustomHex('', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('   ', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('#fff', '#abc')).toBe('#abc')
    expect(sanitizeCategoryCustomHex('not-a-color')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
  })

  it('rejects wrong-length or non-hex six-digit shapes', () => {
    expect(sanitizeCategoryCustomHex('#12')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
    expect(sanitizeCategoryCustomHex('#12345')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
    expect(sanitizeCategoryCustomHex('#1234567')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
    expect(sanitizeCategoryCustomHex('#12345g')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
    expect(sanitizeCategoryCustomHex('#aabbccdd')).toBe(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT)
  })
})
