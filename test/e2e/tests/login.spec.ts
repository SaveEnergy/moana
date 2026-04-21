import { test, expect } from '@playwright/test'

test('login page exposes email, password, and sign-in', async ({ page }) => {
  await page.goto('/login')
  await expect(page.getByRole('heading', { name: /sign in to your account/i })).toBeVisible()
  await expect(page.getByLabel('Email address')).toBeVisible()
  await expect(page.getByLabel('Password')).toBeVisible()
  await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible()
})

test('login page exposes remember me and disabled oauth placeholders', async ({ page }) => {
  await page.goto('/login')
  await expect(page.getByRole('checkbox', { name: /remember me/i })).toBeVisible()
  await expect(
    page.getByRole('group', { name: /OAuth sign-in \(not available in this build\)/i }),
  ).toBeVisible()
  await expect(page.getByRole('button', { name: 'Google' })).toBeDisabled()
  await expect(page.getByRole('button', { name: 'GitHub' })).toBeDisabled()
})

test('login page forgot password control is disabled until wired', async ({ page }) => {
  await page.goto('/login')
  const forgot = page.getByRole('button', { name: 'Forgot password?' })
  await expect(forgot).toBeVisible()
  await expect(forgot).toBeDisabled()
})
