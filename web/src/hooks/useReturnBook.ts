import { useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export const useReturnBook = () => {
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const returnBook = async (payload: { bookId: string; title: string; author: string; condition?: string; description?: string; currentLocationId?: string }) => {
    setLoading(true);
    try {
      // map to backend expected keys (snake_case)
      const body: any = {
        book_id: payload.bookId,
        title: payload.title,
        author: payload.author,
      };
      if (payload.condition) body.condition = payload.condition;
      if (payload.description) body.description = payload.description;
      if (payload.currentLocationId) body.current_location_id = payload.currentLocationId;

      const resp = await api.put('books/return', body);
      showNotification(resp.data?.message || 'Книга возвращена', 'success');
      return resp.data;
    } catch (err: any) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка при возврате книги';
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { returnBook, loading };
};

export default useReturnBook;
