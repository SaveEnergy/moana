import { afterEach, describe, expect, it, vi } from 'vitest'

import { HISTORY_SORT_SELECTOR } from './domSelectors'
import { initHistoryControls, queryHistorySortSelect, wireHistorySortAutoSubmit } from './historyControls'
import { stubDocumentMainLandmark, stubDocumentWithoutMainLandmark } from './stubDocumentMainLandmark'

describe('queryHistorySortSelect', () => {
  it('uses HISTORY_SORT_SELECTOR on root', () => {
    let seen = ''
    const root = {
      querySelector: (sel: string) => {
        seen = sel
        return null
      },
    } as unknown as ParentNode
    queryHistorySortSelect(root)
    expect(seen).toBe(HISTORY_SORT_SELECTOR)
  })
})

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

  it('does not subscribe twice when wiring runs twice on the same select', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement
    wireHistorySortAutoSubmit(select)
    wireHistorySortAutoSubmit(select)
    expect(addEventListener).toHaveBeenCalledTimes(1)
  })
})

describe('initHistoryControls', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('no-ops when #history-sort is absent', () => {
    vi.stubGlobal('document', stubDocumentWithoutMainLandmark())
    expect(() => initHistoryControls()).not.toThrow()
  })

  it('delegates to wireHistorySortAutoSubmit when select is present', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelector: (sel: string) => (sel === HISTORY_SORT_SELECTOR ? select : null),
      }),
    )

    initHistoryControls()

    expect(addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    const onChange = addEventListener.mock.calls[0][1] as () => void
    onChange()
    expect(requestSubmit).toHaveBeenCalledTimes(1)
  })

  it('does not stack change listeners when initHistoryControls runs twice', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelector: (sel: string) => (sel === HISTORY_SORT_SELECTOR ? select : null),
      }),
    )

    initHistoryControls()
    initHistoryControls()

    expect(addEventListener).toHaveBeenCalledTimes(1)
  })

  it('resolves #history-sort under main.app-main when the landmark exists', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement
    const main = {
      querySelector: (sel: string) => (sel === HISTORY_SORT_SELECTOR ? select : null),
    } as unknown as ParentNode

    vi.stubGlobal('document', stubDocumentMainLandmark(main))

    initHistoryControls()

    expect(addEventListener).toHaveBeenCalledWith('change', expect.any(Function))
    const onChange = addEventListener.mock.calls[0][1] as () => void
    onChange()
    expect(requestSubmit).toHaveBeenCalledTimes(1)
  })

  it('does not re-wire when wireHistorySortAutoSubmit ran before init', () => {
    const requestSubmit = vi.fn()
    const form = { requestSubmit } as unknown as HTMLFormElement
    const addEventListener = vi.fn()
    const select = { addEventListener, form } as unknown as HTMLSelectElement

    vi.stubGlobal(
      'document',
      stubDocumentWithoutMainLandmark({
        querySelector: (sel: string) => (sel === HISTORY_SORT_SELECTOR ? select : null),
      }),
    )

    wireHistorySortAutoSubmit(select)
    addEventListener.mockClear()
    initHistoryControls()

    expect(addEventListener).not.toHaveBeenCalled()
  })
})
