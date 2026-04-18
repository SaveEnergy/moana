import { describe, it, expect } from 'vitest'

import { TIMEZONE_COOKIE_NAME, timezoneCookieSegment } from './timezoneCookie'

describe('timezoneCookieSegment', () => {
  it('includes encoded tz and expected attrs', () => {
    const s = timezoneCookieSegment('Europe/Berlin')
    expect(s.startsWith(`${TIMEZONE_COOKIE_NAME}=`)).toBe(true)
    expect(s).toContain(encodeURIComponent('Europe/Berlin'))
    expect(s).toContain('Path=/')
    expect(s).toContain('SameSite=Lax')
  })
})
