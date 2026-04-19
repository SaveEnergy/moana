import { afterEach, describe, expect, it, vi } from 'vitest'

import { HISTORY_SORT_ELEMENT_ID } from './domSelectors'
import { initHistoryControls, wireHistorySortAutoSubmit } from './historyControls'

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

describe('initHistoryControls', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('no-ops when #history-sort is absent', () => {
    vi.stubGlobal('document', {
      getElementById: () => null,
    })
    expect(() => initHistoryControls()).not.toThrow()
  })

  it('delegates to wireHistorySortAutoSubmit when select is present', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement

    vi.stubGlobal('document', {
      getElementById: (id: string) => (id === HISTORY_SORT_ELEMENT_ID ? select : null),
    })

    initHistoryControls()

    expect(addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    const onChange = addEventListener.mock.calls[0][1] as () => void
    onChange()
    expect(requestSubmit).toHaveBeenCalledTimes(1)
  })
})
