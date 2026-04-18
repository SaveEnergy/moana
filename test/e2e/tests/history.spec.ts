import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('history page loads shell and heading', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history')
  await expect(page).toHaveURL(/\/history/)
  await expect(page.getByRole('heading', { name: 'History' })).toBeVisible()
  await expect(page.getByRole('tablist', { name: 'Filter by type' })).toBeVisible()
})

test('history search submits q on GET', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history')
  const q = 'e2e-history-q'
  await page.locator('#history-q').fill(q)
  await page.locator('#history-q').press('Enter')
  await expect(page).toHaveURL(new RegExp(`[?&]q=${encodeURIComponent(q)}`))
})

test('history kind tabs preserve query and switch kind', async ({ page }) => {
  await signInAsTestUser(page)
  const q = 'e2e-kind-nav'
  await page.goto(`/history?q=${encodeURIComponent(q)}`)
  await page.getByRole('tab', { name: 'Expenses' }).click()
  await expect(page).toHaveURL(new RegExp(`[?&]q=${encodeURIComponent(q)}`))
  await expect(page).toHaveURL(/[?&]kind=expense(?:&|$)/)
  await page.getByRole('tab', { name: 'Income' }).click()
  await expect(page).toHaveURL(new RegExp(`[?&]q=${encodeURIComponent(q)}`))
  await expect(page).toHaveURL(/[?&]kind=income(?:&|$)/)
  await page.getByRole('tab', { name: 'All' }).click()
  await expect(page).toHaveURL(new RegExp(`[?&]q=${encodeURIComponent(q)}`))
  await expect(page).toHaveURL(/[?&]kind=all(?:&|$)/)
})
