import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('notifications inbox empty state', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/notifications')
  await expect(page).toHaveURL(/\/notifications/)
  await expect(page.getByRole('heading', { name: 'Notifications' })).toBeVisible()
  await expect(
    page.getByText('Account and activity alerts will appear here.'),
  ).toBeVisible()
  await expect(page.getByText('You have no notifications.')).toBeVisible()
})

test('notifications page sets topbar notifications link as current', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/notifications')
  await expect(page.locator('a.app-topbar-icon-btn[href="/notifications"]')).toHaveAttribute(
    'aria-current',
    'page',
  )
})
