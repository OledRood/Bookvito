import { expect, test } from '@playwright/test'
import { authHeader, backendBaseURL } from './helpers'

test('upload, download and detach image from object storage', async ({ request }) => {
  const api = backendBaseURL()
  const unique = Date.now().toString()

  const uploadResponse = await request.post(`${api}/books/upload`, {
    headers: authHeader('user'),
    multipart: {
      image: {
        name: 'storage-cover.png',
        mimeType: 'image/png',
        buffer: Buffer.from(`storage-${unique}`),
      },
    },
  })
  expect(uploadResponse.status()).toBe(200)
  const uploadPayload = await uploadResponse.json() as { url: string }
  expect(uploadPayload.url).toBeTruthy()

  const downloadResponse = await request.get(uploadPayload.url)
  expect(downloadResponse.status()).toBe(200)
  const downloaded = await downloadResponse.body()
  expect(downloaded.toString()).toContain(`storage-${unique}`)

  const createBookResponse = await request.post(`${api}/books/create`, {
    headers: { ...authHeader('user'), 'Content-Type': 'application/json' },
    data: {
      title: `Storage Book ${unique}`,
      author: 'Storage Bot',
      description: 'Storage flow',
      condition: 'good',
      image_url: uploadPayload.url,
      current_location_id: null,
    },
  })
  expect(createBookResponse.status()).toBe(201)

  const myBooksResponse = await request.get(`${api}/books/my`, { headers: authHeader('user') })
  expect(myBooksResponse.status()).toBe(200)
  const books = await myBooksResponse.json() as Array<{ id: string; title: string }>
  const book = books.find((entry) => entry.title === `Storage Book ${unique}`)
  expect(book).toBeTruthy()

  const detachResponse = await request.delete(`${api}/books/image/${book!.id}`, {
    headers: authHeader('user'),
  })
  expect(detachResponse.status()).toBe(200)
})
