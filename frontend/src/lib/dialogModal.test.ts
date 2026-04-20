import { describe, expect, it, vi } from 'vitest'

import { showModalIfClosed } from './dialogModal'

describe('showModalIfClosed', () => {
  it('calls showModal when open is false', () => {
    const showModal = vi.fn()
    const dialog = { open: false, showModal } as unknown as HTMLDialogElement
    showModalIfClosed(dialog)
    expect(showModal).toHaveBeenCalledTimes(1)
  })

  it('skips showModal when already open', () => {
    const showModal = vi.fn()
    const dialog = { open: true, showModal } as unknown as HTMLDialogElement
    showModalIfClosed(dialog)
    expect(showModal).not.toHaveBeenCalled()
  })

  it('calls showModal when open is undefined (stub without property)', () => {
    const showModal = vi.fn()
    const dialog = { showModal } as unknown as HTMLDialogElement
    showModalIfClosed(dialog)
    expect(showModal).toHaveBeenCalledTimes(1)
  })
})
