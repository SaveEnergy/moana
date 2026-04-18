import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('add-member dialog opens and dismisses', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await page.locator('#settings-add-member-open').click()
  const dlg = page.locator('#settings-add-member-dialog')
  await expect(dlg).toBeVisible()
  await expect(page.locator('#settings-add-member-title')).toHaveText('Add household member')
  await page.locator('#settings-add-member-cancel').click()
  await expect(dlg).toBeHidden()
})

test('add-member dialog closes on Escape', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await page.locator('#settings-add-member-open').click()
  const dlg = page.locator('#settings-add-member-dialog')
  await expect(dlg).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dlg).toBeHidden()
})
