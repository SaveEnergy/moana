/** Mobile drawer shell — matches CSS (`02-shell.css`) and design.md. */
export const MOBILE_SHELL_MAX_WIDTH_PX = 1023 as const

/** Media query for when the sidebar is off-canvas with backdrop. */
export const MOBILE_SHELL_MEDIA_QUERY = `(max-width: ${MOBILE_SHELL_MAX_WIDTH_PX}px)` as const

/**
 * Subscribe to media query changes with `addEventListener` when it is a function,
 * else `addListener` (legacy WebKit / odd environments). Returns unsubscribe for tests and symmetry.
 */
export function onMediaQueryChange(
  mq: MediaQueryList,
  fn: (this: MediaQueryList, ev: MediaQueryListEvent) => void,
): () => void {
  if (typeof mq.addEventListener === 'function') {
    mq.addEventListener('change', fn)
    return () => {
      mq.removeEventListener('change', fn)
    }
  }
  mq.addListener(fn)
  return () => {
    mq.removeListener(fn)
  }
}
