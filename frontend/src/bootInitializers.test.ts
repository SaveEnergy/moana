import { describe, expect, it } from 'vitest'

/**
 * No `vi.mock` for lib modules — `boot.test.ts` mocks them, which replaces initializers with spies
 * named `spy` and would hide accidental reorder in `BOOT_APP_INITIALIZERS`.
 */
import {
  BOOT_APP_INITIALIZERS as BOOT_FROM_BOOT,
  DOCUMENTED_BOOT_INITIALIZER_NAMES as DOCUMENTED_FROM_BOOT,
} from './boot'
import {
  BOOT_APP_INITIALIZERS,
  BOOT_INITIALIZER_NAMES,
  DOCUMENTED_BOOT_INITIALIZER_NAMES,
} from './bootInitializers'

describe('BOOT_APP_INITIALIZERS', () => {
  it('matches DOCUMENTED_BOOT_INITIALIZER_NAMES and derived BOOT_INITIALIZER_NAMES (design.md §2)', () => {
    expect(BOOT_APP_INITIALIZERS.length).toBe(DOCUMENTED_BOOT_INITIALIZER_NAMES.length)
    expect([...BOOT_INITIALIZER_NAMES]).toEqual([...DOCUMENTED_BOOT_INITIALIZER_NAMES])
    expect(BOOT_APP_INITIALIZERS.map((fn) => fn.name)).toEqual([...DOCUMENTED_BOOT_INITIALIZER_NAMES])
  })

  it('boot.ts re-exports DOCUMENTED_BOOT_INITIALIZER_NAMES by reference (stable for tooling imports)', () => {
    expect(DOCUMENTED_FROM_BOOT).toBe(DOCUMENTED_BOOT_INITIALIZER_NAMES)
  })

  it('boot.ts re-exports BOOT_APP_INITIALIZERS by reference (same array as bootInitializers.ts)', () => {
    expect(BOOT_FROM_BOOT).toBe(BOOT_APP_INITIALIZERS)
  })
})
