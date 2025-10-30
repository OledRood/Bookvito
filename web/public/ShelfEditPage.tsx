import React, { useState, useEffect } from 'react';
import { Box, Typography, TextField, Button, MenuItem } from '@mui/material';
import { useParams, useNavigate } from 'react-router-dom';
import bookService from '../src/services/bookService';
import { useNotification } from '../contexts/NotificationContext';

const ShelfEditPage: React.FC = () => {
  const { bookId } = useParams();
  const navigate = useNavigate();
  const { showNotification } = useNotification();
  const [loading, setLoading] = useState(false);
  const [book, setBook] = useState<any | null>(null);
  const [description, setDescription] = useState('');
  const [state, setState] = useState('good');
  const [location, setLocation] = useState('');

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      try {
  const data = await (bookService as any).getBook(String(bookId));
        if (!mounted) return;
        setBook(data);
        setDescription(data?.description || '');
        setLocation(data?.location || '');
        setState(data?.state || 'good');
      } catch (e: any) {
        showNotification('Не удалось загрузить книгу', 'error');
      } finally {
        if (mounted) setLoading(false);
      }
    };
    if (bookId) load();
    return () => { mounted = false; };
  }, [bookId, showNotification]);

  const handleReturn = () => navigate(`/books/return/${bookId}`);

  const handleSave = async () => {
    try {
  await (bookService as any).updateShelf(String(bookId), { description, state, location });
      showNotification('Сохранено', 'success');
    } catch (e: any) {
      showNotification('Ошибка при сохранении', 'error');
    }
  };

  if (!book) return <Typography>Загрузка...</Typography>;

  return (
    <Box sx={{ maxWidth: 640 }}>
      <Typography variant="h5" gutterBottom>Редактирование книги на полке</Typography>
      <Box sx={{ mt: 2 }}>
        <TextField fullWidth label="Название" value={book.title} disabled sx={{ mb: 2 }} />
        <TextField fullWidth label="Автор" value={book.author} disabled sx={{ mb: 2 }} />
        <TextField fullWidth label="Описание" multiline minRows={3} value={description} onChange={(e) => setDescription(e.target.value)} sx={{ mb: 2 }} />
        <TextField fullWidth select label="Состояние" value={state} onChange={(e) => setState(e.target.value)} sx={{ mb: 2 }}>
          <MenuItem value="good">Хорошее</MenuItem>
          <MenuItem value="worn">Плохое</MenuItem>
          <MenuItem value="lost">Утеряна</MenuItem>
        </TextField>
        <TextField fullWidth label="Локация" value={location} onChange={(e) => setLocation(e.target.value)} sx={{ mb: 2 }} />

        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button variant="contained" color="primary" onClick={handleSave}>Сохранить</Button>
          <Button variant="outlined" onClick={handleReturn}>Вернуть</Button>
          <Button variant="text" color="error" onClick={() => navigate('/books/shelf')}>Удалить</Button>
        </Box>
      </Box>
    </Box>
  );
};

export default ShelfEditPage;
