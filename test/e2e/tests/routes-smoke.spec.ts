import { test, expect } from '@playwright/test'

import { E2E_USER_EMAIL, signInAsTestUser } from '../helpers/auth'
import { todayLocalISODate } from '../helpers/dates'
import { APP_CSS_STYLESHEET, APP_JS_MODULE_PRELOAD } from '../helpers/modulePreload'
import {
  SETTINGS_EMAIL,
  TX_EDIT_NOTE,
  TX_INPUT_AMOUNT,
  TX_INPUT_KIND_EXPENSE,
  TX_INPUT_KIND_INCOME,
  TX_INPUT_OCCURRED_ON,
  TX_KIND_LABEL_EXPENSE,
  TX_KIND_LABEL_INCOME,
} from '../helpers/shellSelectors'

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
  await page.locator(TX_KIND_LABEL_EXPENSE).click()
  await expect(page.locator(TX_INPUT_KIND_EXPENSE)).toBeChecked()
  await page.locator(TX_KIND_LABEL_INCOME).click()
  await expect(page.locator(TX_INPUT_KIND_INCOME)).toBeChecked()
})

test('new transaction shows alert for zero amount', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await page.locator(TX_INPUT_AMOUNT).fill('0.00')
  await page.locator(TX_INPUT_OCCURRED_ON).fill(todayLocalISODate())
  await page.getByRole('button', { name: 'Save entry' }).click()
  await expect(page).toHaveURL(/\/transactions/)
  await expect(page.getByRole('alert')).toContainText('Amount must be greater than zero.')
})

test('settings profile loads', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/settings')
  await expect(page).toHaveURL(/\/settings/)
  await expect(page.getByRole('heading', { name: 'Personal' })).toBeVisible()
  await expect(page.locator(SETTINGS_EMAIL)).toHaveValue(E2E_USER_EMAIL)
})

test('transaction edit loads from history and save returns to history', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await page.locator(TX_INPUT_AMOUNT).fill('1.00')
  await page.locator(TX_INPUT_OCCURRED_ON).fill(todayLocalISODate())
  await page.getByRole('button', { name: 'Save entry' }).click()
  await expect(page).toHaveURL(/\/history/)

  await page.getByRole('link', { name: 'Edit' }).first().click()
  await expect(page).toHaveURL(/\/transactions\/\d+\/edit/)
  await expect(page.getByRole('heading', { name: 'Edit entry' })).toBeVisible()
  await expect(page.locator(APP_CSS_STYLESHEET)).toHaveCount(1)
  await expect(page.locator(APP_JS_MODULE_PRELOAD)).toHaveCount(1)

  await page.locator(TX_EDIT_NOTE).fill('e2e edit smoke')
  await page.getByRole('button', { name: 'Save changes' }).click()
  await expect(page).toHaveURL(/\/history/)
})
