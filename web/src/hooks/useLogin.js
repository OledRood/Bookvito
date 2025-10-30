import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import userService from '../services/userService';
import { useNotification } from '../../src/contexts/NotificationContext';

export const useLogin = () => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const navigate = useNavigate();
  const { showNotification } = useNotification();

  const login = async (email, password) => {
    setLoading(true);
    setError(null);
    try {
      const data = await userService.login(email, password);
      try {
        // eslint-disable-next-line no-console
        console.debug('[useLogin] normalized login response', data);
      } catch (e) {}
      if (data && (data.token || data.accessToken)) {
        const t = data.token || data.accessToken;
        // write both common keys used across the app
        localStorage.setItem('token', t);
        localStorage.setItem('accessToken', t);
        // optionally store user
        if (data.user) localStorage.setItem('user', JSON.stringify(data.user));
        showNotification('Успешный вход', 'success');
        navigate('/');
        return true;
      }
      setError('Invalid response');
      showNotification('Неверный ответ от сервера', 'error');
      return false;
    } catch (err) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка входа';
      setError(msg);
      showNotification(msg, 'error');
      return false;
    } finally {
      setLoading(false);
    }
  };

  return { login, loading, error };
};

export default useLogin;
