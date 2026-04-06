import React from 'react'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi } from 'vitest'
import PrivateRoute from './PrivateRoute'

const useAuthMock = vi.fn()

vi.mock('./AuthContext', () => ({
  useAuth: () => useAuthMock(),
}))

describe('PrivateRoute', () => {
  it('redirects anonymous users to login', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: false })

    render(
      <MemoryRouter initialEntries={['/protected']}>
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route path="/protected" element={<PrivateRoute><div>Protected Page</div></PrivateRoute>} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText('Login Page')).toBeInTheDocument()
  })

  it('renders children for authenticated users', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: true })

    render(
      <MemoryRouter initialEntries={['/protected']}>
        <Routes>
          <Route path="/protected" element={<PrivateRoute><div>Protected Page</div></PrivateRoute>} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText('Protected Page')).toBeInTheDocument()
  })
})
