/**
 * Category modal preview icon well: DOM rebuild only depends on the checked icon radio `value`
 * (`""` = auto placeholder). When that value matches the last painted state, {@link initCategoryModal}
 * can skip `innerHTML` / SVG clone after updating the color strip.
 */
export function shouldRepaintCategoryModalIconPreview(
  lastPaintedIconValue: string | undefined,
  currentIconValue: string,
): boolean {
  return lastPaintedIconValue !== currentIconValue
}
