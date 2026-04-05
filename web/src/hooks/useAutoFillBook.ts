import { useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export const useAutoFillBook = () => {
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const autoFill = async (query: string) => {
    setLoading(true);
    try {
      const resp = await api.get('books/auto-fill', {
        params: { q: query },
      });
      return resp.data;
    } catch (err: any) {
      const msg = err?.response?.data?.error || err?.message || 'Не удалось автозаполнить данные';
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { autoFill, loading };
};

export default useAutoFillBook;
