/**
 * Re-export for modules that already import this path (`boot.test.ts`, docs).
 * {@link BOOT_INITIALIZER_NAMES}, {@link BOOT_INITIALIZER_COUNT}, and {@link DOCUMENTED_BOOT_INITIALIZER_NAMES}
 * live alongside {@link BOOT_APP_INITIALIZERS} in `bootInitializers.ts`.
 */
export {
  BOOT_INITIALIZER_COUNT,
  BOOT_INITIALIZER_NAMES,
  DOCUMENTED_BOOT_INITIALIZER_NAMES,
} from './bootInitializers'
