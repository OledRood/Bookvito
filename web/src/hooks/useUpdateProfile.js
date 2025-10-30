import { useState } from 'react';
import userService from '../services/userService';
import { useNotification } from '../../src/contexts/NotificationContext';

export const useUpdateProfile = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const { showNotification } = useNotification();

  const update = async (payload) => {
    setLoading(true);
    setError(null);
    try {
      const data = await userService.updateProfile(payload);
      showNotification(data?.message || 'Профиль обновлен', 'success');
      return data;
    } catch (err) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка обновления профиля';
      setError(msg);
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { update, loading, error };
};

export default useUpdateProfile;
