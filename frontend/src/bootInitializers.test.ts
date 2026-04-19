import { describe, expect, it } from 'vitest'

/**
 * No `vi.mock` for lib modules — `boot.test.ts` mocks them, which replaces initializers with spies
 * named `spy` and would hide accidental reorder in `BOOT_APP_INITIALIZERS`.
 */
import { BOOT_APP_INITIALIZERS } from './boot'
import { BOOT_INITIALIZER_NAMES } from './bootInitializerNames'

describe('BOOT_APP_INITIALIZERS', () => {
  it('lists real initializers in the documented order (design.md)', () => {
    expect(BOOT_APP_INITIALIZERS.map((fn) => fn.name)).toEqual([...BOOT_INITIALIZER_NAMES])
  })
})
