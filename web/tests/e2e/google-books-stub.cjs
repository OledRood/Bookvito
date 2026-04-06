const http = require('node:http')

const portArg = process.argv.find((arg) => arg.startsWith('--port='))
const port = Number((portArg && portArg.split('=')[1]) || process.env.GOOGLE_STUB_PORT || '18081')

const json = (res, status, payload) => {
  res.statusCode = status
  res.setHeader('Content-Type', 'application/json')
  res.end(JSON.stringify(payload))
}

const successPayload = {
  items: [
    {
      volumeInfo: {
        title: 'Stub Title',
        authors: ['Stub Author'],
        description: 'Stub description',
        imageLinks: {
          thumbnail: 'http://example.com/stub.png',
        },
      },
    },
  ],
}

const server = http.createServer((req, res) => {
  const requestURL = new URL(req.url || '/', `http://127.0.0.1:${port}`)

  if (requestURL.pathname !== '/volumes') {
    json(res, 404, { error: 'not found' })
    return
  }

  const query = (requestURL.searchParams.get('q') || '').toLowerCase()

  if (query.includes('stub-timeout')) {
    setTimeout(() => json(res, 200, successPayload), 5500)
    return
  }

  if (query.includes('stub-unavailable')) {
    req.socket.destroy()
    return
  }

  if (query.includes('stub-400')) {
    json(res, 400, { error: 'provider bad request' })
    return
  }

  if (query.includes('stub-500')) {
    json(res, 500, { error: 'provider internal error' })
    return
  }

  if (query.includes('stub-empty')) {
    json(res, 200, { items: [] })
    return
  }

  json(res, 200, successPayload)
})

server.listen(port, '127.0.0.1')

const shutdown = () => {
  server.close(() => process.exit(0))
}

process.on('SIGINT', shutdown)
process.on('SIGTERM', shutdown)
