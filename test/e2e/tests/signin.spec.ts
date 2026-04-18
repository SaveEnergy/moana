import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('sign in and see overview', async ({ page }) => {
  await signInAsTestUser(page)
  await expect(page.getByRole('heading', { name: 'Dashboard' })).toBeVisible()
})
