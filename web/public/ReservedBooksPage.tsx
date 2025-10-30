import React, { useEffect, useState } from 'react';
import { Typography, Box, CircularProgress, Card, CardContent, CardActions, Button, IconButton, Dialog, DialogTitle, DialogContent, DialogActions, Paper } from '@mui/material';
import MapIcon from '@mui/icons-material/Map';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import AutorenewIcon from '@mui/icons-material/Autorenew';
import { useNotification } from '../contexts/NotificationContext';
import bookService from '../src/services/bookService';

const ReservedBooksPage: React.FC = () => {
  const { showNotification } = useNotification();
  const [books, setBooks] = useState<any[]>([]);
  const [loading, setLoading] = useState(false);
  const [qrOpen, setQrOpen] = useState(false);
  const [selected, setSelected] = useState<any>(null);


  useEffect(() => {
    let mounted = true;
    const load = async () => {
      setLoading(true);
      try {
        const data = await bookService.getReservedBooks();
        if (!mounted) return;
        // backend now returns exchanges: { id, user_id, book, expires_at, location }
        // normalize into items used by this component
        const formatLoc = (loc: any, bookLoc: any) => {
          // In Reserved page we only show the address part (no location name)
          const l = loc || bookLoc || null;
          if (!l) return '';
          return l.address || l.name || '';
        };

        const normalized = (data || []).map((ex: any) => {
          const book = ex.book || ex.Book || {};
          const expiresAt = ex.expires_at ? new Date(ex.expires_at) : null;
          let expiresLabel = 'истекло';
          if (expiresAt) {
            const diff = Math.max(0, Math.floor((expiresAt.getTime() - Date.now()) / 1000));
            const days = Math.floor(diff / 86400);
            const hours = Math.floor((diff % 86400) / 3600);
            if (diff <= 0) expiresLabel = 'истекла';
            else if (days > 0) expiresLabel = `через ${days}д ${hours}ч`;
            else expiresLabel = `через ${hours}ч`;
          }
          return {
            id: book.id,
            exchangeId: ex.id,
            title: book.title,
            location: formatLoc(ex.location, book.current_location),
            expiresLabel,
            raw: ex,
          };
        });
        setBooks(normalized || []);
      } catch (e: any) {
        showNotification(e?.response?.data?.error || e.message || 'Ошибка при загрузке', 'error');
      } finally {
        if (mounted) setLoading(false);
      }
    };
    load();
    return () => { mounted = false; };
  }, []);

  if (loading) return <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}><CircularProgress /></Box>;

  const handleDelete = async (id: string) => {
    // cancel reservation
    try {
      await bookService.cancelReservation(id);
      showNotification('Бронь отменена', 'success');
      // reload
      const data = await bookService.getReservedBooks();
      const formatLoc = (loc: any, bookLoc: any) => {
        const l = loc || bookLoc || null;
        if (!l) return '';
        return l.address || l.name || '';
      };
      setBooks((data || []).map((ex: any) => {
        const book = ex.book || ex.Book || {};
        const expiresAt = ex.expires_at ? new Date(ex.expires_at) : null;
        let expiresLabel = 'истекло';
        if (expiresAt) {
          const diff = Math.max(0, Math.floor((expiresAt.getTime() - Date.now()) / 1000));
          const days = Math.floor(diff / 86400);
          const hours = Math.floor((diff % 86400) / 3600);
          if (diff <= 0) expiresLabel = 'истекла';
          else if (days > 0) expiresLabel = `через ${days}д ${hours}ч`;
          else expiresLabel = `через ${hours}ч`;
        }
        return {
          id: book.id,
          exchangeId: ex.id,
          title: book.title,
          location: formatLoc(ex.location, book.current_location),
          expiresLabel,
          raw: ex,
        };
      }));
      // notify other parts of app that reservation state changed
      try { window.dispatchEvent(new Event('books:update')); } catch (e) { /* noop */ }
    } catch (e: any) {
      showNotification(e?.response?.data?.error || e.message || 'Ошибка при отмене', 'error');
    }
  };

  

  const openQr = (item: any) => {
    setSelected(item);
    setQrOpen(true);
  };

  const closeQr = () => {
    setQrOpen(false);
    setSelected(null);
  };

  const handleReceive = async () => {
    if (!selected) return;
    try {
      // call borrow endpoint
      await bookService.borrowBook(selected.id);
      showNotification('Книга помечена как взята', 'success');
      // close dialog immediately
      closeQr();
      // optimistically remove the just-received book from the reserved list so UI matches server
      setBooks((prev) => prev.filter((b) => {
        // match by exchangeId when available, otherwise by book id
        if (selected.exchangeId && b.exchangeId) return b.exchangeId !== selected.exchangeId;
        return b.id !== selected.id;
      }));
      // notify other parts (e.g. MyShelfPage) to refresh their data
      try { window.dispatchEvent(new Event('books:update')); } catch (e) { /* noop */ }
    } catch (e: any) {
      showNotification(e?.response?.data?.error || e.message || 'Ошибка при получении', 'error');
    }
  };

  const handleOnMap = (item: any) => {
    // Open Yandex.Maps in a new tab with the reservation location.
    // Try the most specific fields first: raw.location.path, raw.location.address,
    // fall back to the normalized `location` string we already compute.
    if (!item) return;
    const raw = item.raw || {};
    const loc = raw.location || raw.Location || null;
    const addrCandidate = (loc && (loc.path || loc.address || loc.name)) || item.location || '';
    if (!addrCandidate) {
      showNotification('Локация не указана', 'error');
      return;
    }
    // Build the Yandex.Maps URL and URL-encode the query part.
    const q = encodeURIComponent(String(addrCandidate));
    const url = `https://yandex.ru/maps/?text=${q}`;
    // Open in a new tab safely
    const newWin = window.open(url, '_blank', 'noopener,noreferrer');
    if (newWin) newWin.opener = null;
  };

  const handleExtend = async (id: string) => {
    try {
      await bookService.extendReservation(id);
      showNotification('Бронь продлена на 1 день', 'success');
      // reload
      const data = await bookService.getReservedBooks();
      const formatLoc = (loc: any, bookLoc: any) => {
        const l = loc || bookLoc || null;
        if (!l) return '';
        return l.address || l.name || '';
      };

      const normalized = (data || []).map((ex: any) => {
        const book = ex.book || ex.Book || {};
        const expiresAt = ex.expires_at ? new Date(ex.expires_at) : null;
        let expiresLabel = 'истекло';
        if (expiresAt) {
          const diff = Math.max(0, Math.floor((expiresAt.getTime() - Date.now()) / 1000));
          const days = Math.floor(diff / 86400);
          const hours = Math.floor((diff % 86400) / 3600);
          if (diff <= 0) expiresLabel = 'истекла';
          else if (days > 0) expiresLabel = `через ${days}д ${hours}ч`;
          else expiresLabel = `через ${hours}ч`;
        }
        return {
          id: book.id,
          exchangeId: ex.id,
          title: book.title,
          location: formatLoc(ex.location, book.current_location),
          expiresLabel,
          raw: ex,
        };
      });
      setBooks(normalized || []);
        // reservation changed, notify listeners
        try { window.dispatchEvent(new Event('books:update')); } catch (e) { /* noop */ }
    } catch (e: any) {
      showNotification(e?.response?.data?.error || e.message || 'Ошибка при продлении', 'error');
    }
  };

  return (
    <>
      <Box>
        <Typography variant="h5" gutterBottom>Забронированные</Typography>
        <Box sx={{ mt: 2, display: 'grid', gridTemplateColumns: { xs: 'repeat(1, 1fr)' }, gap: 3 }}>
          {books.map((b: any) => (
            <Card key={b.exchangeId || b.id} sx={{ boxShadow: 3, borderRadius: '24px', p: '20px' }}>
              <CardContent sx={{ pb: 0, pt: 0 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start' }}>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>{b.title || '—'}</Typography>
                  <Typography variant="body2" color="text.secondary">{b.expiresLabel || 'истекло'}</Typography>
                </Box>
                <Typography variant="body1" color="text.primary" sx={{ mt: 1 }}>{b.location || b.address || 'Локация не указана'}</Typography>
              </CardContent>
              <CardActions sx={{ mt: '30px', px: 0, py: 0, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                {/* Delete stays on the left */}
                <Box>
                  <IconButton
                    onClick={() => handleDelete(b.id)}
                    aria-label="Удалить"
                    sx={{ border: '1px solid', borderColor: 'divider', borderRadius: '100px', width: 72, height: 34, p: 0, display: 'flex', justifyContent: 'center', alignItems: 'center' }}
                  >
                    <DeleteOutlineIcon sx={{ fontSize: 20 }} />
                  </IconButton>
                </Box>

                {/* All other actions moved to the right */}
                <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
                  <Button variant="outlined" startIcon={<MapIcon />} onClick={() => handleOnMap(b)} sx={{ height: 34, px: 2, borderColor: 'primary.main', borderRadius: '100px' }}>
                    На карте
                  </Button>

                  {/* Иконка для продления перенесена сюда */}
                  <Button variant="outlined" startIcon={<AutorenewIcon />} onClick={() => handleExtend(b.id)} sx={{ height: 34, px: 2, borderColor: 'primary.main', borderRadius: '100px' }}>
                    Продлить
                  </Button>

                  {/* Получить справа, без иконки */}
                  <Button variant="contained" color="primary" onClick={() => openQr(b)} sx={{ height: 34, px: 2, borderRadius: '100px' }}>
                    Получить
                  </Button>
                </Box>
              </CardActions>
            </Card>
          ))}
        </Box>

        <Typography variant="body1" color="text.primary" align="center" sx={{ mt: 2, fontWeight: 500, opacity: 0.6 }}>
          Книги, которые вы забронировали
        </Typography>
      </Box>

      <Dialog open={qrOpen} onClose={closeQr} maxWidth="xs" fullWidth>
        <DialogTitle>Показать код для получения</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', justifyContent: 'center', my: 2 }}>
            <Paper elevation={3} sx={{ width: 220, height: 220, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
              {/* Simple SVG placeholder QR with exchange id or book id text */}
              <Box sx={{ textAlign: 'center', px: 2 }}>
                <svg width="180" height="180" xmlns="http://www.w3.org/2000/svg">
                  <rect width="100%" height="100%" fill="#f5f5f5" />
                  <text x="50%" y="50%" dominantBaseline="middle" textAnchor="middle" fontSize="10" fill="#222">{selected?.exchangeId || selected?.id}</text>
                </svg>
              </Box>
            </Paper>
          </Box>
          <Typography variant="body2" color="text.secondary" align="center">Покажите этот код у точки выдачи. Нажмите «Получить», чтобы пометить книгу как взятую.</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={closeQr}>Отмена</Button>
          <Button variant="contained" color="primary" onClick={handleReceive}>Получить</Button>
        </DialogActions>
      </Dialog>
    </>
  );
};

export default ReservedBooksPage;