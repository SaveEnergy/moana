import { CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT } from './categoryColor'
import { stripCfTrimEdges } from './unicodeCf'

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
 * @returns `null` when `id` is missing or whitespace-only after **`stripCfTrimEdges`** (incl. NBSP). **`isCustom`** is true when **`data-custom`** is **`"1"`** after **`stripCfTrimEdges`**. **`name`** is left as in **`dataset`** (spacing preserved).
 */
export function readCategoryEditRowDataset(ds: DOMStringMap): CategoryEditRowData | null {
  const id = stripCfTrimEdges(ds.id ?? '')
  if (!id) {
    return null
  }
  return {
    id,
    name: ds.name ?? '',
    rawColor: stripCfTrimEdges(ds.color ?? ''),
    isCustom: stripCfTrimEdges(ds.custom ?? '') === '1',
    customHex: stripCfTrimEdges(ds.customHex ?? CATEGORY_MODAL_CUSTOM_COLOR_INPUT_DEFAULT),
    iconVal: stripCfTrimEdges(ds.icon ?? ''),
  }
}
