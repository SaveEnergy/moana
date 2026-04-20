/**
 * Coalesce rapid callbacks to at most one `requestAnimationFrame` per frame
 * (e.g. native `<input type="color">` `input` events while dragging).
 *
 * Call {@link RafScheduler.cancelPending} before a synchronous refresh so a stale
 * frame cannot overwrite newer state (e.g. opening create/edit modal).
 */
export type RafScheduler = {
  schedule: () => void
  cancelPending: () => void
}

export function createRafScheduler(run: () => void): RafScheduler {
  let id: number | null = null
  return {
    schedule(): void {
      if (id !== null) {
        return
      }
      const raf = globalThis.requestAnimationFrame
      if (typeof raf !== 'function') {
        run()
        return
      }
      id = raf(() => {
        id = null
        run()
      })
    },
    cancelPending(): void {
      if (id === null) {
        return
      }
      const caf = globalThis.cancelAnimationFrame
      if (typeof caf === 'function') {
        caf(id)
      }
      id = null
    },
  }
}
