import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('user menu navigates to settings', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/')
  await page.locator('details.app-user-menu summary.app-user-menu-btn').click()
  await page.getByRole('link', { name: 'Settings' }).click()
  await expect(page).toHaveURL(/\/settings/)
  await expect(page.getByRole('heading', { name: 'Personal' })).toBeVisible()
})
