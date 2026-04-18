export const TIMEZONE_COOKIE_NAME = 'moana_tz'

const TZ_MAX_AGE_SEC = 365 * 24 * 60 * 60

/** Cookie `document.cookie` value segment for one IANA zone (deterministic; unit-tested). */
export function timezoneCookieSegment(tz: string): string {
  return `${TIMEZONE_COOKIE_NAME}=${encodeURIComponent(tz)}; Path=/; Max-Age=${TZ_MAX_AGE_SEC}; SameSite=Lax`
}

/**
 * Read `moana_tz` from a `document.cookie`-style string (testable without JSDOM).
 * Returns null if missing or not decodable.
 */
export function parseMoanaTimezoneCookie(cookieHeader: string): string | null {
  const prefix = `${TIMEZONE_COOKIE_NAME}=`
  for (const part of cookieHeader.split(';')) {
    const s = part.trim()
    if (s.startsWith(prefix)) {
      const raw = s.slice(prefix.length)
      try {
        return decodeURIComponent(raw)
      } catch {
        return null
      }
    }
  }
  return null
}

/** Persist browser IANA timezone for server-side date handling. */
export function setBrowserTimezoneCookie(): void {
  try {
    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
    if (!tz) return
    if (parseMoanaTimezoneCookie(document.cookie) === tz) {
      return
    }
    document.cookie = timezoneCookieSegment(tz)
  } catch {
    // ignore
  }
}
