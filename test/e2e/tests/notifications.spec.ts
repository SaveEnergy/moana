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

test('dashboard notification bell uses layout hook and hides badge when inbox empty', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/')
  const bell = page.locator('a.app-topbar-notif-btn[href="/notifications"]')
  await expect(bell).toBeVisible()
  await expect(bell.locator('.app-notif-badge')).toHaveCount(0)
})

test('notifications page sets topbar notifications link as current', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/notifications')
  await expect(page.locator('a.app-topbar-notif-btn[href="/notifications"]')).toHaveAttribute(
    'aria-current',
    'page',
  )
})
