import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('health endpoint responds OK', async ({ request }) => {
  const res = await request.get('/health')
  expect(res.ok()).toBeTruthy()
})

test('static app.css is served with bytes', async ({ request }) => {
  const res = await request.get('/static/css/app.css')
  expect(res.ok()).toBeTruthy()
  const len = (await res.body()).length
  expect(len).toBeGreaterThan(500)
})

test('static app.js is served with bytes', async ({ request }) => {
  const res = await request.get('/static/js/app.js')
  expect(res.ok()).toBeTruthy()
  const len = (await res.body()).length
  expect(len).toBeGreaterThan(100)
})

test('login page head modulepreloads app.js', async ({ page }) => {
  await page.goto('/login')
  await expect(page.locator('link[rel="modulepreload"][href="/static/js/app.js"]')).toHaveCount(1)
})

const authenticatedShellModulePreload = [
  ['dashboard', '/'],
  ['categories', '/categories'],
  ['history', '/history'],
  ['new transaction', '/transactions'],
  ['notifications', '/notifications'],
  ['settings', '/settings'],
] as const

for (const [label, path] of authenticatedShellModulePreload) {
  test(`layout modulepreloads app.js — ${label}`, async ({ page }) => {
    await signInAsTestUser(page)
    await page.goto(path)
    await expect(page.locator('link[rel="modulepreload"][href="/static/js/app.js"]')).toHaveCount(1)
  })
}
