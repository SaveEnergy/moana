import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('categories page renders', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await expect(page.getByRole('heading', { name: 'Categories', exact: true })).toBeVisible()
  await expect(page.getByRole('button', { name: 'Add category' })).toBeVisible()
})

test('new transaction page renders', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await expect(page.getByRole('heading', { name: 'New entry' })).toBeVisible()
})

test('settings profile loads', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await expect(page).toHaveURL(/\/settings/)
  await expect(page.getByRole('heading', { name: 'Personal' })).toBeVisible()
})
