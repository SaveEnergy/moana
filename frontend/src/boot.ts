import { withBootContentRootScope } from './lib/contentRoot'
import { type BootInitializer, BOOT_APP_INITIALIZERS } from './bootInitializers'

export type { BootInitializer } from './bootInitializers'
export { BOOT_APP_INITIALIZERS, DOCUMENTED_BOOT_INITIALIZER_NAMES } from './bootInitializers'

/**
 * Run an initializer list in order (indexed loop — no iterator allocation on duplicate **`bootApp()`**).
 * Non-empty lists run inside {@link withBootContentRootScope} so each step shares one resolved boot content root.
 * **`bootApp`** passes {@link BOOT_APP_INITIALIZERS}; tests may pass shorter arrays.
 */
export function runBootInitializers(initializers: ReadonlyArray<BootInitializer>): void {
  if (initializers.length === 0) {
    return
  }
  withBootContentRootScope(() => {
    for (let i = 0, n = initializers.length; i < n; i++) {
      initializers[i]()
    }
  })
}

/**
 * Wire all client behaviors; each initializer no-ops when its DOM is missing.
 * Order: timezone cookie → local time labels → shell (listeners) → settings add-member dialog → settings avatar dialog → category modal → history controls → confirm-before-submit forms.
 * Cookie and time text run before interactive modules attach handlers.
 * Name order: {@link DOCUMENTED_BOOT_INITIALIZER_NAMES} and `design.md` §2 **Boot content root**.
 */
export function bootApp(): void {
  runBootInitializers(BOOT_APP_INITIALIZERS)
}
