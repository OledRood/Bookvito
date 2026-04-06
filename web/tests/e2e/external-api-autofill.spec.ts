import path from 'node:path'
import { expect, test } from '@playwright/test'
import { backendBaseURL } from './helpers'

const authDir = path.resolve(__dirname, '.auth')

test.describe('external api autofill backend flow', () => {
  test('returns metadata on provider success', async ({ request }) => {
    const response = await request.get(`${backendBaseURL()}/books/auto-fill?q=stub-success`)
    expect(response.status()).toBe(200)

    const payload = await response.json() as {
      title?: string
      author?: string
      description?: string
    }
    expect(payload.title).toBe('Stub Title')
    expect(payload.author).toBe('Stub Author')
    expect(payload.description).toBe('Stub description')
  })

  test('returns not_found for empty provider response', async ({ request }) => {
    const response = await request.get(`${backendBaseURL()}/books/auto-fill?q=stub-empty`)
    expect(response.status()).toBe(404)

    const payload = await response.json() as { code?: string; error?: string }
    expect(payload.code).toBe('not_found')
    expect(payload.error).toContain('book not found')
  })

  test('returns internal for provider 4xx/5xx/timeout/unavailable', async ({ request }) => {
    for (const query of ['stub-400', 'stub-500', 'stub-timeout', 'stub-unavailable']) {
      const response = await request.get(`${backendBaseURL()}/books/auto-fill?q=${query}`)
      expect(response.status()).toBe(500)

      const payload = await response.json() as { code?: string; error?: string }
      expect(payload.code).toBe('internal')
      expect(payload.error).toBeTruthy()
    }
  })
})

test.describe('external api autofill ui flow', () => {
  test.use({ storageState: path.join(authDir, 'user.json') })

  test('shows error and keeps form interactive when provider fails', async ({ page }) => {
    await page.goto('/books/new')

    const titleInput = page.getByLabel('Название книги')
    const authorInput = page.getByLabel('Автор(ы)')
    const autofillButton = page.getByRole('button', { name: 'Автозаполнение' })

    await titleInput.fill('stub-success')
    await autofillButton.click()
    await expect(authorInput).toHaveValue('Stub Author')

    await titleInput.fill('stub-empty')
    await autofillButton.click()
    await expect(page.getByText('book not found')).toBeVisible()
    await expect(authorInput).toHaveValue('Stub Author')

    await titleInput.fill('stub-500')
    await autofillButton.click()
    await expect(page.getByText('google books responded with status 500')).toBeVisible()

    await expect(page).toHaveURL('/books/new')
    await expect(titleInput).toBeVisible()
    await expect(autofillButton).toBeEnabled()
  })
})
