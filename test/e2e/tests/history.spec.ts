import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import { HISTORY_FROM, HISTORY_Q, HISTORY_SORT, HISTORY_TO } from '../helpers/shellSelectors'

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
  await page.locator(HISTORY_Q).fill(q)
  await page.locator(HISTORY_Q).press('Enter')
  await expect(page).toHaveURL(new RegExp(`[?&]q=${encodeURIComponent(q)}`))
})

test('history sort select submits sort on change', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history')
  await page.locator(HISTORY_SORT).selectOption('oldest')
  await expect(page).toHaveURL(/[?&]sort=oldest(?:&|$)/)
})

test('history clear filters link returns to bare /history', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history?q=e2e-clear&kind=expense&sort=oldest')
  await page.getByRole('link', { name: 'Clear filters' }).click()
  await expect(page).toHaveURL(/\/history$/)
  expect(new URL(page.url()).search).toBe('')
})

test('history apply dates submits from and to on GET', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history')
  await page.locator(HISTORY_FROM).fill('2026-01-10')
  await page.locator(HISTORY_TO).fill('2026-04-19')
  await page.getByRole('button', { name: 'Apply dates' }).click()
  await expect(page).toHaveURL(/[?&]from=2026-01-10(?:&|$)/)
  await expect(page).toHaveURL(/[?&]to=2026-04-19(?:&|$)/)
})

test('history kind tab active class follows kind query', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/history?kind=expense')
  await expect(page.getByRole('tab', { name: 'Expenses' })).toHaveClass(/history-seg-active/)
  await expect(page.getByRole('tab', { name: 'All' })).not.toHaveClass(/history-seg-active/)
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
