import { useEffect, useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export interface BookDetail {
  id: string;
  title: string;
  author: string;
  description?: string;
  condition?: string;
  imageUrl?: string;
  location?: { id: string; name: string; address?: string };
  status?: 'available' | 'reserved' | 'borrowed';
  ownerId?: string;
}

export const useBook = (id?: string | null) => {
  const [book, setBook] = useState<BookDetail | null>(null);
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  const fetchBook = async () => {
    if (!id) return;
    setLoading(true);
    try {
      const resp = await api.get<any>(`books/${id}`);
      const data = resp.data || {};
      // Normalize backend snake_case to camelCase used by components
      const normalized: any = {
        id: data.id,
        title: data.title,
        author: data.author,
        description: data.description,
        condition: data.condition,
        imageUrl: data.image_url || data.imageUrl || '',
        location: data.current_location || data.location || null,
        status: data.status,
        ownerId: data.owner_id || data.ownerId || null,
      };
      setBook(normalized);
    } catch (err: any) {
      const status = err?.response?.status;
      const msg = err?.response?.data?.message || err.message || 'Ошибка при загрузке книги';
      showNotification(msg, 'error');
      if (status === 404) {
        // keep book null and let caller handle navigation or UI
      }
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    let mounted = true;
    if (!id) return;
    // call fetchBook and ignore mounted flag inside since fetchBook manages state safely
    fetchBook();
    return () => { mounted = false; };
  }, [id]);

  return { book, loading, refresh: fetchBook };
};

export default useBook;
