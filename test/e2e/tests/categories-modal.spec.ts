import { test, expect } from '@playwright/test'

import { signInAsTestUser } from '../helpers/auth'
import {
  APP_SHELL,
  CATEGORY_MODAL,
  CATEGORY_MODAL_CLOSE,
  CATEGORY_MODAL_ID,
  CATEGORY_MODAL_NAME,
  CATEGORY_MODAL_OPEN_CREATE,
  CATEGORY_MODAL_SUBMIT,
  CATEGORY_MODAL_TITLE,
} from '../helpers/shellSelectors'

function uniqueCategoryName(prefix: string): string {
  return `${prefix}-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`
}

test('category modal opens and closes', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator(CATEGORY_MODAL_OPEN_CREATE).click()
  const dialog = page.locator(CATEGORY_MODAL)
  await expect(dialog).toBeVisible()
  await expect(page.locator(CATEGORY_MODAL_TITLE)).toHaveText('New category')
  await page.locator(CATEGORY_MODAL_CLOSE).click()
  await expect(dialog).toBeHidden()
})

test('category modal closes on Escape', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator(CATEGORY_MODAL_OPEN_CREATE).click()
  const dialog = page.locator(CATEGORY_MODAL)
  await expect(dialog).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(dialog).toBeHidden()
})

test('category modal closes on backdrop click outside panel', async ({ page }) => {
  await signInAsTestUser(page)
  await page.goto('/categories')
  await page.locator(CATEGORY_MODAL_OPEN_CREATE).click()
  const dialog = page.locator(CATEGORY_MODAL)
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
  const shell = page.locator(APP_SHELL)
  await page.getByRole('button', { name: 'Open navigation menu' }).click()
  await expect(shell).toHaveClass(/sidebar-open/)
  /* Drawer backdrop blocks “Add category”; opening the dialog via showModal() matches the Escape / shell interaction under test. */
  await page.evaluate((id) => {
    document.getElementById(id)?.showModal()
  }, CATEGORY_MODAL_ID)
  await expect(page.locator(CATEGORY_MODAL)).toBeVisible()
  await page.keyboard.press('Escape')
  await expect(page.locator(CATEGORY_MODAL)).toBeHidden()
  await expect(shell).toHaveClass(/sidebar-open/)
})

test('category delete confirm dismiss keeps category', async ({ page }) => {
  await signInAsTestUser(page)
  const name = uniqueCategoryName('e2e-keep')
  await page.goto('/categories')
  await page.locator(CATEGORY_MODAL_OPEN_CREATE).click()
  await page.locator(CATEGORY_MODAL_NAME).fill(name)
  await page.locator(CATEGORY_MODAL_SUBMIT).click()
  await expect(page.getByText(name)).toBeVisible()

  page.once('dialog', async (dialog) => {
    expect(dialog.type()).toBe('confirm')
    expect(dialog.message()).toMatch(/remove this category/i)
    await dialog.dismiss()
  })
  await page
    .locator('.cat-list-row')
    .filter({ hasText: name })
    .locator('form.cat-delete')
    .getByRole('button', { name: 'Remove' })
    .click()
  await expect(page.getByText(name)).toBeVisible()
})

test('category delete confirm accept removes category', async ({ page }) => {
  await signInAsTestUser(page)
  const name = uniqueCategoryName('e2e-gone')
  await page.goto('/categories')
  await page.locator(CATEGORY_MODAL_OPEN_CREATE).click()
  await page.locator(CATEGORY_MODAL_NAME).fill(name)
  await page.locator(CATEGORY_MODAL_SUBMIT).click()
  await expect(page.getByText(name)).toBeVisible()

  page.once('dialog', async (dialog) => {
    expect(dialog.type()).toBe('confirm')
    await dialog.accept()
  })
  await page
    .locator('.cat-list-row')
    .filter({ hasText: name })
    .locator('form.cat-delete')
    .getByRole('button', { name: 'Remove' })
    .click()
  await expect(page.getByText(name)).not.toBeVisible()
})
