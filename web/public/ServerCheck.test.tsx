import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import { http, HttpResponse } from 'msw'
import ServerCheck from './ServerCheck'
import { server } from '../tests/unit/setup'

describe('ServerCheck', () => {
  it('renders children when backend responds', async () => {
    server.use(
      http.get('http://localhost:8080/api/v1/', () => HttpResponse.json({ status: 'ok' })),
    )

    render(
      <ServerCheck>
        <div>Ready</div>
      </ServerCheck>,
    )

    await waitFor(() => expect(screen.getByText('Ready')).toBeInTheDocument())
  })
})
