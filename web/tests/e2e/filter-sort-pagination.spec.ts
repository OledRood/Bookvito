import { expect, test } from '@playwright/test'
import { authHeader, backendBaseURL } from './helpers'

test('filtering, sorting and pagination flow on catalog endpoint', async ({ request }) => {
  const api = backendBaseURL()

  const firstPageResponse = await request.get(`${api}/books/list?limit=1&offset=0&sort_by=title&order=asc`, {
    headers: authHeader('user'),
  })
  expect(firstPageResponse.status()).toBe(200)
  const firstPagePayload = await firstPageResponse.json() as { items: Array<{ id: string }>; has_more: boolean }
  expect(firstPagePayload.items.length).toBe(1)
  expect(typeof firstPagePayload.has_more).toBe('boolean')

  const secondPageResponse = await request.get(`${api}/books/list?limit=1&offset=1&sort_by=title&order=asc`, {
    headers: authHeader('user'),
  })
  expect(secondPageResponse.status()).toBe(200)
  const secondPagePayload = await secondPageResponse.json() as { items: Array<{ id: string }>; has_more: boolean }
  if (secondPagePayload.items.length > 0) {
    expect(secondPagePayload.items[0].id).not.toBe(firstPagePayload.items[0].id)
  }

  const invalidSortResponse = await request.get(`${api}/books/list?limit=1&sort_by=unknown_field&order=asc`, {
    headers: authHeader('user'),
  })
  expect(invalidSortResponse.status()).toBe(400)
  const invalidPayload = await invalidSortResponse.json() as { code?: string; error?: string }
  expect(invalidPayload.code).toBe('validation')
})
