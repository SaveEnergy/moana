import { BOOT_APP_INITIALIZERS } from './bootInitializers'

export type { BootInitializer } from './bootInitializers'
export { BOOT_APP_INITIALIZERS } from './bootInitializers'

/**
 * Wire all client behaviors; each initializer no-ops when its DOM is missing.
 * Order: timezone cookie → local time labels → shell (listeners) → settings dialog → category modal → history controls → confirm-before-submit forms.
 * Cookie and time text run before interactive modules attach handlers.
 */
export function bootApp(): void {
  for (const init of BOOT_APP_INITIALIZERS) {
    init()
  }
}
