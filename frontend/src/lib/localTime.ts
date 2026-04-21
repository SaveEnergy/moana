import { queryBootContentAll, resolveContentQueryRoot } from './contentRoot'
import { LOCAL_TIME_ELEMENTS_SELECTOR, TIME_DATETIME_ATTRIBUTE } from './domSelectors'
import { trimEdgesIfNeeded } from './trimEdges'

/** Reused across rows so hydrating many `<time>` nodes does not allocate a formatter per call. */
const localTimeFormatter = new Intl.DateTimeFormat(undefined, {
  hour: 'numeric',
  minute: '2-digit',
})

/**
 * Normalize `datetime` attribute text for parsing / memo keys — **`trim`** only when edges are whitespace.
 * (Server-rendered ISO values are usually already tight; avoids an allocation on hot paths.)
 */
export function normalizeIsoDatetimeAttr(iso: string): string | null {
  const t = trimEdgesIfNeeded(iso)
  return t ? t : null
}

/** Parse + format when **`t`** is already **`normalizeIsoDatetimeAttr`** output (memo cache miss — skips second edge probe). */
function formatLocalTimeLabelNormalized(t: string): string | null {
  const d = new Date(t)
  if (Number.isNaN(d.getTime())) return null
  return localTimeFormatter.format(d)
}

/** Format ISO datetime for inline display (matches previous `toLocaleTimeString` behavior). */
export function formatLocalTimeLabel(iso: string): string | null {
  const t = normalizeIsoDatetimeAttr(iso)
  if (!t) return null
  return formatLocalTimeLabelNormalized(t)
}

/**
 * Memoize `formatLocalTimeLabel` per distinct ISO string (exports for unit tests).
 * Valid outputs live in a `Map`; known-invalid attrs are recorded in a `Set` so repeated bad values do not re-hit `Date` parsing.
 */
export function createLocalTimeLabelMemo(): (iso: string) => string | undefined {
  const ok = new Map<string, string>()
  const bad = new Set<string>()
  return (iso: string) => {
    const key = normalizeIsoDatetimeAttr(iso)
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
    const formatted = formatLocalTimeLabelNormalized(key)
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
 * When `root` is the global `document` (including the default no-arg boot call), uses {@link queryBootContentAll} with **`LOCAL_TIME_ELEMENTS_SELECTOR`**; otherwise {@link resolveContentQueryRoot} then **`querySelectorAll`**.
 * Skips assigning `textContent` when it already equals the computed label (less DOM churn on repeat boot / re-scan).
 */
export function applyLocalTimeElements(root: ParentNode = document): void {
  const nodes =
    typeof document !== 'undefined' && root === document
      ? queryBootContentAll<HTMLTimeElement>(LOCAL_TIME_ELEMENTS_SELECTOR)
      : resolveContentQueryRoot(root).querySelectorAll<HTMLTimeElement>(LOCAL_TIME_ELEMENTS_SELECTOR)
  if (nodes.length === 0) {
    return
  }
  for (let i = 0, n = nodes.length; i < n; i++) {
    const el = nodes[i]
    const iso = el.getAttribute(TIME_DATETIME_ATTRIBUTE)
    if (!iso) {
      continue
    }
    const label = labelForIso(iso)
    if (!label) {
      continue
    }
    if (el.textContent !== label) {
      el.textContent = label
    }
  }
}
