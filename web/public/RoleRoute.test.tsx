import React from 'react'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi } from 'vitest'
import RoleRoute from './RoleRoute'

const useAuthMock = vi.fn()

vi.mock('./AuthContext', () => ({
  useAuth: () => useAuthMock(),
}))

describe('RoleRoute', () => {
  it('shows 403 stub when role is insufficient', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: true, user: { role: 'user' } })

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <Routes>
          <Route path="/admin" element={<RoleRoute requiredRole="admin"><div>Admin Page</div></RoleRoute>} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText('403')).toBeInTheDocument()
  })

  it('allows admin to access admin route', () => {
    useAuthMock.mockReturnValue({ isAuthenticated: true, user: { role: 'admin' } })

    render(
      <MemoryRouter initialEntries={['/admin']}>
        <Routes>
          <Route path="/admin" element={<RoleRoute requiredRole="admin"><div>Admin Page</div></RoleRoute>} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByText('Admin Page')).toBeInTheDocument()
  })
})
