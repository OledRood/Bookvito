import { expect, test } from '@playwright/test'

test('user can log in, reload, and log out', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('Email адрес').fill('user@bookvito.test')
  await page.getByLabel('Пароль').fill('password123')
  await page.getByRole('button', { name: 'Войти' }).click()

  await expect(page).toHaveURL('/')

  await page.goto('/profile')
  await expect(page.getByText('Выйти')).toBeVisible()

  await page.reload()
  await expect(page.getByText('Выйти')).toBeVisible()

  await page.getByText('Выйти').click()
  await expect(page).toHaveURL(/\/login/)
})

test('session is invalidated when refresh token is expired or invalid', async ({ page }) => {
  await page.goto('/login')
  await page.getByLabel('Email адрес').fill('user@bookvito.test')
  await page.getByLabel('Пароль').fill('password123')
  await page.getByRole('button', { name: 'Войти' }).click()
  await expect(page).toHaveURL('/')

  await page.evaluate(() => {
    localStorage.setItem('accessToken', 'invalid-token')
    localStorage.setItem('token', 'invalid-token')
    localStorage.setItem('refreshToken', 'expired-refresh-token')
  })

  await page.goto('/books/my')
  await expect(page).toHaveURL(/\/login/)
})
