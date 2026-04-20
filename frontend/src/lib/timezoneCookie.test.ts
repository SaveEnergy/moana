import { afterEach, describe, expect, it, vi } from 'vitest'

import {
  TIMEZONE_COOKIE_NAME,
  parseMoanaTimezoneCookie,
  setBrowserTimezoneCookie,
  timezoneCookieSegment,
} from './timezoneCookie'

describe('timezoneCookieSegment', () => {
  it('includes encoded tz and expected attrs', () => {
    const s = timezoneCookieSegment('Europe/Berlin')
    expect(s.startsWith(`${TIMEZONE_COOKIE_NAME}=`)).toBe(true)
    expect(s).toContain(encodeURIComponent('Europe/Berlin'))
    expect(s).toContain('Path=/')
    expect(s).toContain('SameSite=Lax')
  })
})

describe('parseMoanaTimezoneCookie', () => {
  it('returns null for empty or unrelated cookies', () => {
    expect(parseMoanaTimezoneCookie('')).toBeNull()
    expect(parseMoanaTimezoneCookie('other=1')).toBeNull()
  })

  it('decodes moana_tz when present', () => {
    expect(parseMoanaTimezoneCookie(`${TIMEZONE_COOKIE_NAME}=Europe%2FBerlin`)).toBe('Europe/Berlin')
    expect(
      parseMoanaTimezoneCookie(`a=1; ${TIMEZONE_COOKIE_NAME}=America%2FNew_York; b=2`),
    ).toBe('America/New_York')
  })

  it('handles a leading semicolon segment (empty first chunk)', () => {
    expect(parseMoanaTimezoneCookie(`; ${TIMEZONE_COOKIE_NAME}=Europe%2FBerlin`)).toBe('Europe/Berlin')
  })

  it('returns null on invalid percent encoding', () => {
    expect(parseMoanaTimezoneCookie(`${TIMEZONE_COOKIE_NAME}=%`)).toBeNull()
  })

  it('returns the first decodable moana_tz when the name is repeated', () => {
    expect(
      parseMoanaTimezoneCookie(
        `${TIMEZONE_COOKIE_NAME}=Europe%2FBerlin; ${TIMEZONE_COOKIE_NAME}=America%2FNew_York`,
      ),
    ).toBe('Europe/Berlin')
  })

  it('skips an undecodable segment and uses a later valid moana_tz', () => {
    expect(
      parseMoanaTimezoneCookie(`${TIMEZONE_COOKIE_NAME}=%; ${TIMEZONE_COOKIE_NAME}=Europe%2FBerlin`),
    ).toBe('Europe/Berlin')
  })

  it('finds moana_tz after many unrelated segments (no split allocation regression)', () => {
    const junk = Array.from({ length: 40 }, (_, i) => `k${i}=v`).join('; ')
    const tail = `${TIMEZONE_COOKIE_NAME}=Europe%2FBerlin`
    expect(parseMoanaTimezoneCookie(`${junk}; ${tail}`)).toBe('Europe/Berlin')
  })

  it('returns null for long unrelated jars without moana_tz= (prefix fast path)', () => {
    const junk = Array.from({ length: 40 }, (_, i) => `k${i}=v`).join('; ')
    expect(parseMoanaTimezoneCookie(junk)).toBeNull()
  })

  it('does not treat xmoana_tz as moana_tz', () => {
    expect(parseMoanaTimezoneCookie('xmoana_tz=1')).toBeNull()
  })
})

describe('setBrowserTimezoneCookie', () => {
  afterEach(() => {
    vi.restoreAllMocks()
    vi.unstubAllGlobals()
  })

  function mockResolvedTimeZone(tz: string) {
    vi.spyOn(Intl.DateTimeFormat.prototype, 'resolvedOptions').mockReturnValue({
      calendar: 'gregory',
      locale: 'en-US',
      numberingSystem: 'latn',
      timeZone: tz,
    } as Intl.ResolvedDateTimeFormatOptions)
  }

  it('writes moana_tz when missing or different from the browser zone', () => {
    mockResolvedTimeZone('Pacific/Auckland')
    let jar = ''
    vi.stubGlobal('document', {
      get cookie() {
        return jar
      },
      set cookie(v: string) {
        jar = v
      },
    })

    setBrowserTimezoneCookie()

    expect(jar).toContain(`${TIMEZONE_COOKIE_NAME}=`)
    expect(jar).toContain(encodeURIComponent('Pacific/Auckland'))
  })

  it('does not overwrite document.cookie when moana_tz already matches', () => {
    mockResolvedTimeZone('Europe/Berlin')
    const existing = `${TIMEZONE_COOKIE_NAME}=${encodeURIComponent('Europe/Berlin')}`
    let jar = existing
    let sets = 0
    vi.stubGlobal('document', {
      get cookie() {
        return jar
      },
      set cookie(v: string) {
        sets += 1
        jar = v
      },
    })

    setBrowserTimezoneCookie()

    expect(sets).toBe(0)
    expect(jar).toBe(existing)
  })

  it('does not write twice when run twice and moana_tz already matches', () => {
    mockResolvedTimeZone('Europe/Berlin')
    const existing = `${TIMEZONE_COOKIE_NAME}=${encodeURIComponent('Europe/Berlin')}`
    let jar = existing
    let sets = 0
    vi.stubGlobal('document', {
      get cookie() {
        return jar
      },
      set cookie(v: string) {
        sets += 1
        jar = v
      },
    })

    setBrowserTimezoneCookie()
    setBrowserTimezoneCookie()

    expect(sets).toBe(0)
    expect(jar).toBe(existing)
  })

  it('no-ops when the runtime reports an empty time zone', () => {
    mockResolvedTimeZone('')
    let sets = 0
    vi.stubGlobal('document', {
      get cookie() {
        return ''
      },
      set cookie(_v: string) {
        sets += 1
      },
    })

    setBrowserTimezoneCookie()

    expect(sets).toBe(0)
  })

  it('no-ops when resolvedOptions throws (try/catch in setBrowserTimezoneCookie)', () => {
    vi.spyOn(Intl.DateTimeFormat.prototype, 'resolvedOptions').mockImplementation(() => {
      throw new Error('unavailable')
    })
    let sets = 0
    vi.stubGlobal('document', {
      get cookie() {
        return ''
      },
      set cookie(_v: string) {
        sets += 1
      },
    })

    expect(() => setBrowserTimezoneCookie()).not.toThrow()
    expect(sets).toBe(0)
  })
})
