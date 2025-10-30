import { useEffect, useState } from 'react';
import api from '../services/api';
import { useNotification } from '../../contexts/NotificationContext';

export interface BookSummary {
  id: string;
  title: string;
  author: string;
  imageUrl?: string;
  condition?: string;
  locationId?: string;
}

export const useBooksList = () => {
  const [books, setBooks] = useState<BookSummary[] | null>(null);
  const [loading, setLoading] = useState(false);
  const { showNotification } = useNotification();

  useEffect(() => {
    let mounted = true;
    const fetch = async () => {
      setLoading(true);
      try {
        const resp = await api.get<any[]>('books/summary');
        if (!mounted) return;
        // Normalize snake_case from backend (image_url) to camelCase used in the UI (imageUrl)
        const normalized = (resp.data || []).map((b: any) => ({
          id: b.id,
          title: b.title,
          author: b.author,
          imageUrl: b.image_url || b.imageUrl || '',
          condition: b.condition,
          locationId: b.current_location_id || b.locationId || null,
        }));
        setBooks(normalized);
      } catch (err: any) {
        const msg = err?.response?.data?.message || err.message || 'Ошибка при загрузке книг';
        showNotification(msg, 'error');
        setBooks([]);
      } finally {
        if (mounted) setLoading(false);
      }
    };
    fetch();
    return () => { mounted = false; };
  }, [showNotification]);

  return { books, loading };
};

export default useBooksList;
