import { expect, test } from '@playwright/test'
import { authHeader, backendBaseURL, readRoleState } from './helpers'

test('CRUD operations by roles', async ({ request }) => {
  const userState = readRoleState('user')
  const api = backendBaseURL()
  const unique = Date.now().toString()

  const uploadResponse = await request.post(`${api}/books/upload`, {
    headers: authHeader('user'),
    multipart: {
      image: {
        name: 'crud-cover.png',
        mimeType: 'image/png',
        buffer: Buffer.from(`cover-${unique}`),
      },
    },
  })
  expect(uploadResponse.status()).toBe(200)
  const uploadPayload = await uploadResponse.json()
  expect(uploadPayload.url).toBeTruthy()

  const createResponse = await request.post(`${api}/books/create`, {
    headers: { ...authHeader('user'), 'Content-Type': 'application/json' },
    data: {
      title: `E2E CRUD ${unique}`,
      author: 'E2E Bot',
      description: 'CRUD scenario',
      condition: 'good',
      image_url: uploadPayload.url,
      current_location_id: null,
    },
  })
  expect(createResponse.status()).toBe(201)

  const myBooksResponse = await request.get(`${api}/books/my`, {
    headers: authHeader('user'),
  })
  expect(myBooksResponse.status()).toBe(200)
  const books = (await myBooksResponse.json()) as Array<{ id: string; title: string }>
  const createdBook = books.find((book) => book.title === `E2E CRUD ${unique}`)
  expect(createdBook).toBeTruthy()

  const updateResponse = await request.put(`${api}/books/${createdBook!.id}`, {
    headers: { ...authHeader('user'), 'Content-Type': 'application/json' },
    data: {
      description: 'CRUD scenario updated',
    },
  })
  expect(updateResponse.status()).toBe(200)

  const adminUsersResponse = await request.get(`${api}/admin/users`, {
    headers: authHeader('admin'),
  })
  expect(adminUsersResponse.status()).toBe(200)
  const users = (await adminUsersResponse.json()) as Array<{ id: string; email: string; role: string }>
  const user = users.find((entry) => entry.email === userState.user.email)
  expect(user).toBeTruthy()

  const elevateRoleResponse = await request.put(`${api}/admin/users/${user!.id}/role`, {
    headers: { ...authHeader('admin'), 'Content-Type': 'application/json' },
    data: { role: 'moder' },
  })
  expect(elevateRoleResponse.status()).toBe(200)

  const restoreRoleResponse = await request.put(`${api}/admin/users/${user!.id}/role`, {
    headers: { ...authHeader('admin'), 'Content-Type': 'application/json' },
    data: { role: 'user' },
  })
  expect(restoreRoleResponse.status()).toBe(200)

  const reportsResponse = await request.get(`${api}/moder/reports`, {
    headers: authHeader('moder'),
  })
  expect(reportsResponse.status()).toBe(200)
  const reports = (await reportsResponse.json()) as Array<{ id: string }>
  if (reports.length > 0) {
    const resolveResponse = await request.put(`${api}/moder/reports/${reports[0].id}/resolve`, {
      headers: authHeader('moder'),
    })
    expect(resolveResponse.status()).toBe(200)
  }
})
