/** Reused across rows so hydrating many `<time>` nodes does not allocate a formatter per call. */
const localTimeFormatter = new Intl.DateTimeFormat(undefined, {
  hour: 'numeric',
  minute: '2-digit',
})

/** Format ISO datetime for inline display (matches previous `toLocaleTimeString` behavior). */
export function formatLocalTimeLabel(iso: string): string | null {
  const d = new Date(iso)
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
    if (!iso) {
      return undefined
    }
    if (ok.has(iso)) {
      return ok.get(iso)
    }
    if (bad.has(iso)) {
      return undefined
    }
    const formatted = formatLocalTimeLabel(iso)
    if (!formatted) {
      bad.add(iso)
      return undefined
    }
    ok.set(iso, formatted)
    return formatted
  }
}

/** Fill all `time.js-local-time[datetime]` elements in `root`. */
export function applyLocalTimeElements(root: ParentNode = document): void {
  const nodes = root.querySelectorAll<HTMLTimeElement>('time.js-local-time[datetime]')
  if (nodes.length === 0) {
    return
  }
  const labelFor = createLocalTimeLabelMemo()
  for (const el of nodes) {
    const iso = el.getAttribute('datetime')
    if (!iso) {
      continue
    }
    const label = labelFor(iso)
    if (!label) {
      continue
    }
    el.textContent = label
  }
}
