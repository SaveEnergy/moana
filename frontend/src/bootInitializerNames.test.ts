import { describe, expect, it } from 'vitest'

import { BOOT_INITIALIZER_NAMES, DOCUMENTED_BOOT_INITIALIZER_NAMES } from './bootInitializerNames'

describe('BOOT_INITIALIZER_NAMES', () => {
  it('matches DOCUMENTED_BOOT_INITIALIZER_NAMES (re-export parity with bootInitializers.test)', () => {
    expect([...DOCUMENTED_BOOT_INITIALIZER_NAMES]).toEqual([...BOOT_INITIALIZER_NAMES])
  })

  it('has no duplicate entries (boot order must be unambiguous)', () => {
    const unique = new Set(BOOT_INITIALIZER_NAMES)
    expect(unique.size).toBe(BOOT_INITIALIZER_NAMES.length)
  })

  it('lists only non-empty string names', () => {
    for (const name of BOOT_INITIALIZER_NAMES) {
      expect(name.length).toBeGreaterThan(0)
      expect(name.trim()).toBe(name)
    }
  })
})
