import { describe, expect, it, vi } from 'vitest'

import { attachShowModalOnClick, closeDialogIfOpen, showModalIfClosed } from './dialogModal'

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

describe('attachShowModalOnClick', () => {
  it('registers click that calls showModalIfClosed', () => {
    const showModal = vi.fn()
    const dialog = { open: false, showModal } as unknown as HTMLDialogElement
    const handlers: Array<() => void> = []
    const openBtn = {
      addEventListener: (_type: string, fn: () => void) => {
        handlers.push(fn)
      },
    } as unknown as HTMLElement

    attachShowModalOnClick(openBtn, dialog)

    expect(handlers).toHaveLength(1)
    handlers[0]()
    expect(showModal).toHaveBeenCalledTimes(1)
  })

  it('click handler skips showModal when dialog is already open', () => {
    const showModal = vi.fn()
    const dialog = { open: true, showModal } as unknown as HTMLDialogElement
    const handlers: Array<() => void> = []
    const openBtn = {
      addEventListener: (_type: string, fn: () => void) => {
        handlers.push(fn)
      },
    } as unknown as HTMLElement

    attachShowModalOnClick(openBtn, dialog)
    handlers[0]()

    expect(showModal).not.toHaveBeenCalled()
  })

  it('no-ops when open button is null', () => {
    const dialog = { open: false, showModal: vi.fn() } as unknown as HTMLDialogElement
    expect(() => attachShowModalOnClick(null, dialog)).not.toThrow()
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
