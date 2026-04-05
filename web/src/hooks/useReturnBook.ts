import { useState } from 'react';
import { returnBook as returnBookRequest } from '../services/bookService';
import { useNotification } from '../../contexts/NotificationContext';

export const useReturnBook = () => {
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const returnBook = async (payload: {
    bookId: string;
    title: string;
    author: string;
    condition?: string;
    description?: string;
    currentLocationId?: string;
    imageUrl?: string;
  }, options?: { silent?: boolean }) => {
    setLoading(true);
    try {
      const resp = await returnBookRequest(payload);
      if (!options?.silent) {
        showNotification(resp?.message || 'Книга возвращена', 'success');
      }
      return resp;
    } catch (err: any) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка при возврате книги';
      if (!options?.silent) {
        showNotification(msg, 'error');
      }
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { returnBook, loading };
};

export default useReturnBook;
