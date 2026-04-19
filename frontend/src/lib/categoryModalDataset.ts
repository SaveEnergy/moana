import { CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT } from './categoryColor'

/**
 * Normalized payload from a category row **Edit** button (`categories.html` `data-*`).
 * Pure helper so `initCategoryModal` wiring stays testable without a full DOM harness.
 */
export type CategoryEditRowData = {
  id: string
  name: string
  rawColor: string
  isCustom: boolean
  customHex: string
  iconVal: string
}

/**
 * Parse `dataset` from `.cat-modal-open-edit` (or any element with the same attributes).
 * @returns `null` when `id` is missing or whitespace-only after trim.
 */
export function readCategoryEditRowDataset(ds: DOMStringMap): CategoryEditRowData | null {
  const id = (ds.id ?? '').trim()
  if (!id) {
    return null
  }
  return {
    id,
    name: ds.name ?? '',
    rawColor: (ds.color ?? '').trim(),
    isCustom: ds.custom === '1',
    customHex: (ds.customHex ?? CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT).trim(),
    iconVal: (ds.icon ?? '').trim(),
  }
}
