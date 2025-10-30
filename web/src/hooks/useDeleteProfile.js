import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import userService from '../services/userService';
import { useNotification } from '../../src/contexts/NotificationContext';

export const useDeleteProfile = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const { showNotification } = useNotification();

  const remove = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await userService.deleteProfile();
      showNotification(data?.message || 'Профиль удалён', 'success');
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      navigate('/login');
      return data;
    } catch (err) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка удаления профиля';
      setError(msg);
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { remove, loading, error };
};

export default useDeleteProfile;
