import { expect, test } from '@playwright/test'
import path from 'node:path'

const authDir = path.resolve(__dirname, '.auth')

test.describe('role-based routes', () => {
  test.use({ storageState: path.join(authDir, 'user.json') })

  test('user sees 403 on admin route', async ({ page }) => {
    await page.goto('/admin')
    await expect(page.getByText('Недостаточно прав для просмотра этой страницы')).toBeVisible()
  })
})

test.describe('moder route', () => {
  test.use({ storageState: path.join(authDir, 'moder.json') })

  test('moder can open moderation page', async ({ page }) => {
    await page.goto('/moderation')
    await expect(page.getByText('Жалобы на объявления')).toBeVisible()
  })
})

test.describe('admin route', () => {
  test.use({ storageState: path.join(authDir, 'admin.json') })

  test('admin can open admin page', async ({ page }) => {
    await page.goto('/admin')
    await expect(page.getByText('Панель администратора')).toBeVisible()
  })
})
