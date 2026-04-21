import { trimEdgesIfNeeded } from './trimEdges'

export const TIMEZONE_COOKIE_NAME = 'moana_tz'

const TZ_MAX_AGE_SEC = 365 * 24 * 60 * 60

/** Lazy singleton for **`resolvedOptions().timeZone`** — duplicate **`bootApp()`** must not allocate another `Intl.DateTimeFormat`. */
let browserTimeZoneFormat: Intl.DateTimeFormat | undefined

/** Cookie `document.cookie` value segment for one IANA zone (deterministic; unit-tested). */
export function timezoneCookieSegment(tz: string): string {
  return `${TIMEZONE_COOKIE_NAME}=${encodeURIComponent(tz)}; Path=/; Max-Age=${TZ_MAX_AGE_SEC}; SameSite=Lax`
}

/**
 * Read `moana_tz` from a `document.cookie`-style string (testable without JSDOM).
 * Returns null if missing or not decodable; an empty value after `=` yields `""` (valid `decodeURIComponent`).
 * Each segment uses the first `=` as the name/value boundary with **trimmed** names (RFC-style spaces around `=`).
 * Scans semicolon-separated segments without `split` on the full header so large `document.cookie` avoids a segment array.
 */
export function parseMoanaTimezoneCookie(cookieHeader: string): string | null {
  /* Skip the segment walk when the cookie name cannot appear (first-visit / unrelated jars). */
  if (!cookieHeader.includes(TIMEZONE_COOKIE_NAME)) {
    return null
  }
  let pos = 0
  while (pos < cookieHeader.length) {
    const semi = cookieHeader.indexOf(';', pos)
    const end = semi === -1 ? cookieHeader.length : semi
    const s = trimEdgesIfNeeded(cookieHeader.slice(pos, end))
    const eq = s.indexOf('=')
    if (eq !== -1) {
      const segName = trimEdgesIfNeeded(s.slice(0, eq))
      if (segName === TIMEZONE_COOKIE_NAME) {
        const raw = trimEdgesIfNeeded(s.slice(eq + 1))
        try {
          return decodeURIComponent(raw)
        } catch {
          // Malformed segment — try another `moana_tz` later in the header.
        }
      }
    }
    if (semi === -1) {
      break
    }
    pos = semi + 1
  }
  return null
}

/** Persist browser IANA timezone for server-side date handling. Snapshots **`document.cookie`** once when comparing to the resolved zone. */
export function setBrowserTimezoneCookie(): void {
  try {
    if (!browserTimeZoneFormat) {
      browserTimeZoneFormat = new Intl.DateTimeFormat()
    }
    const tz = browserTimeZoneFormat.resolvedOptions().timeZone
    if (!tz) return
    const jar = document.cookie
    if (parseMoanaTimezoneCookie(jar) === tz) {
      return
    }
    document.cookie = timezoneCookieSegment(tz)
  } catch {
    // ignore
  }
}
