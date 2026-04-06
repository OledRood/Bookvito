import { beforeEach, describe, expect, it, vi } from 'vitest'

const axiosPostMock = vi.fn()

vi.mock('axios', () => {
  const instance: any = vi.fn(() => Promise.resolve({ data: {} }))
  instance.defaults = {
    baseURL: 'http://127.0.0.1:18080/api/v1/',
    headers: { common: {} },
  }
  instance.interceptors = {
    request: {
      use: (handler: any) => {
        instance.__requestHandler = handler
      },
    },
    response: {
      use: (fulfilled: any, rejected: any) => {
        instance.__responseFulfilled = fulfilled
        instance.__responseRejected = rejected
      },
    },
  }

  return {
    default: {
      create: vi.fn(() => instance),
      post: axiosPostMock,
    },
  }
})

describe('api service', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.resetModules()
    vi.clearAllMocks()
  })

  it('attaches access token from localStorage', async () => {
    localStorage.setItem('token', 'access-token')
    const api = (await import('./api')).default as any

    const config = await api.__requestHandler({
      headers: {},
      method: 'get',
      url: 'books/list',
    })

    expect(config.headers.Authorization).toBe('Bearer access-token')
  })

  it('refreshes token and retries the original request on 401', async () => {
    localStorage.setItem('token', 'old-token')
    localStorage.setItem('accessToken', 'old-token')
    localStorage.setItem('refreshToken', 'refresh-token')

    const api = (await import('./api')).default as any
    const retriedResponse = { data: { ok: true } }
    api.mockResolvedValueOnce(retriedResponse)
    axiosPostMock.mockResolvedValueOnce({
      data: {
        access_token: 'new-token',
        refresh_token: 'new-refresh',
      },
    })

    const result = await api.__responseRejected({
      response: { status: 401 },
      config: { headers: {}, url: 'books/my' },
      message: 'unauthorized',
    })

    expect(axiosPostMock).toHaveBeenCalled()
    expect(api).toHaveBeenCalled()
    expect(localStorage.getItem('token')).toBe('new-token')
    expect(localStorage.getItem('refreshToken')).toBe('new-refresh')
    expect(result).toEqual(retriedResponse)
  })

  it('clears tokens and dispatches unauthorized when refresh token is missing', async () => {
    localStorage.setItem('token', 'stale-token')
    const unauthorizedListener = vi.fn()
    window.addEventListener('auth:unauthorized', unauthorizedListener)

    const api = (await import('./api')).default as any

    await expect(
      api.__responseRejected({
        response: { status: 401 },
        config: { headers: {}, url: 'books/my' },
        message: 'unauthorized',
      }),
    ).rejects.toBeTruthy()

    expect(localStorage.getItem('token')).toBeNull()
    expect(unauthorizedListener).toHaveBeenCalled()
    window.removeEventListener('auth:unauthorized', unauthorizedListener)
  })

  it('forces logout when refresh request returns no token', async () => {
    localStorage.setItem('token', 'stale-token')
    localStorage.setItem('accessToken', 'stale-token')
    localStorage.setItem('refreshToken', 'stale-refresh')
    const unauthorizedListener = vi.fn()
    window.addEventListener('auth:unauthorized', unauthorizedListener)

    const api = (await import('./api')).default as any
    axiosPostMock.mockResolvedValueOnce({ data: {} })

    await expect(
      api.__responseRejected({
        response: { status: 401 },
        config: { headers: {}, url: 'books/my' },
        message: 'unauthorized',
      }),
    ).rejects.toBeTruthy()

    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('accessToken')).toBeNull()
    expect(localStorage.getItem('refreshToken')).toBeNull()
    expect(unauthorizedListener).toHaveBeenCalled()
    window.removeEventListener('auth:unauthorized', unauthorizedListener)
  })

  it('forces logout when retried request still returns 401', async () => {
    localStorage.setItem('token', 'stale-token')
    localStorage.setItem('accessToken', 'stale-token')
    localStorage.setItem('refreshToken', 'stale-refresh')
    const unauthorizedListener = vi.fn()
    window.addEventListener('auth:unauthorized', unauthorizedListener)

    const api = (await import('./api')).default as any

    await expect(
      api.__responseRejected({
        response: { status: 401 },
        config: { headers: {}, url: 'books/my', _retry: true },
        message: 'unauthorized',
      }),
    ).rejects.toBeTruthy()

    expect(localStorage.getItem('token')).toBeNull()
    expect(localStorage.getItem('accessToken')).toBeNull()
    expect(localStorage.getItem('refreshToken')).toBeNull()
    expect(unauthorizedListener).toHaveBeenCalled()
    window.removeEventListener('auth:unauthorized', unauthorizedListener)
  })
})
