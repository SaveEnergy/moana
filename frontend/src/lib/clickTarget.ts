/** Minimal shape for delegated `click` handlers and Vitest stubs (see `dialogDismiss`, `mobileShellDismiss`). */
export type ClickTargetEvent = Pick<Event, 'target'>

/** Vitest: arbitrary `target` as {@link ClickTargetEvent} (not a live DOM event). */
export function stubClickTargetEvent(target: unknown): ClickTargetEvent {
  return { target } as unknown as ClickTargetEvent
}

/**
 * Resolve a `click` event target to an `Element` for delegated handlers.
 * When the hit target is a `Text` node (e.g. label or × character), `event.target`
 * is not an `Element` and has no `closest()` — use `parentElement` instead.
 *
 * Accepts any event with `target` so `click` listeners typed as `Event` need no cast.
 */
export function clickEventTargetElement(e: ClickTargetEvent): Element | null {
  const t = e.target
  if (t == null || typeof t !== 'object') {
    return null
  }
  if (typeof Element !== 'undefined' && t instanceof Element) {
    return t
  }
  const maybeEl = t as { closest?: (selector: string) => Element | null }
  if (typeof maybeEl.closest === 'function') {
    return t as Element
  }
  const parent = (t as { parentElement?: Element | null }).parentElement
  return parent ?? null
}
