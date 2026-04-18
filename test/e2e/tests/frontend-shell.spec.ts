import { test, expect, type Page } from '@playwright/test'

async function signIn(page: Page) {
  await page.goto('/login')
  await page.locator('input[name="email"]').fill('e2e@moana.test')
  await page.locator('input[name="password"]').fill('password123')
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(/\/$/)
}

test.beforeEach(async ({ page }) => {
  await signIn(page)
})

test('dashboard loads design tokens and overview', async ({ page }) => {
  await page.goto('/')
  const primary = await page.evaluate(() =>
    getComputedStyle(document.documentElement).getPropertyValue('--primary').trim(),
  )
  expect(primary.toLowerCase()).toBe('#306369')
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
})

test('notifications link is reachable', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('link', { name: 'Notifications' }).click()
  await expect(page).toHaveURL(/\/notifications/)
  await expect(page.getByRole('heading', { name: 'Notifications' })).toBeVisible()
})

test('mobile sidebar toggles', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator('#app-shell')
  await expect(shell).not.toHaveClass(/sidebar-open/)
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.keyboard.press('Escape')
  await expect(shell).not.toHaveClass(/sidebar-open/)
})
