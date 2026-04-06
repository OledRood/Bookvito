import { http, HttpResponse } from 'msw'

export const handlers = [
  http.get('http://localhost:8080/api/v1/', () => HttpResponse.json({ status: 'ok' })),
  http.get('http://127.0.0.1:18080/api/v1/', () => HttpResponse.json({ status: 'ok' })),
  http.get('http://localhost:8080/api/v1/locations/getAll', () => HttpResponse.json([])),
  http.get('http://127.0.0.1:18080/api/v1/locations/getAll', () => HttpResponse.json([])),
]
