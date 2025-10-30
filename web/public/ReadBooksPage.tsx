import React, { useEffect, useState } from 'react';
import { Typography, Box, CircularProgress, List, ListItem, ListItemText, Paper } from '@mui/material';
import { useNotification } from '../contexts/NotificationContext';
import bookService from '../src/services/bookService';
import BookCard from './BookCard';

const ReadBooksPage: React.FC = () => {
  const { showNotification } = useNotification();
  const [books, setBooks] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      try {
        const data = await bookService.getReadBooks();
        if (!mounted) return;
        setBooks(data || []);
      } catch (e: any) {
        showNotification(e?.response?.data?.error || e.message || 'Ошибка при загрузке', 'error');
      } finally {
        if (mounted) setLoading(false);
      }
    };
    load();
    return () => { mounted = false; };
  }, [showNotification]);

  if (loading) return <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}><CircularProgress /></Box>;

  return (
    <Box>
      <Typography variant="h5" gutterBottom>Прочитанное</Typography>
      <Paper sx={{ borderRadius: 3, overflow: 'hidden' }}>
        <List disablePadding>
          {books.map((b: any) => (
            <ListItem key={b.id} divider>
              <ListItemText primary={b.title} secondary={b.author} primaryTypographyProps={{ color: 'text.primary' }} secondaryTypographyProps={{ color: 'text.secondary' }} />
            </ListItem>
          ))}
        </List>
      </Paper>
      <Typography variant="body1" color="text.primary" align="center" sx={{ mt: 2, fontWeight: 500, opacity: 0.6 }}>
        Прочитанные и возвращённые книги
      </Typography>
    </Box>
  );
};

export default ReadBooksPage;