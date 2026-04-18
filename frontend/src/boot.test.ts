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
})
