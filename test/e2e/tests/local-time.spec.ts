import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import { todayLocalISODate } from '../helpers/dates'
import { TX_INPUT_AMOUNT, TX_INPUT_OCCURRED_ON } from '../helpers/shellSelectors'

test('transaction create shows amount on history', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await page.locator(TX_INPUT_AMOUNT).fill('7.50')
  await page.locator(TX_INPUT_OCCURRED_ON).fill(todayLocalISODate())
  await page.getByRole('button', { name: 'Save entry' }).click()
  await expect(page).toHaveURL(/\/history/)

  await expect(page.getByText('€7.50').first()).toBeVisible()
})
