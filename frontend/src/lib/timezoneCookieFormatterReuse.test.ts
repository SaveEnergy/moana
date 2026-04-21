import { afterEach, describe, expect, it, vi } from 'vitest'

/**
 * Isolated file: `vi.resetModules()` would desync static imports in `timezoneCookie.test.ts`.
 */
describe('setBrowserTimezoneCookie Intl reuse', () => {
  afterEach(() => {
    vi.restoreAllMocks()
    vi.unstubAllGlobals()
    vi.resetModules()
  })

  it('constructs at most one Intl.DateTimeFormat across repeated calls', async () => {
    const ctor = vi.spyOn(Intl, 'DateTimeFormat')
    vi.spyOn(Intl.DateTimeFormat.prototype, 'resolvedOptions').mockReturnValue({
      calendar: 'gregory',
      locale: 'en-US',
      numberingSystem: 'latn',
      timeZone: 'Europe/Berlin',
    } as Intl.ResolvedDateTimeFormatOptions)

    let jar = ''
    vi.stubGlobal('document', {
      get cookie() {
        return jar
      },
      set cookie(v: string) {
        jar = v
      },
    })

    const { setBrowserTimezoneCookie } = await import('./timezoneCookie')
    setBrowserTimezoneCookie()
    setBrowserTimezoneCookie()

    expect(ctor.mock.calls.length).toBe(1)
  })
})
