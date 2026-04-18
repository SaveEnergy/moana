import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('sign out from user menu redirects to login', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/')
  await page.locator('details.app-user-menu summary.app-user-menu-btn').click()
  await page.getByRole('button', { name: 'Sign out' }).click()
  await expect(page).toHaveURL(/\/login/)
  await expect(page.getByRole('heading', { name: /sign in/i })).toBeVisible()
})

test('session gate: dashboard redirects anonymous users to login', async ({ page }) => {
  await page.goto('/')
  await expect(page).toHaveURL(/\/login/)
})
