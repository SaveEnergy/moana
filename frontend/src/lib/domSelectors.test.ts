import { describe, expect, it } from 'vitest'

import {
  CATEGORY_MODAL_ELEMENT_ID,
  CATEGORY_MODAL_FORM_ELEMENT_ID,
  CATEGORY_MODAL_FORM_SELECTOR,
  CATEGORY_MODAL_ID_INPUT_ELEMENT_ID,
  CATEGORY_MODAL_ID_INPUT_SELECTOR,
  CATEGORY_MODAL_NAME_ELEMENT_ID,
  CATEGORY_MODAL_NAME_SELECTOR,
  CATEGORY_MODAL_OPEN_CREATE_ELEMENT_ID,
  CATEGORY_MODAL_OPEN_CREATE_SELECTOR,
  CATEGORY_MODAL_PREVIEW_ELEMENT_ID,
  CATEGORY_MODAL_PREVIEW_ICON_ELEMENT_ID,
  CATEGORY_MODAL_PREVIEW_ICON_SELECTOR,
  CATEGORY_MODAL_PREVIEW_SELECTOR,
  CATEGORY_MODAL_SELECTOR,
  CATEGORY_MODAL_SUBMIT_ELEMENT_ID,
  CATEGORY_MODAL_SUBMIT_SELECTOR,
  CATEGORY_MODAL_TITLE_ELEMENT_ID,
  CATEGORY_MODAL_TITLE_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_COLOR_RADIOS_SELECTOR,
  CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
  CATEGORY_PAGE_INTRO_SECTION_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
  CATEGORY_MODAL_ICON_RADIOS_SELECTOR,
  DATA_CONFIRM_ATTRIBUTE,
  FORM_DATA_CONFIRM_SELECTOR,
  APP_SHELL_ELEMENT_ID,
  APP_SIDEBAR_BACKDROP_ELEMENT_ID,
  APP_SIDEBAR_TOGGLE_ELEMENT_ID,
  APP_SHELL_SELECTOR,
  APP_SIDEBAR_BACKDROP_SELECTOR,
  APP_SIDEBAR_TOGGLE_SELECTOR,
  HISTORY_SORT_ELEMENT_ID,
  HISTORY_SORT_SELECTOR,
  SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID,
  SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID,
  SETTINGS_ADD_MEMBER_DIALOG_SELECTOR,
  SETTINGS_ADD_MEMBER_OPEN_SELECTOR,
  LOCAL_TIME_ELEMENTS_SELECTOR,
  TIME_DATETIME_ATTRIBUTE,
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

describe('domSelectors categories page', () => {
  it('keeps intro section class aligned with categories.html', () => {
    expect(CATEGORY_PAGE_INTRO_SECTION_SELECTOR).toBe('.cat-page-intro')
  })
})

describe('domSelectors category modal chrome', () => {
  it('builds #id selectors for dialog, form, and fields', () => {
    expect(CATEGORY_MODAL_SELECTOR).toBe(`#${CATEGORY_MODAL_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_FORM_SELECTOR).toBe(`#${CATEGORY_MODAL_FORM_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_ID_INPUT_SELECTOR).toBe(`#${CATEGORY_MODAL_ID_INPUT_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_TITLE_SELECTOR).toBe(`#${CATEGORY_MODAL_TITLE_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_SUBMIT_SELECTOR).toBe(`#${CATEGORY_MODAL_SUBMIT_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_PREVIEW_SELECTOR).toBe(`#${CATEGORY_MODAL_PREVIEW_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_PREVIEW_ICON_SELECTOR).toBe(`#${CATEGORY_MODAL_PREVIEW_ICON_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_NAME_SELECTOR).toBe(`#${CATEGORY_MODAL_NAME_ELEMENT_ID}`)
    expect(CATEGORY_MODAL_OPEN_CREATE_SELECTOR).toBe(`#${CATEGORY_MODAL_OPEN_CREATE_ELEMENT_ID}`)
  })
})

describe('domSelectors data-confirm forms', () => {
  it('builds the form query from the data-confirm attribute name', () => {
    expect(FORM_DATA_CONFIRM_SELECTOR).toBe(`form[${DATA_CONFIRM_ATTRIBUTE}]`)
  })
})

describe('domSelectors local time', () => {
  it('targets js-local-time elements with the datetime attribute', () => {
    expect(LOCAL_TIME_ELEMENTS_SELECTOR).toBe(`time.js-local-time[${TIME_DATETIME_ATTRIBUTE}]`)
  })
})

describe('domSelectors app shell', () => {
  it('builds #id selectors for shell, toggle, and backdrop', () => {
    expect(APP_SHELL_SELECTOR).toBe(`#${APP_SHELL_ELEMENT_ID}`)
    expect(APP_SIDEBAR_TOGGLE_SELECTOR).toBe(`#${APP_SIDEBAR_TOGGLE_ELEMENT_ID}`)
    expect(APP_SIDEBAR_BACKDROP_SELECTOR).toBe(`#${APP_SIDEBAR_BACKDROP_ELEMENT_ID}`)
  })
})

describe('domSelectors history sort', () => {
  it('builds #id from HISTORY_SORT_ELEMENT_ID', () => {
    expect(HISTORY_SORT_SELECTOR).toBe(`#${HISTORY_SORT_ELEMENT_ID}`)
  })
})

describe('domSelectors settings add-member', () => {
  it('builds #id selectors for dialog and open control', () => {
    expect(SETTINGS_ADD_MEMBER_DIALOG_SELECTOR).toBe(`#${SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID}`)
    expect(SETTINGS_ADD_MEMBER_OPEN_SELECTOR).toBe(`#${SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID}`)
  })
})
