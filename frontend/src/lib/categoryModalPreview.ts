import { resolveCategoryModalPreviewBackground, shouldUpdateCategoryModalPreviewBackground } from './categoryColor'
import { shouldRepaintCategoryModalIconPreview } from './categoryModalIconPreview'
import {
  CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
  CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM,
  CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR,
  CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
  CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS,
  MOANA_ICON_CAT_PREVIEW_CLASS,
  MOANA_ICON_SVG_SELECTOR,
} from './domSelectors'
import { getFormRadioGroupValue } from './radioMap'
import { createRafScheduler, type RafScheduler } from './scheduleAnimationFrame'

export type CategoryModalPreviewSyncOpts = {
  colorRadioTarget?: HTMLInputElement
  iconRadioTarget?: HTMLInputElement
}

export type CategoryModalPreviewController = {
  sync: (opts?: CategoryModalPreviewSyncOpts) => void
  /** Coalesced custom-color `input` → {@link sync}. */
  raf: RafScheduler
  resetPaintState: () => void
}

/**
 * Live color strip + icon well for the category modal: reads radio groups, skips redundant DOM when helpers allow,
 * coalesces custom-color **`input`** via {@link createRafScheduler}.
 */
export function createCategoryModalPreviewController(deps: {
  form: HTMLFormElement
  colorNativeInput: HTMLInputElement | null | undefined
  iconRadioByValue: ReadonlyMap<string, HTMLInputElement>
  preview: HTMLElement
  iconWrap: HTMLElement
}): CategoryModalPreviewController {
  const { form: catForm, colorNativeInput, iconRadioByValue, preview: catPreview, iconWrap: catIconWrap } = deps

  let lastPaintedIconGroupValue: string | undefined
  let lastResolvedPreviewBackground: string | undefined

  function sync(opts?: CategoryModalPreviewSyncOpts) {
    const colorVal = getFormRadioGroupValue(
      catForm,
      CATEGORY_MODAL_COLOR_RADIO_GROUP_NAME,
      opts?.colorRadioTarget,
    )
    const nativeForPreview =
      colorVal === CATEGORY_MODAL_COLOR_RADIO_VALUE_CUSTOM ? colorNativeInput?.value : undefined
    const nextBg = resolveCategoryModalPreviewBackground(colorVal || undefined, nativeForPreview)
    if (shouldUpdateCategoryModalPreviewBackground(lastResolvedPreviewBackground, nextBg)) {
      catPreview.style.background = nextBg
      lastResolvedPreviewBackground = nextBg
    }

    const iconVal = getFormRadioGroupValue(
      catForm,
      CATEGORY_MODAL_ICON_RADIO_GROUP_NAME,
      opts?.iconRadioTarget,
    )
    if (!shouldRepaintCategoryModalIconPreview(lastPaintedIconGroupValue, iconVal)) {
      return
    }
    lastPaintedIconGroupValue = iconVal

    const ir =
      iconRadioByValue.get(iconVal) ??
      catForm.querySelector<HTMLInputElement>(CATEGORY_MODAL_ICON_RADIO_CHECKED_SELECTOR)
    catIconWrap.innerHTML = ''
    const isAutoIcon = !ir?.value
    catIconWrap.classList.toggle(CATEGORY_MODAL_PREVIEW_ICON_AUTO_CLASS, isAutoIcon)
    if (isAutoIcon) {
      catIconWrap.textContent = 'A'
      return
    }
    const label = ir.closest('label')
    const svg = label?.querySelector(MOANA_ICON_SVG_SELECTOR)
    if (svg) {
      const clone = svg.cloneNode(true) as SVGElement
      clone.classList.add(MOANA_ICON_CAT_PREVIEW_CLASS)
      catIconWrap.appendChild(clone)
    }
  }

  const raf = createRafScheduler(sync)

  function resetPaintState() {
    lastPaintedIconGroupValue = undefined
    lastResolvedPreviewBackground = undefined
  }

  return { sync, raf, resetPaintState }
}
