import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Button,
  CircularProgress,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  RadioGroup,
  FormControlLabel,
  Radio,
} from '@mui/material';
import { useNavigate } from 'react-router-dom';
import { useNotification } from '../contexts/NotificationContext';
import { useAuth } from './AuthContext';
import useDeleteBook from '../src/hooks/useDeleteBook';
import bookService from '../src/services/bookService';

type Book = {
  id: string;
  title: string;
  author: string;
  image_url?: string;
  description?: string;
  location?: { id: string; name: string; address?: string } | null;
};

const MyShelfPage: React.FC = () => {
  const navigate = useNavigate();
  const { showNotification } = useNotification();
  const { isAuthenticated } = useAuth();
  const { deleteBook, loading: deleting } = useDeleteBook();
  // return flow moved to dedicated page; related hooks removed from this component
  const [selectedBook, setSelectedBook] = useState<Book | null>(null);
  const [openDelete, setOpenDelete] = useState(false);
  const [deleteReason, setDeleteReason] = useState<string>('bad');
  const [deleteOtherText, setDeleteOtherText] = useState<string>('');
  const [loading, setLoading] = useState<boolean>(false);
  const [books, setBooks] = useState<Book[]>([]);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    const load = async () => {
      try {
        setLoading(true);
        const data = await bookService.getShelfBooks();
        if (!mounted) return;
        setBooks(data || []);
      } catch (e: any) {
        showNotification(e?.response?.data?.error || e.message || 'Ошибка при загрузке полки', 'error');
      } finally {
        if (mounted) setLoading(false);
      }
    };
    load();
    // listen for global updates (e.g. a book was borrowed from ReservedBooksPage)
    const onBooksUpdate = () => { load(); };
    window.addEventListener('books:update', onBooksUpdate);
    return () => { mounted = false; window.removeEventListener('books:update', onBooksUpdate); };
  }, [showNotification]);

  const handleOpenReturn = (book: Book) => {
    // Navigate to dedicated return page instead of opening dialog
    navigate(`/books/return/${book.id}`);
  };

  const handleOpenDelete = (book: Book) => {
    setSelectedBook(book);
    setOpenDelete(true);
  };

  return (
    <>
      <Box sx={{ p: '20px' }}>
        <Typography variant="h5" gutterBottom>
          Моя полка
        </Typography>
      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Box sx={{ display: 'flex', flexDirection: 'column', gap: 1 }}>
          {books.map((book) => (
            <Paper key={book.id} elevation={1} sx={{ p: '20px', borderRadius: 2 }}>
              <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2 }}>
                <Box sx={{ flex: 1 }} onClick={() => handleOpenReturn(book)}>
                  <Typography variant="body1" sx={{ fontWeight: 600 }}>{book.title}</Typography>
                </Box>
                <Box sx={{ display: 'flex', gap: 1 }}>
                  <Button
                    variant="outlined"
                    color="error"
                    size="small"
                    onClick={(e: React.MouseEvent<HTMLButtonElement>) => { e.stopPropagation(); handleOpenDelete(book); }}
                    disabled={deletingId === book.id}
                  >
                    Удалить
                  </Button>
                  <Button
                    variant="contained"
                    color="primary"
                    size="small"
                    onClick={(e: React.MouseEvent<HTMLButtonElement>) => { e.stopPropagation(); handleOpenReturn(book); }}
                    disabled={deletingId === book.id}
                  >
                    Вернуть
                  </Button>
                </Box>
              </Box>
            </Paper>
          ))}
        </Box>
      )}
        <Typography variant="body1" color="text.primary" align="center" sx={{ mt: 2, fontWeight: 500, opacity: 0.6 }}>
          Книги, которые сейчас у вас
        </Typography>
        {!isAuthenticated && (
          <Box sx={{ mt: 2, display: 'flex', gap: 1, alignItems: 'center' }}>
            <Typography variant="body2" color="text.secondary">Чтобы управлять своей полкой, пожалуйста, войдите.</Typography>
            <Button variant="text" size="small" onClick={() => navigate('/login')}>Войти</Button>
          </Box>
        )}

      {/* Note: delete dialog and per-item actions were intentionally removed to show only title blocks.
          Clicking a block opens the return flow via dedicated page. Deletion and other management
          are still available via the edit/management screens if needed. */}

      {/* Delete Dialog */}
      <Dialog fullWidth maxWidth="xs" open={openDelete} onClose={() => setOpenDelete(false)}>
        <DialogTitle>Причина удаления</DialogTitle>
        <DialogContent>
          <Box sx={{ mt: 1 }}>
            <RadioGroup value={deleteReason} name="delete-reason" onChange={(e) => setDeleteReason(e.target.value)}>
              <FormControlLabel value="bad" control={<Radio />} label="Плохое состояние" />
              <FormControlLabel value="lost" control={<Radio />} label="Утеряна" />
              <FormControlLabel value="other" control={<Radio />} label="Другое" />
            </RadioGroup>
            {deleteReason === 'other' && (
              <TextField
                label="Если другое, укажите"
                value={deleteOtherText}
                onChange={(e) => setDeleteOtherText(e.target.value)}
                fullWidth
                sx={{ mt: 1 }}
                error={deleteReason === 'other' && deleteOtherText.trim().length > 0 && deleteOtherText.trim().length < 5}
                helperText={deleteReason === 'other' && deleteOtherText.trim().length > 0 && deleteOtherText.trim().length < 5 ? 'Уточните причину минимум 5 символов' : ''}
              />
            )}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpenDelete(false)}>Отмена</Button>
          <Button
            variant="outlined"
            color="error"
            onClick={async () => {
              if (!selectedBook) return;
              if (deleteReason === 'other' && deleteOtherText.trim().length < 5) {
                return;
              }
              try {
                setDeletingId(selectedBook.id);
                const reason = deleteReason === 'other' ? deleteOtherText : deleteReason;
                await deleteBook(selectedBook.id, reason);
                setBooks((prev) => prev.filter((b) => b.id !== selectedBook.id));
                try { window.dispatchEvent(new Event('books:update')); } catch (e) {}
                setOpenDelete(false);
              } catch (err) {
                // handled by hook
              } finally {
                setDeletingId(null);
              }
            }}
            disabled={deletingId === selectedBook?.id || (deleteReason === 'other' && deleteOtherText.trim().length < 5)}
          >
            Удалить
          </Button>
        </DialogActions>
      </Dialog>
      </Box>
    </>
  );
};

export default MyShelfPage;