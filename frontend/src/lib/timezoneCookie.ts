export const TIMEZONE_COOKIE_NAME = 'moana_tz'

const TZ_MAX_AGE_SEC = 365 * 24 * 60 * 60

/** Cookie `document.cookie` value segment for one IANA zone (deterministic; unit-tested). */
export function timezoneCookieSegment(tz: string): string {
  return `${TIMEZONE_COOKIE_NAME}=${encodeURIComponent(tz)}; Path=/; Max-Age=${TZ_MAX_AGE_SEC}; SameSite=Lax`
}

/** Persist browser IANA timezone for server-side date handling. */
export function setBrowserTimezoneCookie(): void {
  try {
    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
    if (!tz) return
    document.cookie = timezoneCookieSegment(tz)
  } catch {
    // ignore
  }
}
