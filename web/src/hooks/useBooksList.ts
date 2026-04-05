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

let booksCache: BookSummary[] | null = null;
let booksPromise: Promise<BookSummary[]> | null = null;

const fetchBookSummaries = async (): Promise<BookSummary[]> => {
  if (booksCache) return booksCache;
  if (booksPromise) return booksPromise;

  booksPromise = api.get<any[]>('books/summary')
    .then((resp) => {
      const rawData = Array.isArray(resp.data)
        ? resp.data
        : Array.isArray(resp.data?.items)
          ? resp.data.items
          : [];

      const normalized = rawData.map((b: any) => ({
        id: b.id,
        title: b.title,
        author: b.author,
        imageUrl: b.image_url || b.imageUrl || '',
        condition: b.condition,
        locationId: b.current_location_id || b.locationId || null,
      }));

      booksCache = normalized;
      return normalized;
    })
    .finally(() => {
      booksPromise = null;
    });

  return booksPromise;
};

export const useBooksList = () => {
  const [books, setBooks] = useState<BookSummary[] | null>(null);
  const [loading, setLoading] = useState(!booksCache);
  const { showNotification } = useNotification();

  useEffect(() => {
    let mounted = true;

    if (booksCache) {
      setBooks(booksCache);
      setLoading(false);
      return () => { mounted = false; };
    }

    const fetch = async () => {
      setLoading(true);
      try {
        const normalized = await fetchBookSummaries();
        if (!mounted) return;
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
