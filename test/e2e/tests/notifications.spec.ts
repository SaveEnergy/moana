import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('notifications inbox empty state', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/notifications')
  await expect(page).toHaveURL(/\/notifications/)
  await expect(page.getByRole('heading', { name: 'Notifications' })).toBeVisible()
  await expect(page.getByText('You have no notifications.')).toBeVisible()
})
