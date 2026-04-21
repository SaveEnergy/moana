import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import {
  APP_GLOBAL_SEARCH,
  APP_SHELL,
  APP_SIDEBAR_BACKDROP,
  APP_SIDEBAR_CLOSE,
  APP_SIDEBAR_NAV,
  NOTIFICATIONS_PATH,
  TOPBAR_NOTIFICATIONS_LINK,
} from '../helpers/shellSelectors'

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

test('dashboard period active class follows period query', async ({ page }) => {
  await page.goto('/?period=12m')
  await expect(page.getByRole('link', { name: '12 months' })).toHaveClass(/is-active/)
  await expect(page.getByRole('link', { name: '30 days' })).not.toHaveClass(/is-active/)
})

test('sidebar main nav links reach primary routes', async ({ page }) => {
  await page.goto('/')
  const nav = page.locator(APP_SIDEBAR_NAV)
  await nav.getByRole('link', { name: 'Transactions' }).click()
  await expect(page).toHaveURL(/\/transactions/)
  await nav.getByRole('link', { name: 'History' }).click()
  await expect(page).toHaveURL(/\/history/)
  await nav.getByRole('link', { name: 'Categories' }).click()
  await expect(page).toHaveURL(/\/categories/)
  await nav.getByRole('link', { name: 'Dashboard' }).click()
  await expect(page).toHaveURL(/\/$/)
})

test('sidebar active link matches current route', async ({ page }) => {
  await page.goto('/history')
  await expect(page.locator(`${APP_SIDEBAR_NAV} a[href="/history"]`)).toHaveClass(/sidebar-link-active/)
  await expect(page.locator(`${APP_SIDEBAR_NAV} a[href="/"]`)).not.toHaveClass(/sidebar-link-active/)

  await page.goto('/transactions')
  await expect(page.locator(`${APP_SIDEBAR_NAV} a[href="/transactions"]`)).toHaveClass(/sidebar-link-active/)

  await page.goto('/categories')
  await expect(page.locator(`${APP_SIDEBAR_NAV} a[href="/categories"]`)).toHaveClass(/sidebar-link-active/)

  await page.goto('/')
  await expect(page.locator(`${APP_SIDEBAR_NAV} a[href="/"]`)).toHaveClass(/sidebar-link-active/)
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
  await page.locator(APP_GLOBAL_SEARCH).fill(q)
  await page.locator(APP_GLOBAL_SEARCH).press('Enter')
  await expect(page).toHaveURL(/\/history(\?|$)/)
  expect(new URL(page.url()).searchParams.get('q')).toBe(q)
})

test('notifications link is reachable', async ({ page }) => {
  await page.goto('/')
  await expect(page.locator(TOPBAR_NOTIFICATIONS_LINK)).toBeVisible()
  await page.locator(TOPBAR_NOTIFICATIONS_LINK).click()
  await expect(page).toHaveURL(new RegExp(`${NOTIFICATIONS_PATH}$`))
  await expect(page.getByRole('heading', { name: 'Notifications' })).toBeVisible()
})

test('mobile sidebar toggles', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator(APP_SHELL)
  await expect(shell).not.toHaveClass(/sidebar-open/)
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.keyboard.press('Escape')
  await expect(shell).not.toHaveClass(/sidebar-open/)
})

test('mobile sidebar closes on backdrop click', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator(APP_SHELL)
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.locator(APP_SIDEBAR_BACKDROP).click()
  await expect(shell).not.toHaveClass(/sidebar-open/)
})

test('mobile sidebar closes on drawer close control', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator(APP_SHELL)
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.locator(APP_SIDEBAR_CLOSE).click()
  await expect(shell).not.toHaveClass(/sidebar-open/)
})

test('mobile: Escape does not collapse drawer while account menu is open', async ({ page }) => {
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/')
  const shell = page.locator(APP_SHELL)
  const menu = page.locator('details.app-user-menu')
  await page.locator('details.app-user-menu summary.app-user-menu-btn').click()
  await expect(menu).toHaveAttribute('open', '')
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  await page.keyboard.press('Escape')
  /* Drawer must stay up while account details[open] is observed (native menu may not dismiss on Esc here). */
  await expect(shell).toHaveClass(/sidebar-open/)
})
