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

/** Fill all `time.js-local-time[datetime]` elements in `root`. */
export function applyLocalTimeElements(root: ParentNode = document): void {
  const nodes = root.querySelectorAll<HTMLTimeElement>('time.js-local-time[datetime]')
  if (nodes.length === 0) {
    return
  }
  /** Same ISO often repeats (same-minute entries); format once per distinct attribute value. */
  const byIso = new Map<string, string>()
  for (const el of nodes) {
    const iso = el.getAttribute('datetime')
    if (!iso) continue
    let label = byIso.get(iso)
    if (label === undefined) {
      const formatted = formatLocalTimeLabel(iso)
      if (!formatted) continue
      byIso.set(iso, formatted)
      label = formatted
    }
    el.textContent = label
  }
}
