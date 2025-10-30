import { useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export const useRequestBook = () => {
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const requestBook = async (bookId: string, locationId?: string) => {
    setLoading(true);
    try {
      const payload: any = { bookId };
      if (locationId) payload.locationId = locationId;
      const resp = await api.post('books/request', payload);
      showNotification(resp.data?.message || 'Книга забронирована', 'success');
      return resp.data;
    } catch (err: any) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка при бронировании';
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { requestBook, loading };
};

export default useRequestBook;
