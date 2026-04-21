import { beforeEach, describe, expect, it, vi } from 'vitest'

const stubs = vi.hoisted(() => ({
  applyLocalTimeElements: vi.fn(),
  setBrowserTimezoneCookie: vi.fn(),
  initShellSidebar: vi.fn(),
  initSettingsMemberDialog: vi.fn(),
  initCategoryModal: vi.fn(),
  initHistoryControls: vi.fn(),
  initConfirmSubmitForms: vi.fn(),
}))

vi.mock('./lib/localTime', () => ({ applyLocalTimeElements: stubs.applyLocalTimeElements }))
vi.mock('./lib/timezoneCookie', () => ({ setBrowserTimezoneCookie: stubs.setBrowserTimezoneCookie }))
vi.mock('./lib/shellSidebar', () => ({ initShellSidebar: stubs.initShellSidebar }))
vi.mock('./lib/settingsMemberDialog', () => ({ initSettingsMemberDialog: stubs.initSettingsMemberDialog }))
vi.mock('./lib/categoryModal', () => ({ initCategoryModal: stubs.initCategoryModal }))
vi.mock('./lib/historyControls', () => ({ initHistoryControls: stubs.initHistoryControls }))
vi.mock('./lib/confirmSubmitForms', () => ({ initConfirmSubmitForms: stubs.initConfirmSubmitForms }))

import { BOOT_APP_INITIALIZERS, bootApp } from './boot'
import { BOOT_APP_INITIALIZERS as bootInitializersModuleList } from './bootInitializers'
import { BOOT_INITIALIZER_COUNT, BOOT_INITIALIZER_NAMES } from './bootInitializerNames'

describe('bootApp', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('BOOT_APP_INITIALIZERS lists one entry per boot step', () => {
    expect(BOOT_INITIALIZER_COUNT).toBe(BOOT_INITIALIZER_NAMES.length)
    expect(BOOT_APP_INITIALIZERS).toHaveLength(BOOT_INITIALIZER_COUNT)
  })

  it('re-exports the same initializer array as bootInitializers.ts (stable module boundary)', () => {
    expect(BOOT_APP_INITIALIZERS).toBe(bootInitializersModuleList)
  })

  it('invokes each initializer once', () => {
    bootApp()
    expect(stubs.setBrowserTimezoneCookie).toHaveBeenCalledTimes(1)
    expect(stubs.applyLocalTimeElements).toHaveBeenCalledTimes(1)
    expect(stubs.initShellSidebar).toHaveBeenCalledTimes(1)
    expect(stubs.initSettingsMemberDialog).toHaveBeenCalledTimes(1)
    expect(stubs.initCategoryModal).toHaveBeenCalledTimes(1)
    expect(stubs.initHistoryControls).toHaveBeenCalledTimes(1)
    expect(stubs.initConfirmSubmitForms).toHaveBeenCalledTimes(1)
  })

  it('invokes each initializer once per bootApp call (orchestrator is not implicitly single-run)', () => {
    bootApp()
    bootApp()
    expect(stubs.setBrowserTimezoneCookie).toHaveBeenCalledTimes(2)
    expect(stubs.applyLocalTimeElements).toHaveBeenCalledTimes(2)
    expect(stubs.initShellSidebar).toHaveBeenCalledTimes(2)
    expect(stubs.initSettingsMemberDialog).toHaveBeenCalledTimes(2)
    expect(stubs.initCategoryModal).toHaveBeenCalledTimes(2)
    expect(stubs.initHistoryControls).toHaveBeenCalledTimes(2)
    expect(stubs.initConfirmSubmitForms).toHaveBeenCalledTimes(2)
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
    stubs.initHistoryControls.mockImplementation(() => {
      order.push('historyControls')
    })
    stubs.initConfirmSubmitForms.mockImplementation(() => {
      order.push('confirmSubmitForms')
    })
    bootApp()
    expect(order).toEqual([
      'timezone',
      'localTime',
      'shell',
      'settingsDialog',
      'categoryModal',
      'historyControls',
      'confirmSubmitForms',
    ])
  })
})
