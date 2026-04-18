import { test, expect, type Page } from '@playwright/test'

async function signIn(page: Page) {
  await page.goto('/login')
  await page.locator('input[name="email"]').fill('e2e@moana.test')
  await page.locator('input[name="password"]').fill('password123')
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(/\/$/)
}

test('history page loads shell and heading', async ({ page }) => {
  await signIn(page)
  await page.goto('/history')
  await expect(page).toHaveURL(/\/history/)
  await expect(page.getByRole('heading', { name: 'History' })).toBeVisible()
  await expect(page.getByRole('tablist', { name: 'Filter by type' })).toBeVisible()
})
