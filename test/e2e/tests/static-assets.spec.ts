import { test, expect } from '@playwright/test'

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
