import api from './api';

export const getAdminStats = async () => {
  const resp = await api.get('admin/stats');
  return resp.data || {};
};

export const getAdminUsers = async () => {
  const resp = await api.get('admin/users');
  return resp.data || [];
};

export const updateUserRole = async (userId, role) => {
  const resp = await api.put(`admin/users/${userId}/role`, { role });
  return resp.data;
};

export default { getAdminStats, getAdminUsers, updateUserRole };
