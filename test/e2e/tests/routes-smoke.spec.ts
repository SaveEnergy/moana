import { test, expect, type Page } from '@playwright/test'

async function signIn(page: Page) {
  await page.goto('/login')
  await page.locator('input[name="email"]').fill('e2e@moana.test')
  await page.locator('input[name="password"]').fill('password123')
  await page.getByRole('button', { name: /sign in/i }).click()
  await expect(page).toHaveURL(/\/$/)
}

test('categories page renders', async ({ page }) => {
  await signIn(page)
  await page.goto('/categories')
  await expect(page.getByRole('heading', { name: 'Categories', exact: true })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Add category' })).toBeVisible()
})

test('new transaction page renders', async ({ page }) => {
  await signIn(page)
  await page.goto('/transactions')
  await expect(page.getByRole('heading', { name: 'New entry' })).toBeVisible()
})

test('settings profile loads', async ({ page }) => {
  await signIn(page)
  await page.goto('/settings')
  await expect(page).toHaveURL(/\/settings/)
  await expect(page.getByRole('heading', { name: 'Personal' })).toBeVisible()
})
