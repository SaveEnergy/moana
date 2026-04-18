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
  for (const el of nodes) {
    const iso = el.getAttribute('datetime')
    if (!iso) continue
    const label = formatLocalTimeLabel(iso)
    if (label) el.textContent = label
  }
}
