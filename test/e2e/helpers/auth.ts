import { expect, type Page } from '@playwright/test'

/** Must match `test/e2e/scripts/start-server.sh` seeded user (`go run … user add`). */
export const E2E_USER_EMAIL = 'e2e@moana.test'
export const E2E_USER_PASSWORD = 'password123'

/** Complete login flow; lands on dashboard (`/`). */
export async function signInAsTestUser(page: Page): Promise<void> {
  await page.goto('/login')
  await page.locator('input[name="email"]').fill(E2E_USER_EMAIL)
  await page.locator('input[name="password"]').fill(E2E_USER_PASSWORD)
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(/\/$/)
}
