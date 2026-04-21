import { type BootInitializer, BOOT_APP_INITIALIZERS } from './bootInitializers'

export type { BootInitializer } from './bootInitializers'
export { BOOT_APP_INITIALIZERS } from './bootInitializers'

/**
 * Run an initializer list in order (indexed loop — no iterator allocation on duplicate **`bootApp()`**).
 * **`bootApp`** passes {@link BOOT_APP_INITIALIZERS}; tests may pass shorter arrays.
 */
export function runBootInitializers(initializers: ReadonlyArray<BootInitializer>): void {
  for (let i = 0, n = initializers.length; i < n; i++) {
    initializers[i]()
  }
}

/**
 * Wire all client behaviors; each initializer no-ops when its DOM is missing.
 * Order: timezone cookie → local time labels → shell (listeners) → settings dialog → category modal → history controls → confirm-before-submit forms.
 * Cookie and time text run before interactive modules attach handlers.
 */
export function bootApp(): void {
  runBootInitializers(BOOT_APP_INITIALIZERS)
}
