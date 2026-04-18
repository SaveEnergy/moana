import { beforeEach, describe, expect, it, vi } from 'vitest'

const stubs = vi.hoisted(() => ({
  applyLocalTimeElements: vi.fn(),
  setBrowserTimezoneCookie: vi.fn(),
  initShellSidebar: vi.fn(),
  initSettingsMemberDialog: vi.fn(),
  initCategoryModal: vi.fn(),
}))

vi.mock('./lib/localTime', () => ({ applyLocalTimeElements: stubs.applyLocalTimeElements }))
vi.mock('./lib/timezoneCookie', () => ({ setBrowserTimezoneCookie: stubs.setBrowserTimezoneCookie }))
vi.mock('./lib/shellSidebar', () => ({ initShellSidebar: stubs.initShellSidebar }))
vi.mock('./lib/settingsMemberDialog', () => ({ initSettingsMemberDialog: stubs.initSettingsMemberDialog }))
vi.mock('./lib/categoryModal', () => ({ initCategoryModal: stubs.initCategoryModal }))

import { bootApp } from './boot'

describe('bootApp', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('invokes each initializer once', () => {
    bootApp()
    expect(stubs.setBrowserTimezoneCookie).toHaveBeenCalledTimes(1)
    expect(stubs.applyLocalTimeElements).toHaveBeenCalledTimes(1)
    expect(stubs.initShellSidebar).toHaveBeenCalledTimes(1)
    expect(stubs.initSettingsMemberDialog).toHaveBeenCalledTimes(1)
    expect(stubs.initCategoryModal).toHaveBeenCalledTimes(1)
  })

  it('runs initializers in a stable order', () => {
    const order: string[] = []
    stubs.setBrowserTimezoneCookie.mockImplementation(() => {
      order.push('timezone')
    })
    stubs.applyLocalTimeElements.mockImplementation(() => {
      order.push('localTime')
    })
    stubs.initShellSidebar.mockImplementation(() => {
      order.push('shell')
    })
    stubs.initSettingsMemberDialog.mockImplementation(() => {
      order.push('settingsDialog')
    })
    stubs.initCategoryModal.mockImplementation(() => {
      order.push('categoryModal')
    })
    bootApp()
    expect(order).toEqual([
      'timezone',
      'localTime',
      'shell',
      'settingsDialog',
      'categoryModal',
    ])
  })
})
