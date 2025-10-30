import axios from 'axios';

// plain axios instance used by services
// Read base URL from Vite env (VITE_API_BASE_URL) when available, otherwise
// fall back to localhost:8080 which is the backend default in this project.
const BASE_URL = (typeof import.meta !== 'undefined' && import.meta.env && import.meta.env.VITE_API_BASE_URL)
  ? import.meta.env.VITE_API_BASE_URL
  : 'http://localhost:8080/api/v1/';

const api = axios.create({
  baseURL: BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 15000,
});

// request interceptor: attach token from localStorage
api.interceptors.request.use((config) => {
  try {
    // support both `token` (new) and `accessToken` (existing public/AuthContext)
    const token = localStorage.getItem('token') || localStorage.getItem('accessToken');
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
  } catch (e) {
    // ignore
  }
  // debug: print outgoing request in browser console
  try {
    // eslint-disable-next-line no-console
    console.debug('[api] request', { method: config.method, url: config.url, headers: config.headers });
  } catch (e) {}
  return config;
});

// response interceptor: handle 401 globally
api.interceptors.response.use(
  (resp) => resp,
  (error) => {
    const { response } = error || {};
    // debug: log response error
    try {
      // eslint-disable-next-line no-console
      console.error('[api] response error', { message: error.message, response: response && { status: response.status, data: response.data } });
    } catch (e) {}
    if (response && response.status === 401) {
      const originalRequest = error.config || {};
      // don't try to refresh for the refresh endpoint itself
      if (originalRequest._retry) {
        try {
          localStorage.removeItem('token');
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
        } catch (e) {}
        window.dispatchEvent(new CustomEvent('auth:unauthorized'));
        return Promise.reject(error);
      }

      // queue for single refresh call
      if (!api.isRefreshing) api.isRefreshing = false;
      if (!api.refreshSubscribers) api.refreshSubscribers = [];

      const subscribeTokenRefresh = (cb) => {
        api.refreshSubscribers.push(cb);
      };

      const onRefreshed = (token) => {
        api.refreshSubscribers.forEach((cb) => cb(token));
        api.refreshSubscribers = [];
      };

      const refreshToken = localStorage.getItem('refreshToken');
      if (!refreshToken) {
        // no refresh token available, force logout
        try {
          localStorage.removeItem('token');
          localStorage.removeItem('accessToken');
          localStorage.removeItem('refreshToken');
        } catch (e) {}
        window.dispatchEvent(new CustomEvent('auth:unauthorized'));
        return Promise.reject(error);
      }

      if (api.isRefreshing) {
        // already refreshing, return a promise that resolves once refreshed
        return new Promise((resolve, reject) => {
          subscribeTokenRefresh((token) => {
            // set the new header and retry
            originalRequest.headers = originalRequest.headers || {};
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(api(originalRequest));
          });
        });
      }

      api.isRefreshing = true;

      // attempt refresh
      return axios.post(`${api.defaults.baseURL}users/refresh`, { refresh_token: refreshToken })
        .then((res) => {
          const data = res.data || {};
          const newToken = data.access_token || data.token || data.accessToken;
          const newRefresh = data.refresh_token || data.refreshToken;
          if (newToken) {
            try {
              localStorage.setItem('token', newToken);
              localStorage.setItem('accessToken', newToken);
              if (newRefresh) localStorage.setItem('refreshToken', newRefresh);
            } catch (e) {}
            api.defaults.headers.common.Authorization = `Bearer ${newToken}`;
            onRefreshed(newToken);
            // mark request as retried
            originalRequest._retry = true;
            originalRequest.headers = originalRequest.headers || {};
            originalRequest.headers.Authorization = `Bearer ${newToken}`;
            return api(originalRequest);
          }
          // if refresh didn't return token, force logout
          try {
            localStorage.removeItem('token');
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
          } catch (e) {}
          window.dispatchEvent(new CustomEvent('auth:unauthorized'));
          return Promise.reject(error);
        })
        .catch((refreshError) => {
          try {
            localStorage.removeItem('token');
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
          } catch (e) {}
          window.dispatchEvent(new CustomEvent('auth:unauthorized'));
          return Promise.reject(refreshError);
        })
        .finally(() => {
          api.isRefreshing = false;
        });
    }
    return Promise.reject(error);
  }
);

export default api;
