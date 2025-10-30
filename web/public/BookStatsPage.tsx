import React, { useEffect, useState } from 'react';
import { Box, Typography, Button, Paper } from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import bookService from '../src/services/bookService';
import { useNotification } from '../contexts/NotificationContext';

const BookStatsPage: React.FC = () => {
  const { bookId } = useParams();
  const navigate = useNavigate();
  const { showNotification } = useNotification();
  const [stats, setStats] = useState<any | null>(null);
  const [book, setBook] = useState<any | null>(null);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      try {
        if (!bookId) return;
        const data = await (bookService as any).getBookStats(bookId);
        if (!mounted) return;
        setStats(data);
      } catch (e: any) {
        showNotification(e?.response?.data?.error || e.message || 'Ошибка при получении статистики', 'error');
      }
    };
    load();
    return () => { mounted = false; };
  }, [bookId, showNotification]);

  useEffect(() => {
    let mounted = true;
    const loadBook = async () => {
      try {
        if (!bookId) return;
        const b = await (bookService as any).getBook(bookId);
        if (!mounted) return;
        setBook(b);
      } catch (e: any) {
        // don't block stats on book fetch error
      }
    };
    loadBook();
    return () => { mounted = false; };
  }, [bookId]);

  return (
    <Box>
      <Typography variant="h5" gutterBottom>
        Статистика:&nbsp;
        {book?.title ? (
          <Box component="span" sx={{ fontWeight: 700 }}>&quot;{book.title}&quot;</Box>
        ) : (
          'Статистика'
        )}
      </Typography>

      {!stats && <Typography>Загрузка...</Typography>}

      {stats && (
        <Box sx={{ display: 'flex', gap: 2, flexWrap: 'wrap', alignItems: 'flex-start' }}>
          {[
            { title: 'Всего записей', value: stats.total_books ?? 0 },
            { title: 'Читатели', value: stats.total_unique_borrowers ?? 0 },
            { title: 'Статус', value: stats.current_status ?? '—' },
            { title: 'Состояние', value: stats.current_condition ?? '—' },
            { title: 'Первая локация', value: stats.first_location ?? '—' },
            { title: 'Последняя локация', value: stats.last_location ?? '—' },
          ].map((it) => (
            <Paper
              key={it.title}
              sx={(theme) => ({
                p: 1.25,
                px: 2,
                borderRadius: 2,
                flex: '0 0 auto',
                display: 'inline-flex',
                alignItems: 'center',
                gap: 1,
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                // on narrow screens stretch full width and space items apart
                [theme.breakpoints.down('sm')]: {
                  display: 'flex',
                  width: '100%',
                  justifyContent: 'space-between',
                },
              })}
              elevation={2}
            >
              <Typography variant="subtitle2" sx={{ color: 'text.secondary', display: 'inline' }}>{it.title}:</Typography>
              <Typography
                component="span"
                variant="subtitle2"
                sx={(theme) => ({
                  fontWeight: 700,
                  ml: 1,
                  display: 'inline-block',
                  maxWidth: 420,
                  overflow: 'hidden',
                  textOverflow: 'ellipsis',
                  // push value to the right on small screens
                  [theme.breakpoints.down('sm')]: {
                    ml: 'auto',
                  },
                })}
              >
                {it.value}
              </Typography>
            </Paper>
          ))}
        </Box>
      )}

      
    </Box>
  );
};

export default BookStatsPage;
