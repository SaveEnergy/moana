import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'

test('category modal opens and closes', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator('#cat-modal-open-create').click()
  const dialog = page.locator('#cat-modal')
  await expect(dialog).toBeVisible()
  await expect(page.locator('#cat-modal-title')).toHaveText('New category')
  await page.locator('#cat-modal-close').click()
  await expect(dialog).toBeHidden()
})

test('category modal closes on Escape', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator('#cat-modal-open-create').click()
  const dialog = page.locator('#cat-modal')
  await expect(dialog).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dialog).toBeHidden()
})

test('category modal closes on backdrop click outside panel', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator('#cat-modal-open-create').click()
  const dialog = page.locator('#cat-modal')
  const panel = page.locator('.cat-modal-panel')
  await expect(dialog).toBeVisible()
  const box = await panel.boundingBox()
  expect(box).not.toBeNull()
  /* Viewport pixel above the card — ::backdrop / dimmed area (regresses if dismiss logic misses engines). */
  await page.mouse.click(box!.x + box!.width / 2, Math.max(8, box!.y - 24))
  await expect(dialog).toBeHidden()
})

test('mobile: first Escape closes modal with sidebar left open', async ({ page }) => {
  await signInAsTestUser(page)
  await page.setViewportSize({ width: 600, height: 800 })
  await page.goto('/categories')
  const shell = page.locator('#app-shell')
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  /* Drawer backdrop blocks “Add category”; opening the dialog via showModal() matches the Escape / shell interaction under test. */
  await page.evaluate(() => {
    document.getElementById('cat-modal')?.showModal()
  })
  await expect(page.locator('#cat-modal')).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(page.locator('#cat-modal')).toBeHidden()
  await expect(shell).toHaveClass(/sidebar-open/)
})
