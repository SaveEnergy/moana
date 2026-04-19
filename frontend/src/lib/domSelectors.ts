/**
 * Shared selector strings for DOM that is also authored in HTML templates (`internal/assets/templates/`).
 * Centralizes hooks so refactors don’t miss a string copy across keyboard, shell, and dismiss logic.
 */

/** Topbar account menu when expanded (`layout.html` `details.app-user-menu[open]`). */
export const APP_USER_MENU_OPEN_SELECTOR = 'details.app-user-menu[open]' as const

/** Mobile shell root (`layout.html` `id="app-shell"`). */
export const APP_SHELL_ELEMENT_ID = 'app-shell' as const

/** `aria-controls` target; mobile menu affordance (`layout.html`). */
export const APP_SIDEBAR_TOGGLE_ELEMENT_ID = 'app-sidebar-toggle' as const

/** Dimmed overlay behind the off-canvas drawer (`layout.html`). */
export const APP_SIDEBAR_BACKDROP_ELEMENT_ID = 'app-sidebar-backdrop' as const

/** In-drawer close control (`layout.html`); use with `Element.closest()` / `#…` form below. */
export const APP_SIDEBAR_CLOSE_ELEMENT_ID = 'app-sidebar-close' as const

/** `#` + {@link APP_SIDEBAR_CLOSE_ELEMENT_ID} — `mobileShellDismiss` uses `closest()`, not `getElementById`. */
export const APP_SIDEBAR_CLOSE_SELECTOR = `#${APP_SIDEBAR_CLOSE_ELEMENT_ID}` as const

/** History page sort `<select>` (`history.html`). */
export const HISTORY_SORT_ELEMENT_ID = 'history-sort' as const

/** `#` + {@link HISTORY_SORT_ELEMENT_ID} — `queryHistorySortSelect` / `querySelector`. */
export const HISTORY_SORT_SELECTOR = `#${HISTORY_SORT_ELEMENT_ID}` as const

/** Settings → add household member dialog (`settings.html`). */
export const SETTINGS_ADD_MEMBER_DIALOG_ELEMENT_ID = 'settings-add-member-dialog' as const
export const SETTINGS_ADD_MEMBER_OPEN_ELEMENT_ID = 'settings-add-member-open' as const
export const SETTINGS_ADD_MEMBER_CLOSE_ELEMENT_ID = 'settings-add-member-close' as const
export const SETTINGS_ADD_MEMBER_CANCEL_ELEMENT_ID = 'settings-add-member-cancel' as const

export const SETTINGS_ADD_MEMBER_CLOSE_SELECTOR = `#${SETTINGS_ADD_MEMBER_CLOSE_ELEMENT_ID}` as const
export const SETTINGS_ADD_MEMBER_CANCEL_SELECTOR = `#${SETTINGS_ADD_MEMBER_CANCEL_ELEMENT_ID}` as const

/** Arguments for `attachNativeDialogDismiss` on the add-member `dialog`. */
export const SETTINGS_ADD_MEMBER_DISMISS_SELECTORS: readonly string[] = [
  SETTINGS_ADD_MEMBER_CLOSE_SELECTOR,
  SETTINGS_ADD_MEMBER_CANCEL_SELECTOR,
]

/** Categories list + modal (`categories.html`). */
export const CATEGORY_LIST_SECTION_SELECTOR = '.cat-list-section' as const
export const CATEGORY_MODAL_OPEN_EDIT_SELECTOR = '.cat-modal-open-edit' as const

/** Preview strip when icon = Auto (`categoryModal.ts` toggles on `#cat-modal-preview-icon`). */
export const CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS = 'cat-modal-preview-icon--auto' as const

/** Cloned Lucide glyph in the modal preview well. */
export const MOANA_ICON_CAT_PREVIEW_CLASS = 'moana-icon--cat-preview' as const

/** Lucide-derived SVG in templates (`design.md` Icons); category modal clones for preview. */
export const MOANA_ICON_SVG_SELECTOR = 'svg.moana-icon' as const

/** Native `<input type="color">` in the category form (`input` delegation). */
export const CATEGORY_COLOR_NATIVE_CLASS = 'cat-color-native' as const

/** Custom swatch row wrapping the native color input. */
export const CATEGORY_COLOR_SWATCH_CUSTOM_SELECTOR = '.cat-color-swatch--custom' as const

/** Color radio `value` that shows the native color picker (`categories.html`). */
export const CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM = 'custom' as const

/** Selects that radio when the native color input fires `input`. */
export const CATEGORY_MODAL_COLOR_RADIO_CUSTOM_SELECTOR =
  `input[type="radio"][value="${CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM}"]` as const

/** Radio `name` on color presets + custom row (`categories.html`). */
export const CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME = 'color' as const
/** Radio `name` on icon grid (`categories.html`). */
export const CATEGORY_MODAL_ICON_RADIO_GROUP_NAME = 'icon' as const

/** All color radios in the modal form — `radioMap` index by preset (`categories.html`). */
export const CATEGORY_MODAL_COLOR_RADIOS_SELECTOR =
  `input[type="radio"][name="${CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME}"]` as const
/** All icon radios in the modal form (`categories.html`). */
export const CATEGORY_MODAL_ICON_RADIOS_SELECTOR =
  `input[type="radio"][name="${CATEGORY_MODAL_ICON_RADIO_GROUP_NAME}"]` as const

/** Currently selected color preset (`categoryModal.ts` preview). */
export const CATEGORY_MODAL_COLOR_RADIO_CHECKED_SELECTOR =
  `input[type="radio"][name="${CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME}"]:checked` as const
/** Currently selected icon (`categoryModal.ts` preview). */
export const CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR =
  `input[type="radio"][name="${CATEGORY_MODAL_ICON_RADIO_GROUP_NAME}"]:checked` as const

export const CATEGORY_MODAL_ELEMENT_ID = 'cat-modal' as const
export const CATEGORY_MODAL_FORM_ELEMENT_ID = 'cat-modal-form' as const
export const CATEGORY_MODAL_ID_INPUT_ELEMENT_ID = 'cat-modal-id' as const
export const CATEGORY_MODAL_TITLE_ELEMENT_ID = 'cat-modal-title' as const
export const CATEGORY_MODAL_SUBMIT_ELEMENT_ID = 'cat-modal-submit' as const
export const CATEGORY_MODAL_PREVIEW_ELEMENT_ID = 'cat-modal-preview' as const
export const CATEGORY_MODAL_PREVIEW_ICON_ELEMENT_ID = 'cat-modal-preview-icon' as const
export const CATEGORY_MODAL_NAME_ELEMENT_ID = 'cat-modal-name' as const
export const CATEGORY_MODAL_OPEN_CREATE_ELEMENT_ID = 'cat-modal-open-create' as const
export const CATEGORY_MODAL_COLOR_NATIVE_ELEMENT_ID = 'cat-modal-color-native' as const
export const CATEGORY_MODAL_CLOSE_ELEMENT_ID = 'cat-modal-close' as const

export const CATEGORY_MODAL_COLOR_NATIVE_SELECTOR = `#${CATEGORY_MODAL_COLOR_NATIVE_ELEMENT_ID}` as const
export const CATEGORY_MODAL_CLOSE_SELECTOR = `#${CATEGORY_MODAL_CLOSE_ELEMENT_ID}` as const

export const CATEGORY_MODAL_DISMISS_SELECTORS: readonly string[] = [CATEGORY_MODAL_CLOSE_SELECTOR]

/** Destructive / irreversible POST forms opt in with this attribute (`confirmSubmitForms.ts`). */
export const DATA_CONFIRM_ATTRIBUTE = 'data-confirm' as const

/** Forms that wire `confirm()` before submit; matches `[attribute]` across templates. */
export const FORM_DATA_CONFIRM_SELECTOR = `form[${DATA_CONFIRM_ATTRIBUTE}]` as const

/** `<time class="js-local-time">` with an instant (`history.html` and related rows; `localTime.ts`). */
export const LOCAL_TIME_ELEMENTS_SELECTOR = 'time.js-local-time[datetime]' as const

/** ISO timestamp on `<time>` — read by `applyLocalTimeElements`. */
export const TIME_DATETIME_ATTRIBUTE = 'datetime' as const
