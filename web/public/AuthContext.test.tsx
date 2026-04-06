import React from 'react'
import { act, render, screen, waitFor } from '@testing-library/react'
import { vi } from 'vitest'
import { AuthProvider, useAuth } from './AuthContext'
import userService from '../src/services/userService'

vi.mock('../src/services/userService', () => ({
  default: {
    login: vi.fn(),
    register: vi.fn(),
    getProfile: vi.fn(),
    logout: vi.fn(),
  },
}))

const AuthProbe = () => {
  const {
    user,
    isAuthenticated,
    isAdmin,
    isModer,
    hasRole,
    login,
    register,
    logout,
    saveTokens,
    clearTokens,
  } = useAuth()

  return (
    <div>
      <span data-testid="auth">{String(isAuthenticated)}</span>
      <span data-testid="name">{user?.name || 'anonymous'}</span>
      <span data-testid="role">{user?.role || 'none'}</span>
      <span data-testid="is-admin">{String(isAdmin)}</span>
      <span data-testid="is-moder">{String(isModer)}</span>
      <span data-testid="has-role-moder">{String(hasRole('moder'))}</span>
      <button
        onClick={() => {
          login('reader@example.com', 'password123').then((result) => {
            ;(window as any).__lastLoginResult = result
          })
        }}
      >
        login
      </button>
      <button
        onClick={() => {
          register('Reader', 'reader@example.com', 'password123').then((result) => {
            ;(window as any).__lastRegisterResult = result
          })
        }}
      >
        register
      </button>
      <button onClick={() => logout()}>logout</button>
      <button
        onClick={() =>
          saveTokens('saved-token', 'saved-refresh', {
            id: 'saved-id',
            name: 'Saved User',
            email: 'saved@bookvito.test',
            role: 'moder',
            avatar: 'avatar1.png',
          })
        }
      >
        save
      </button>
      <button onClick={() => clearTokens()}>clear</button>
    </div>
  )
}

describe('AuthContext', () => {
  const mockedUserService = vi.mocked(userService, { deep: true })

  beforeEach(() => {
    localStorage.clear()
    ;(window as any).__lastLoginResult = undefined
    ;(window as any).__lastRegisterResult = undefined
    mockedUserService.login.mockReset()
    mockedUserService.register.mockReset()
    mockedUserService.getProfile.mockReset()
    mockedUserService.logout.mockReset()
    mockedUserService.logout.mockResolvedValue(undefined as any)
  })

  it('restores auth state from localStorage and clears it on unauthorized event', async () => {
    localStorage.setItem('accessToken', 'token')
    localStorage.setItem('refreshToken', 'refresh')
    localStorage.setItem('user', JSON.stringify({ id: '1', name: 'Reader', role: 'user' }))

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    )

    await waitFor(() => {
      expect(screen.getByTestId('auth')).toHaveTextContent('true')
      expect(screen.getByTestId('name')).toHaveTextContent('Reader')
    })

    act(() => {
      window.dispatchEvent(new CustomEvent('auth:unauthorized'))
    })

    await waitFor(() => {
      expect(screen.getByTestId('auth')).toHaveTextContent('false')
      expect(screen.getByTestId('name')).toHaveTextContent('anonymous')
    })
  })

  it('logs in, fetches profile and sets role flags', async () => {
    mockedUserService.login.mockResolvedValue({
      token: 'new-token',
      refreshToken: 'new-refresh',
    } as any)
    mockedUserService.getProfile.mockResolvedValue({
      id: '1',
      name: 'Admin Reader',
      email: 'admin@bookvito.test',
      role: 'admin',
      avatar: 'avatar1.png',
    } as any)

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    )

    act(() => {
      screen.getByRole('button', { name: 'login' }).click()
    })

    await waitFor(() => {
      expect((window as any).__lastLoginResult).toBe(true)
      expect(screen.getByTestId('auth')).toHaveTextContent('true')
      expect(screen.getByTestId('name')).toHaveTextContent('Admin Reader')
      expect(screen.getByTestId('is-admin')).toHaveTextContent('true')
      expect(screen.getByTestId('is-moder')).toHaveTextContent('true')
      expect(screen.getByTestId('has-role-moder')).toHaveTextContent('true')
    })

    expect(localStorage.getItem('accessToken')).toBe('new-token')
    expect(localStorage.getItem('token')).toBe('new-token')
    expect(localStorage.getItem('refreshToken')).toBe('new-refresh')
  })

  it('returns false on login when token is absent in response', async () => {
    mockedUserService.login.mockResolvedValue({} as any)

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    )

    act(() => {
      screen.getByRole('button', { name: 'login' }).click()
    })

    await waitFor(() => {
      expect((window as any).__lastLoginResult).toBe(false)
      expect(screen.getByTestId('auth')).toHaveTextContent('false')
      expect(screen.getByTestId('name')).toHaveTextContent('anonymous')
    })
  })

  it('registers and restores user from profile endpoint', async () => {
    mockedUserService.register.mockResolvedValue({
      success: true,
      data: {
        token: 'reg-token',
        refreshToken: 'reg-refresh',
      },
    } as any)
    mockedUserService.getProfile.mockResolvedValue({
      id: '2',
      name: 'Moder Reader',
      email: 'moder@bookvito.test',
      role: 'moder',
      avatar: 'avatar2.png',
    } as any)

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    )

    act(() => {
      screen.getByRole('button', { name: 'register' }).click()
    })

    await waitFor(() => {
      expect((window as any).__lastRegisterResult).toEqual({ success: true })
      expect(screen.getByTestId('auth')).toHaveTextContent('true')
      expect(screen.getByTestId('name')).toHaveTextContent('Moder Reader')
      expect(screen.getByTestId('role')).toHaveTextContent('moder')
      expect(screen.getByTestId('is-admin')).toHaveTextContent('false')
      expect(screen.getByTestId('is-moder')).toHaveTextContent('true')
    })
  })

  it('returns register error payload when registration fails', async () => {
    mockedUserService.register.mockResolvedValue({
      success: false,
      error: 'email already exists',
    } as any)

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    )

    act(() => {
      screen.getByRole('button', { name: 'register' }).click()
    })

    await waitFor(() => {
      expect((window as any).__lastRegisterResult).toEqual({
        success: false,
        error: 'email already exists',
      })
      expect(screen.getByTestId('auth')).toHaveTextContent('false')
    })
  })

  it('clears local state via save/clear/logout actions', async () => {
    localStorage.setItem('accessToken', 'initial-token')
    localStorage.setItem('refreshToken', 'initial-refresh')
    localStorage.setItem('user', JSON.stringify({ id: '1', name: 'Reader', role: 'user' }))

    render(
      <AuthProvider>
        <AuthProbe />
      </AuthProvider>,
    )

    act(() => {
      screen.getByRole('button', { name: 'save' }).click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('auth')).toHaveTextContent('true')
      expect(screen.getByTestId('name')).toHaveTextContent('Saved User')
      expect(localStorage.getItem('accessToken')).toBe('saved-token')
      expect(localStorage.getItem('refreshToken')).toBe('saved-refresh')
    })

    act(() => {
      screen.getByRole('button', { name: 'clear' }).click()
    })

    await waitFor(() => {
      expect(screen.getByTestId('auth')).toHaveTextContent('false')
      expect(localStorage.getItem('accessToken')).toBeNull()
      expect(localStorage.getItem('refreshToken')).toBeNull()
      expect(localStorage.getItem('user')).toBeNull()
    })

    act(() => {
      screen.getByRole('button', { name: 'save' }).click()
    })
    act(() => {
      screen.getByRole('button', { name: 'logout' }).click()
    })

    await waitFor(() => {
      expect(mockedUserService.logout).toHaveBeenCalled()
      expect(screen.getByTestId('auth')).toHaveTextContent('false')
      expect(localStorage.getItem('accessToken')).toBeNull()
      expect(localStorage.getItem('refreshToken')).toBeNull()
      expect(localStorage.getItem('user')).toBeNull()
    })
  })
})
