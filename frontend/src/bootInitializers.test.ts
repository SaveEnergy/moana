import { describe, expect, it } from 'vitest'

/**
 * No `vi.mock` for lib modules — `boot.test.ts` mocks them, which replaces initializers with spies
 * named `spy` and would hide accidental reorder in `BOOT_APP_INITIALIZERS`.
 */
import { BOOT_APP_INITIALIZERS, BOOT_INITIALIZER_NAMES } from './bootInitializers'

/** Must match `design.md` §2 JS row and **Boot content root** table — update when changing {@link BOOT_APP_INITIALIZERS}. */
const DOCUMENTED_BOOT_NAME_ORDER = [
  'setBrowserTimezoneCookie',
  'applyLocalTimeElements',
  'initShellSidebar',
  'initSettingsMemberDialog',
  'initCategoryModal',
  'initHistoryControls',
  'initConfirmSubmitForms',
] as const

describe('BOOT_APP_INITIALIZERS', () => {
  it('matches design.md sequence and derived BOOT_INITIALIZER_NAMES', () => {
    expect(BOOT_APP_INITIALIZERS.length).toBe(DOCUMENTED_BOOT_NAME_ORDER.length)
    expect([...BOOT_INITIALIZER_NAMES]).toEqual([...DOCUMENTED_BOOT_NAME_ORDER])
    expect(BOOT_APP_INITIALIZERS.map((fn) => fn.name)).toEqual([...DOCUMENTED_BOOT_NAME_ORDER])
  })
})
