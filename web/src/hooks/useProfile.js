import { useState, useEffect, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import userService from '../services/userService';
import { useNotification } from '../../src/contexts/NotificationContext';

export const useProfile = () => {
  const [profile, setProfile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const { showNotification } = useNotification();

  const fetchProfile = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await userService.getProfile();
      setProfile(data);
      return data;
    } catch (err) {
      const status = err?.response?.status;
      const msg = err?.response?.data?.message || err.message || 'Ошибка получения профиля';
      setError(msg);
      showNotification(msg, 'error');
      if (status === 401) {
        try {
          localStorage.removeItem('token');
        } catch (e) {}
        navigate('/login');
      }
      return null;
    } finally {
      setLoading(false);
    }
  }, [navigate, showNotification]);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      fetchProfile();
    }
  }, [fetchProfile]);

  return { profile, loading, error, fetchProfile };
};

export default useProfile;
