import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test.beforeEach(async ({ page }) => {
  await signInAsTestUser(page)
})

test('dashboard loads design tokens and overview', async ({ page }) => {
  await page.goto('/')
  const primary = await page.evaluate(() =>
    getComputedStyle(document.documentElement).getPropertyValue('--primary').trim(),
  )
  expect(primary.toLowerCase()).toBe('#306369')
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
  await expect(
    page.getByRole('group', { name: /Statistics and outflow period/i }),
  ).toBeVisible()
  await expect(page.getByText('Total Income', { exact: true })).toBeVisible()
  await expect(page.getByRole('heading', { name: 'Money movement' })).toBeVisible()
})

test('dashboard period links set period query', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('link', { name: '12 months' }).click()
  await expect(page).toHaveURL(/[?&]period=12m(?:&|$)/)
  await page.getByRole('link', { name: '30 days' }).click()
  await expect(page).toHaveURL(/[?&]period=30d(?:&|$)/)
})

test('sidebar main nav links reach primary routes', async ({ page }) => {
  await page.goto('/')
  const nav = page.locator('#app-sidebar-nav')
  await nav.getByRole('link', { name: 'Transactions' }).click()
  await expect(page).toHaveURL(/\/transactions/)
  await nav.getByRole('link', { name: 'History' }).click()
  await expect(page).toHaveURL(/\/history/)
  await nav.getByRole('link', { name: 'Categories' }).click()
  await expect(page).toHaveURL(/\/categories/)
  await nav.getByRole('link', { name: 'Dashboard' }).click()
  await expect(page).toHaveURL(/\/$/)
})

test('sidebar FAB links to new transaction', async ({ page }) => {
  await page.goto('/')
  await page.getByRole('link', { name: 'Add transaction' }).click()
  await expect(page).toHaveURL(/\/transactions/)
  await expect(page.getByRole('heading', { name: 'New entry' })).toBeVisible()
})

test('site footer legal nav exposes MIT and repo links', async ({ page }) => {
  await page.goto('/')
  const nav = page.getByRole('navigation', { name: 'Legal and source' })
  const mit = nav.getByRole('link', { name: 'MIT License' })
  await expect(mit).toHaveAttribute('href', 'https://opensource.org/licenses/MIT')
  await expect(mit).toHaveAttribute('target', '_blank')
  const gh = nav.getByRole('link', { name: 'GitHub' })
  await expect(gh).toHaveAttribute('href', /github\.com\//)
  await expect(gh).toHaveAttribute('target', '_blank')
})

test('global search control is in shell', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByRole('search')).toBeVisible()
  await expect(page.getByPlaceholder('Search')).toBeVisible()
})

test('topbar search submits to history with q', async ({ page }) => {
  await page.goto('/')
  const q = 'e2e-topbar-search'
  await page.locator('#app-global-search').fill(q)
  await page.locator('#app-global-search').press('Enter')
  await expect(page).toHaveURL(/\/history(\?|$)/)
  expect(new URL(page.url()).searchParams.get('q')).toBe(q)
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

test('mobile sidebar closes on backdrop click', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator('#app-shell')
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.locator('#app-sidebar-backdrop').click()
  await expect(shell).not.toHaveClass(/sidebar-open/)
})

test('mobile sidebar closes on drawer close control', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator('#app-shell')
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.locator('#app-sidebar-close').click()
  await expect(shell).not.toHaveClass(/sidebar-open/)
})

test('mobile: Escape does not collapse drawer while account menu is open', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator('#app-shell')
  const menu = page.locator('details.app-user-menu')
  await page.locator('details.app-user-menu summary.app-user-menu-btn').click()
  await expect(menu).toHaveAttribute('open', '')
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.keyboard.press('Escape')
  /* Drawer must stay up while account details[open] is observed (native menu may not dismiss on Esc here). */
  await expect(shell).toHaveClass(/sidebar-open/)
})
