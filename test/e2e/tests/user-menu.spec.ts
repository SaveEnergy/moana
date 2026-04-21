import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import { APP_USER_MENU_SUMMARY, USER_MENU_SETTINGS_LINK } from '../helpers/shellSelectors'

test('user menu navigates to settings', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/')
  await page.locator(APP_USER_MENU_SUMMARY).click()
  await page.getByRole('link', { name: 'Settings' }).click()
  await expect(page).toHaveURL(/\/settings/)
  await expect(page.getByRole('heading', { name: 'Personal' })).toBeVisible()
})

test('settings page sets user-menu Settings link as current', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await expect(page.locator(USER_MENU_SETTINGS_LINK)).toHaveAttribute(
    'aria-current',
    'page',
  )
})
