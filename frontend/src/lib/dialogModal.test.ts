import { describe, expect, it, vi } from 'vitest'

import { closeDialogIfOpen, showModalIfClosed } from './dialogModal'

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

describe('closeDialogIfOpen', () => {
  it('calls close when open is true', () => {
    const close = vi.fn()
    const dialog = { open: true, close } as unknown as HTMLDialogElement
    closeDialogIfOpen(dialog)
    expect(close).toHaveBeenCalledTimes(1)
  })

  it('skips close when open is false', () => {
    const close = vi.fn()
    const dialog = { open: false, close } as unknown as HTMLDialogElement
    closeDialogIfOpen(dialog)
    expect(close).not.toHaveBeenCalled()
  })

  it('calls close when open is undefined (stub without property)', () => {
    const close = vi.fn()
    const dialog = { close } as unknown as HTMLDialogElement
    closeDialogIfOpen(dialog)
    expect(close).toHaveBeenCalledTimes(1)
  })
})
