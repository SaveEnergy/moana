import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import { todayLocalISODate } from '../helpers/dates'
import { LOCAL_TIME_DISPLAY, TX_INPUT_AMOUNT, TX_INPUT_OCCURRED_ON } from '../helpers/shellSelectors'

test('transaction create hydrates local time on history', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await page.locator(TX_INPUT_AMOUNT).fill('7.50')
  await page.locator(TX_INPUT_OCCURRED_ON).fill(todayLocalISODate())
  await page.getByRole('button', { name: 'Save entry' }).click()
  await expect(page).toHaveURL(/\/history/)

  const stamp = page.locator(LOCAL_TIME_DISPLAY).first()
  await expect(stamp).toBeVisible()
  await expect(stamp).not.toHaveText('…')
})
