import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import {
  NOTIFICATIONS_PATH,
  TOPBAR_NOTIFICATION_BADGE,
  TOPBAR_NOTIFICATIONS_LINK,
} from '../helpers/shellSelectors'

test('notifications inbox empty state', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto(NOTIFICATIONS_PATH)
  await expect(page).toHaveURL(new RegExp(`${NOTIFICATIONS_PATH}$`))
  await expect(page.getByRole('heading', { name: 'Notifications' })).toBeVisible()
  await expect(
    page.getByText('Account and activity alerts will appear here.'),
  ).toBeVisible()
  await expect(page.getByText('You have no notifications.')).toBeVisible()
})

test('dashboard notification bell uses layout hook and hides badge when inbox empty', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/')
  const bell = page.locator(TOPBAR_NOTIFICATIONS_LINK)
  await expect(bell).toBeVisible()
  await expect(bell.locator(TOPBAR_NOTIFICATION_BADGE)).toHaveCount(0)
})

test('notifications page sets topbar notifications link as current', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto(NOTIFICATIONS_PATH)
  await expect(page.locator(TOPBAR_NOTIFICATIONS_LINK)).toHaveAttribute(
    'aria-current',
    'page',
  )
})
