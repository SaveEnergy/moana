import { resolveBootContentQueryRoot, resolveContentQueryRoot } from './contentRoot'
import { LOCAL_TIME_ELEMENTS_SELECTOR, TIME_DATETIME_ATTRIBUTE } from './domSelectors'

/** Reused across rows so hydrating many `<time>` nodes does not allocate a formatter per call. */
const localTimeFormatter = new Intl.DateTimeFormat(undefined, {
  hour: 'numeric',
  minute: '2-digit',
})

/** Format ISO datetime for inline display (matches previous `toLocaleTimeString` behavior); trims first. */
export function formatLocalTimeLabel(iso: string): string | null {
  const t = iso.trim()
  if (!t) return null
  const d = new Date(t)
  if (Number.isNaN(d.getTime())) return null
  return localTimeFormatter.format(d)
}

/**
 * Memoize `formatLocalTimeLabel` per distinct ISO string (exports for unit tests).
 * Valid outputs live in a `Map`; known-invalid attrs are recorded in a `Set` so repeated bad values do not re-hit `Date` parsing.
 */
export function createLocalTimeLabelMemo(): (iso: string) => string | undefined {
  const ok = new Map<string, string>()
  const bad = new Set<string>()
  return (iso: string) => {
    const key = iso.trim()
    if (!key) {
      return undefined
    }
    const cached = ok.get(key)
    if (cached !== undefined) {
      return cached
    }
    if (bad.has(key)) {
      return undefined
    }
    const formatted = formatLocalTimeLabel(key)
    if (!formatted) {
      bad.add(key)
      return undefined
    }
    ok.set(key, formatted)
    return formatted
  }
}

/** One memo for the page so repeat `applyLocalTimeElements` calls reuse parsed ISO labels. */
const labelForIso = createLocalTimeLabelMemo()

/**
 * Fill local clock labels for every `<time>` matching `LOCAL_TIME_ELEMENTS_SELECTOR` in `root`.
 * When `root` is the global `document` (including the default no-arg boot call), uses {@link resolveBootContentQueryRoot}; otherwise {@link resolveContentQueryRoot}.
 */
export function applyLocalTimeElements(root: ParentNode = document): void {
  const scope =
    typeof document !== 'undefined' && root === document
      ? resolveBootContentQueryRoot()
      : resolveContentQueryRoot(root)
  const nodes = scope.querySelectorAll<HTMLTimeElement>(LOCAL_TIME_ELEMENTS_SELECTOR)
  if (nodes.length === 0) {
    return
  }
  for (const el of nodes) {
    const iso = el.getAttribute(TIME_DATETIME_ATTRIBUTE)
    if (!iso) {
      continue
    }
    const label = labelForIso(iso)
    if (!label) {
      continue
    }
    el.textContent = label
  }
}
