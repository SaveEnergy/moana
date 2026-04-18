import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

function todayLocalISODate(): string {
  const d = new Date()
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

test('transaction create hydrates local time on history', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/transactions')
  await page.locator('input[name="amount"]').fill('7.50')
  await page.locator('input[name="occurred_on"]').fill(todayLocalISODate())
  await page.getByRole('button', { name: 'Save entry' }).click()
  await expect(page).toHaveURL(/\/history/)

  const stamp = page.locator('time.js-local-time').first()
  await expect(stamp).toBeVisible()
  await expect(stamp).not.toHaveText('…')
})
