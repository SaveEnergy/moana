import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import { todayLocalISODate } from '../helpers/dates'

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

test('new transaction kind pills switch income and expense', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  // Inputs are visually hidden (styled pill labels); click the labels like a user.
  await page.locator('label.kind-option:has(input[name="kind"][value="expense"])').click()
  await expect(page.locator('input[name="kind"][value="expense"]')).toBeChecked()
  await page.locator('label.kind-option:has(input[name="kind"][value="income"])').click()
  await expect(page.locator('input[name="kind"][value="income"]')).toBeChecked()
})

test('settings profile loads', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await expect(page).toHaveURL(/\/settings/)
  await expect(page.getByRole('heading', { name: 'Personal' })).toBeVisible()
})

test('transaction edit loads from history and save returns to history', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await page.locator('input[name="amount"]').fill('1.00')
  await page.locator('input[name="occurred_on"]').fill(todayLocalISODate())
  await page.getByRole('button', { name: 'Save entry' }).click()
  await expect(page).toHaveURL(/\/history/)

  await page.getByRole('link', { name: 'Edit' }).first().click()
  await expect(page).toHaveURL(/\/transactions\/\d+\/edit/)
  await expect(page.getByRole('heading', { name: 'Edit entry' })).toBeVisible()

  await page.locator('#tx-edit-note').fill('e2e edit smoke')
  await page.getByRole('button', { name: 'Save changes' }).click()
  await expect(page).toHaveURL(/\/history/)
})
