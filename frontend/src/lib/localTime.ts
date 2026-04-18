/** Format ISO datetime for inline display (matches previous main.ts behavior). */
export function formatLocalTimeLabel(iso: string): string | null {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return null
  return d.toLocaleTimeString(undefined, { hour: 'numeric', minute: '2-digit' })
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
