const TZ_COOKIE = 'moana_tz'
const TZ_MAX_AGE_SEC = 365 * 24 * 60 * 60

/** Persist browser IANA timezone for server-side date handling. */
export function setBrowserTimezoneCookie(): void {
  try {
    const tz = Intl.DateTimeFormat().resolvedOptions().timeZone
    if (!tz) return
    document.cookie = `${TZ_COOKIE}=${encodeURIComponent(tz)}; Path=/; Max-Age=${TZ_MAX_AGE_SEC}; SameSite=Lax`
  } catch {
    // ignore
  }
}
