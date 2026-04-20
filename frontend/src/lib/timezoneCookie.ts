export const TIMEZONE_COOKIE_NAME = 'moana_tz'

const TZ_MAX_AGE_SEC = 365 * 24 * 60 * 60

/** Cookie `document.cookie` value segment for one IANA zone (deterministic; unit-tested). */
export function timezoneCookieSegment(tz: string): string {
  return `${TIMEZONE_COOKIE_NAME}=${encodeURIComponent(tz)}; Path=/; Max-Age=${TZ_MAX_AGE_SEC}; SameSite=Lax`
}

/**
 * Read `moana_tz` from a `document.cookie`-style string (testable without JSDOM).
 * Returns null if missing or not decodable; an empty value after `=` yields `""` (valid `decodeURIComponent`).
 * Scans semicolon-separated segments without `split` so large `document.cookie` strings avoid an extra array allocation.
 */
export function parseMoanaTimezoneCookie(cookieHeader: string): string | null {
  const prefix = `${TIMEZONE_COOKIE_NAME}=`
  /* One linear scan — skip segment walk when the name never appears (common on first visit / unrelated cookies). */
  if (!cookieHeader.includes(prefix)) {
    return null
  }
  let pos = 0
  while (pos < cookieHeader.length) {
    const semi = cookieHeader.indexOf(';', pos)
    const end = semi === -1 ? cookieHeader.length : semi
    const s = cookieHeader.slice(pos, end).trim()
    if (s.startsWith(prefix)) {
      const raw = s.slice(prefix.length)
      try {
        return decodeURIComponent(raw)
      } catch {
        // Malformed segment — try another `moana_tz=` later in the header.
      }
    }
    if (semi === -1) {
      break
    }
    pos = semi + 1
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
