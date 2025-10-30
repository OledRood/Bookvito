import React, { useEffect, useState } from 'react';
import { Typography, Box, CircularProgress, List, ListItem, ListItemText, IconButton, Paper } from '@mui/material';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { useNotification } from '../contexts/NotificationContext';
import bookService from '../src/services/bookService';
import { Link as RouterLink } from 'react-router-dom';

// MyBooksPage: render as simple list tiles (no images) with chevron -> stats placeholder

const MyBooksPage: React.FC = () => {
  const { showNotification } = useNotification();
  const [books, setBooks] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      try {
        const data = await bookService.getMyBooks();
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
      <Typography variant="h5" gutterBottom>Мои книги</Typography>
      <Paper sx={{ borderRadius: 3, overflow: 'hidden' }}>
        <List disablePadding>
          {books.map((b: any, idx: number) => (
            <ListItem
              key={b.id}
              divider
              sx={{
                ...(idx === 0 ? { pt: 2 } : {}),
                ...(idx === books.length - 1 ? { pb: 2 } : {}),
              }}
              secondaryAction={
                <IconButton edge="end" component={RouterLink} to={`/books/my/stats/${b.id}`} aria-label="stats">
                  <ChevronRightIcon />
                </IconButton>
              }
              component={RouterLink}
              to={`/books/my/stats/${b.id}`}
            >
              <ListItemText primary={b.title} secondary={b.author} primaryTypographyProps={{ color: 'text.primary' }} secondaryTypographyProps={{ color: 'text.secondary' }} />
            </ListItem>
          ))}
        </List>
      </Paper>
      <Typography variant="body1" color="text.primary" align="center" sx={{ mt: 2, fontWeight: 500, opacity: 0.6 }}>
        Все книги, с которыми вы взаимодействовали
      </Typography>
    </Box>
  );
};

export default MyBooksPage;