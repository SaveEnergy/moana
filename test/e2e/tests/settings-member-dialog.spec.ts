import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import {
  APP_SHELL,
  SETTINGS_ADD_MEMBER_CANCEL,
  SETTINGS_ADD_MEMBER_DIALOG,
  SETTINGS_ADD_MEMBER_DIALOG_HEADER,
  SETTINGS_ADD_MEMBER_DIALOG_ID,
  SETTINGS_ADD_MEMBER_OPEN,
  SETTINGS_ADD_MEMBER_TITLE,
} from '../helpers/shellSelectors'

test('add-member dialog opens and dismisses', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await page.locator(SETTINGS_ADD_MEMBER_OPEN).click()
  const dlg = page.locator(SETTINGS_ADD_MEMBER_DIALOG)
  await expect(dlg).toBeVisible()
  await expect(page.locator(SETTINGS_ADD_MEMBER_TITLE)).toHaveText('Add household member')
  await page.locator(SETTINGS_ADD_MEMBER_CANCEL).click()
  await expect(dlg).toBeHidden()
})

test('add-member dialog closes on Escape', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await page.locator(SETTINGS_ADD_MEMBER_OPEN).click()
  const dlg = page.locator(SETTINGS_ADD_MEMBER_DIALOG)
  await expect(dlg).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dlg).toBeHidden()
})

test('add-member dialog closes on backdrop click outside card', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await page.locator(SETTINGS_ADD_MEMBER_OPEN).click()
  const dlg = page.locator(SETTINGS_ADD_MEMBER_DIALOG)
  await expect(dlg).toBeVisible()
  const header = page.locator(SETTINGS_ADD_MEMBER_DIALOG_HEADER)
  const box = await header.boundingBox()
  expect(box).not.toBeNull()
  await page.mouse.click(box!.x + box!.width / 2, Math.max(8, box!.y - 24))
  await expect(dlg).toBeHidden()
})

test('mobile: first Escape closes add-member dialog with sidebar left open', async ({ page }) => {
  await signInAsTestUser(page)
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/settings')
  const shell = page.locator(APP_SHELL)
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.evaluate((id) => {
    document.getElementById(id)?.showModal()
  }, SETTINGS_ADD_MEMBER_DIALOG_ID)
  const dlg = page.locator(SETTINGS_ADD_MEMBER_DIALOG)
  await expect(dlg).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dlg).toBeHidden()
  await expect(shell).toHaveClass(/sidebar-open/)
})
