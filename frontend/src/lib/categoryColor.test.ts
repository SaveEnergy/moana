import { describe, expect, it } from 'vitest'

import {
  CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT,
  CATEGORY_MODAL_DEFAULT_PREVIEW_BG,
  resolveCategoryModalPreviewBackground,
  sanitizeCategoryCustomHex,
} from './categoryColor'
import { CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM } from './domSelectors'

describe('category modal custom color default', () => {
  it('matches categories.html native color input initial value', () => {
    expect(CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT).toBe('#818cf8')
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
