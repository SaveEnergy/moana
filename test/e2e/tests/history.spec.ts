import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('history page loads shell and heading', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history')
  await expect(page).toHaveURL(/\/history/)
  await expect(page.getByRole('heading', { name: 'History' })).toBeVisible()
  await expect(page.getByRole('tablist', { name: 'Filter by type' })).toBeVisible()
})
