import { useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export const useCreateBook = () => {
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const createBook = async (payload: any) => {
    setLoading(true);
    try {
      const resp = await api.post('books/create', payload);
      showNotification(resp.data?.message || 'Книга создана', 'success');
      return resp.data;
    } catch (err: any) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка при создании книги';
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { createBook, loading };
};

export default useCreateBook;
