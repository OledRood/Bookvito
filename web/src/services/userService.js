import api from './api';

export const login = async (email, password) => {
  const { data } = await api.post('users/login', { email, password });
  try {
    // eslint-disable-next-line no-console
    console.debug('[userService] login raw response', data);
  } catch (e) {}
  // normalize backend naming (access_token / refresh_token) to token / refreshToken
  return {
    token: data.token || data.access_token || data.accessToken,
    refreshToken: data.refreshToken || data.refresh_token,
    user: data.user || data.user,
    raw: data,
  };
};

export const register = async (name, email, password) => {
  try {
    const { data } = await api.post('users/registration', { name, email, password });
    try {
      // eslint-disable-next-line no-console
      console.debug('[userService] register raw response', data);
    } catch (e) {}

    const normalized = {
      token: data.token || data.access_token || data.accessToken,
      refreshToken: data.refreshToken || data.refresh_token,
      user: data.user || data.user,
      raw: data,
    };
    return { success: true, data: normalized };
  } catch (err) {
    // try to extract meaningful message from axios error
    const message = err && err.response && err.response.data && (err.response.data.error || err.response.data.message)
      ? err.response.data.error || err.response.data.message
      : (err && err.message) || 'Registration failed';
    try {
      // eslint-disable-next-line no-console
      console.error('[userService] register error', { message, err });
    } catch (e) {}
    return { success: false, error: message, rawError: err };
  }
};

export const getProfile = async () => {
  const { data } = await api.get('users/me');
  return data;
};

export const updateProfile = async (payload) => {
  const { data } = await api.put('users/me', payload);
  return data;
};

export const deleteProfile = async () => {
  const { data } = await api.delete('users/me');
  return data;
};

export default {
  login,
  register,
  getProfile,
  updateProfile,
  deleteProfile,
};
