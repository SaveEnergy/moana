import { test, expect } from '@playwright/test'

test('sign in and see overview', async ({ page }) => {
  await page.goto('/login')
  await page.locator('input[name="email"]').fill('e2e@moana.test')
  await page.locator('input[name="password"]').fill('password123')
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(/\/$/)
  await expect(page.getByText('Running total')).toBeVisible()
})
