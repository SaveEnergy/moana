import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import {
  APP_CSS_STYLESHEET,
  APP_JS_MODULE_PRELOAD,
  STATIC_APP_CSS_PATH,
  STATIC_APP_JS_PATH,
} from '../helpers/modulePreload'

test('health endpoint responds OK', async ({ request }) => {
  const res = await request.get('/health')
  expect(res.ok()).toBeTruthy()
})

test('static app.css is served with bytes', async ({ request }) => {
  const res = await request.get(STATIC_APP_CSS_PATH)
  expect(res.ok()).toBeTruthy()
  const len = (await res.body()).length
  expect(len).toBeGreaterThan(500)
})

test('built app.css retains design-system shadow custom properties', async ({ request }) => {
  const res = await request.get(STATIC_APP_CSS_PATH)
  expect(res.ok()).toBeTruthy()
  const text = await res.text()
  expect(text).toContain('--shadow-modal-panel:')
  expect(text).toContain('--shadow-category-tile-selected:')
  expect(text).toContain('--shadow-drawer-edge:')
  expect(text).toContain('--shadow-primary-fab:')
  expect(text).toContain('--ring-float-field-focus:')
  expect(text).toContain('--shadow-admin-add-dialog:')
  expect(text).toContain('--backdrop-cat-modal:')
  expect(text).toContain('--inset-donut-hole:')
  expect(text).toContain('--shadow-input-underline-focus:')
})

test('static app.js is served with bytes', async ({ request }) => {
  const res = await request.get(STATIC_APP_JS_PATH)
  expect(res.ok()).toBeTruthy()
  const len = (await res.body()).length
  expect(len).toBeGreaterThan(100)
})

test('favicon SVG is embedded and served (production image must keep file pre-go:embed)', async ({
  request,
}) => {
  const res = await request.get('/static/favicon-mo.svg')
  expect(res.ok()).toBeTruthy()
  const text = await res.text()
  expect(text).toContain('<svg')
})

test('GET /favicon.ico redirects to static SVG', async ({ request }) => {
  const res = await request.get('/favicon.ico', { maxRedirects: 0 })
  expect(res.status()).toBe(302)
  expect(res.headers().location).toBe('/static/favicon-mo.svg')
})

test('login page head links app.css and modulepreloads app.js', async ({ page }) => {
  await page.goto('/login')
  await expect(page.locator(APP_CSS_STYLESHEET)).toHaveCount(1)
  await expect(page.locator(APP_JS_MODULE_PRELOAD)).toHaveCount(1)
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
  test(`layout head links app.css and modulepreloads app.js — ${label}`, async ({ page }) => {
    await signInAsTestUser(page)
    await page.goto(path)
    await expect(page.locator(APP_CSS_STYLESHEET)).toHaveCount(1)
    await expect(page.locator(APP_JS_MODULE_PRELOAD)).toHaveCount(1)
  })
}
