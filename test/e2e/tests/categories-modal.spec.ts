import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('category modal opens and closes', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator('#cat-modal-open-create').click()
  const dialog = page.locator('#cat-modal')
  await expect(dialog).toBeVisible()
  await expect(page.locator('#cat-modal-title')).toHaveText('New category')
  await page.locator('#cat-modal-close').click()
  await expect(dialog).toBeHidden()
})

test('category modal closes on Escape', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator('#cat-modal-open-create').click()
  const dialog = page.locator('#cat-modal')
  await expect(dialog).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dialog).toBeHidden()
})
