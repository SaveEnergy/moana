import { afterEach, describe, expect, it, vi } from 'vitest'

import { createCategoryModalOpenFlow } from './categoryModalOpenFlow'
import type { CategoryModalPreviewController } from './categoryModalPreview'
import * as categoryModalDataset from './categoryModalDataset'
import * as dialogModal from './dialogModal'

describe('createCategoryModalOpenFlow', () => {
  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('openCreateModal resets form, targets create action, and shows modal', () => {
    vi.spyOn(dialogModal, 'showModalIfClosed').mockImplementation(() => {})

    const form = {
      action: '',
      reset: vi.fn(),
    } as unknown as HTMLFormElement
    const idInput = { value: 'old' } as HTMLInputElement
    const titleEl = { textContent: '' } as HTMLElement
    const submitBtn = { textContent: '' } as HTMLElement
    const nameInput = { value: '', focus: vi.fn() } as unknown as HTMLInputElement
    const colorNativeInput = { value: '' } as HTMLInputElement
    const modal = {} as HTMLDialogElement

    const previewCtl: Pick<CategoryModalPreviewController, 'raf' | 'resetPaintState' | 'sync'> = {
      raf: { cancelPending: vi.fn(), schedule: vi.fn() },
      resetPaintState: vi.fn(),
      sync: vi.fn(),
    }

    const { openCreateModal } = createCategoryModalOpenFlow({
      modal,
      form,
      idInput,
      titleEl,
      submitBtn,
      nameInput,
      colorNativeInput,
      colorRadioByValue: new Map(),
      iconRadioByValue: new Map(),
      previewCtl: previewCtl as CategoryModalPreviewController,
    })

    openCreateModal()

    expect(form.action).toBe('/categories')
    expect(idInput.value).toBe('')
    expect(titleEl.textContent).toBe('New category')
    expect(submitBtn.textContent).toBe('Create category')
    expect(form.reset).toHaveBeenCalledTimes(1)
    expect(previewCtl.sync).toHaveBeenCalledTimes(1)
    expect(nameInput.focus).toHaveBeenCalledTimes(1)
    expect(dialogModal.showModalIfClosed).toHaveBeenCalledWith(modal)
  })

  it('openEditModal no-ops when dataset row is invalid', () => {
    vi.spyOn(dialogModal, 'showModalIfClosed').mockImplementation(() => {})
    vi.spyOn(categoryModalDataset, 'readCategoryEditRowDataset').mockReturnValue(null)

    const form = { action: '/categories' } as HTMLFormElement
    const previewCtl: Pick<CategoryModalPreviewController, 'raf' | 'resetPaintState' | 'sync'> = {
      raf: { cancelPending: vi.fn(), schedule: vi.fn() },
      resetPaintState: vi.fn(),
      sync: vi.fn(),
    }

    const { openEditModal } = createCategoryModalOpenFlow({
      modal: {} as HTMLDialogElement,
      form,
      idInput: { value: '' } as HTMLInputElement,
      titleEl: { textContent: '' } as HTMLElement,
      submitBtn: { textContent: '' } as HTMLElement,
      nameInput: { value: '', focus: vi.fn() } as unknown as HTMLInputElement,
      colorNativeInput: null,
      colorRadioByValue: new Map(),
      iconRadioByValue: new Map(),
      previewCtl: previewCtl as CategoryModalPreviewController,
    })

    openEditModal({} as HTMLElement)

    expect(form.action).toBe('/categories')
    expect(previewCtl.sync).not.toHaveBeenCalled()
    expect(dialogModal.showModalIfClosed).not.toHaveBeenCalled()
  })
})
