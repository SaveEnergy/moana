import { test, expect, type Page } from '@playwright/test'

async function signIn(page: Page) {
  await page.goto('/login')
  await page.locator('input[name="email"]').fill('e2e@moana.test')
  await page.locator('input[name="password"]').fill('password123')
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(/\/$/)
}

test('notifications inbox empty state', async ({ page }) => {
  await signIn(page)
  await page.goto('/notifications')
  await expect(page).toHaveURL(/\/notifications/)
  await expect(page.getByRole('heading', { name: 'Notifications' })).toBeVisible()
  await expect(page.getByText('You have no notifications.')).toBeVisible()
})
