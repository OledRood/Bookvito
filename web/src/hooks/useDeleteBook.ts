import { useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export const useDeleteBook = () => {
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const deleteBook = async (id: string, reason?: string) => {
    setLoading(true);
    try {
      // backend accepts DELETE /books/delete with JSON body { book_id: <id>, reason: <optional> }
      const body: any = { book_id: id };
      if (reason) body.reason = reason;
      const resp = await api.delete('books/delete', { data: body });
      showNotification(resp.data?.message || 'Книга удалена', 'success');
      return resp.data;
    } catch (err: any) {
      const msg = err?.response?.data?.message || err.message || 'Ошибка при удалении книги';
      showNotification(msg, 'error');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return { deleteBook, loading };
};

export default useDeleteBook;
