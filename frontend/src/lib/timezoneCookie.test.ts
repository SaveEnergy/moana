import { describe, it, expect } from 'vitest'

import {
  TIMEZONE_COOKIE_NAME,
  parseMoanaTimezoneCookie,
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

  it('returns null on invalid percent encoding', () => {
    expect(parseMoanaTimezoneCookie(`${TIMEZONE_COOKIE_NAME}=%`)).toBeNull()
  })
})
