import { describe, expect, it, vi } from 'vitest'

import { wireHistorySortAutoSubmit } from './historyControls'

describe('wireHistorySortAutoSubmit', () => {
  it('no-ops on null', () => {
    wireHistorySortAutoSubmit(null)
  })

  it('requests submit on change when form exists', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement
    wireHistorySortAutoSubmit(select)
    expect(addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    const onChange = addEventListener.mock.calls[0][1] as () => void
    onChange()
    expect(requestSubmit).toHaveBeenCalledTimes(1)
  })

  it('does not subscribe when form is missing', () => {
    const addEventListener = vi.fn()
    const select = { addEventListener, form: null } as unknown as HTMLSelectElement
    wireHistorySortAutoSubmit(select)
    expect(addEventListener).not.toHaveBeenCalled()
  })
})
