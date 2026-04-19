import { describe, expect, it } from 'vitest'

import {
  CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
  CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
} from './domSelectors'

/**
 * Guardrails: `internal/assets/templates/categories.html` uses `name="color"` / `name="icon"`
 * on preset radios; `categoryModal.ts` and these selectors must stay aligned.
 */
describe('domSelectors category modal radios', () => {
  it('exports HTML group names matching categories.html', () => {
    expect(CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME).toBe('color')
    expect(CATEGORY_MODAL_ICON_RADIO_GROUP_NAME).toBe('icon')
  })

  it('builds group and :checked selectors from group names', () => {
    expect(CATEGORY_MODAL_COLOR_RADIOS_SELECTOR).toBe(
      `input[type="radio"][name="${CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME}"]`,
    )
    expect(CATEGORY_MODAL_ICON_RADIOS_SELECTOR).toBe(
      `input[type="radio"][name="${CATEGORY_MODAL_ICON_RADIO_GROUP_NAME}"]`,
    )
    expect(CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR).toBe(
      `${CATEGORY_MODAL_COLOR_RADIOS_SELECTOR}:checked`,
    )
    expect(CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR).toBe(
      `${CATEGORY_MODAL_ICON_RADIOS_SELECTOR}:checked`,
    )
  })

  it('keeps custom color radio value in sync with template', () => {
    expect(CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM).toBe('custom')
  })
})
