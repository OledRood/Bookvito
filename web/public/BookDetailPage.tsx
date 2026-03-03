import React, { useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Container,
  Paper,
  Button,
  Chip,
  CircularProgress,
  Alert,
  Divider,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
} from '@mui/material';
import FlagIcon from '@mui/icons-material/Flag';
import useMediaQuery from '@mui/material/useMediaQuery';
import { useTheme } from '@mui/material/styles';
import RoomIcon from '@mui/icons-material/Room';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import resolveImageUrl from '../src/utils/imageUrl';
import { useNotification } from './NotificationContext';
import { useAuth } from './AuthContext';
import useBook from '../src/hooks/useBook';
// useLocationsList removed: location is read-only from book
import useRequestBook from '../src/hooks/useRequestBook';
import moderService from '../src/services/moderService';
import ButtonSpinner from '../src/components/ButtonSpinner';

const BookDetailPage: React.FC = () => {
  const { bookId } = useParams<{ bookId: string }>();
  const navigate = useNavigate();
  const { showNotification } = useNotification();
  const { isAuthenticated } = useAuth();
  const { book, loading, refresh } = useBook(bookId || undefined);
  const { requestBook, loading: requesting } = useRequestBook();
  const [reportOpen, setReportOpen] = useState(false);
  const [reportReason, setReportReason] = useState('');
  const [reporting, setReporting] = useState(false);
  const theme = useTheme();
  const isMdUp = useMediaQuery(theme.breakpoints.up('md'));
  const titleRef = React.useRef<HTMLDivElement | null>(null);
  const blocksRef = React.useRef<HTMLDivElement | null>(null);
  const [heroMaxHeight, setHeroMaxHeight] = React.useState<number | undefined>(undefined);

  const handleReport = async () => {
    if (!reportReason.trim()) return;
    setReporting(true);
    try {
      await moderService.reportBook(bookId as string, reportReason.trim());
      showNotification('Жалоба отправлена модератору', 'success');
      setReportOpen(false);
      setReportReason('');
    } catch (e: any) {
      showNotification(e?.response?.data?.error || 'Ошибка при отправке жалобы', 'error');
    } finally {
      setReporting(false);
    }
  };

  const handleBooking = async () => {
    const locId = book?.location?.id;
    if (!locId) {
      showNotification('Локация для получения не задана.', 'warning');
      return;
    }
    try {
      // pass explicit locationId so backend can persist chosen pickup location
      await requestBook(bookId as string, locId);
      // refresh book details to reflect new status (reserved)
      try {
        await refresh();
      } catch (e) {
        // ignore refresh errors (notification already shown by useBook)
      }
      navigate('/books');
    } catch (err) {
      // requestBook already shows notification
    }
  };

  if (loading) {
    return (
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'center', p: 4 }}>
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (!book) {
    return (
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Alert severity="error" action={<Button onClick={() => navigate(-1)}>Вернуться</Button>}>
          Книга не найдена.
        </Alert>
      </Container>
    );
  }

  // Map backend condition values to user-friendly labels and statuses
  const conditionMap: Record<string, { label: string; key: string; color: 'error' | 'warning' | 'success' }> = {
    bad: { label: 'Потрепанное', key: 'bad', color: 'error' },
    worn: { label: 'Потрепанное', key: 'bad', color: 'error' },
    normal: { label: 'Нормальное', key: 'normal', color: 'warning' },
    ok: { label: 'Нормальное', key: 'normal', color: 'warning' },
    new: { label: 'Отличное', key: 'excellent', color: 'success' },
    excellent: { label: 'Отличное', key: 'excellent', color: 'success' },
  };

  const cond = book.condition ? (conditionMap[book.condition] || { label: book.condition, key: book.condition, color: 'warning' as const }) : null;

  // We avoid dynamic DOM measurements (they caused hook/effect instability).
  // Instead use CSS-based constraints for the hero height (see below).

  const openLocationInYandexMaps = (loc?: any) => {
    // prefer explicit location object passed, otherwise try book.location
    const l = loc || book?.location || null;
    if (!l) {
      showNotification('Локация не указана', 'warning');
      return;
    }
    const candidate = l.path || l.address || l.name || '';
    if (!candidate) {
      showNotification('Локация не указана', 'warning');
      return;
    }
    const q = encodeURIComponent(String(candidate));
    const url = `https://yandex.ru/maps/?text=${q}`;
    const newWin = window.open(url, '_blank', 'noopener,noreferrer');
    if (newWin) newWin.opener = null;
  };

  

      return (
        <Box>
          {isMdUp ? (
            // Desktop/tablet layout: image left, blocks right
            <Container maxWidth="lg" sx={{ mt: 4, mb: 6 }}>
              <Box sx={{ display: 'flex', gap: 4, alignItems: 'flex-start' }}>
                <Box sx={{ width: '40%', borderRadius: 2, overflow: 'hidden' }}>
                  <Box
                    component="img"
                    src={resolveImageUrl(book.imageUrl)}
                    alt={book.title}
                    sx={{ width: '100%', height: '100%', display: 'block', objectFit: 'cover' }}
                  />
                </Box>
                <Box sx={{ width: '60%', display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <Box>
                    <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>{book.title}</Typography>
                    <Box sx={{ width: 60, height: 4, bgcolor: 'divider', my: 1, borderRadius: 2 }} />
                    <Typography variant="subtitle1" color="text.secondary">{book.author}</Typography>
                  </Box>
                  <Paper sx={{ px: { xs: 2, md: 4 }, py: { xs: 0.5, md: 1 }, borderRadius: 0.75, width: '100%', maxWidth: '100%' }} elevation={1}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2 }}>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Состояние</Typography>
                      <Typography
                        variant="body2"
                        sx={(t) => ({ textTransform: 'lowercase', fontWeight: 600, color: cond ? t.palette[cond.color].main : t.palette.text.secondary })}
                      >
                        {cond ? cond.label.toLowerCase() : '—'}
                      </Typography>
                    </Box>
                  </Paper>
                  <Paper sx={{ px: { xs: 2, md: 4 }, py: { xs: 0.5, md: 1 }, borderRadius: 0.75 }} elevation={1}>
                    <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>Описание</Typography>
                    <Typography variant="body1">{book.description || 'Описание отсутствует.'}</Typography>
                  </Paper>
                  <Paper
                    onClick={() => openLocationInYandexMaps()}
                    sx={{
                      px: { xs: 2, md: 4 },
                      py: { xs: 0.5, md: 1 },
                      borderRadius: 0.75,
                      display: 'flex',
                      alignItems: 'center',
                      gap: 1,
                      cursor: 'pointer',
                      '&:hover': { bgcolor: 'action.hover' },
                    }}
                    elevation={1}
                  >
                    <RoomIcon color="action" />
                    <Box sx={{ flex: 1 }}>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Локация</Typography>
                      <Typography variant="body2" color="text.secondary">
                        {(() => {
                          const loc = book.location;
                          if (!loc) return 'Не указана';
                          if (loc.address && loc.name) return `${loc.address} (${loc.name})`;
                          return loc.address || loc.name || 'Не указана';
                        })()}
                      </Typography>
                    </Box>
                    <ChevronRightIcon color="action" />
                  </Paper>

                  <Button
                    variant="contained"
                    color="primary"
                    size="large"
                    sx={{ mt: 1, borderRadius: 3, py: 1.5, alignSelf: 'flex-start', width: 220 }}
                    onClick={handleBooking}
                    disabled={requesting || !book.location?.id || !isAuthenticated}
                  >
                    {requesting ? <ButtonSpinner /> : 'Забронировать'}
                  </Button>

                  {!isAuthenticated && (
                    <Box sx={{ display: 'flex', gap: 1, mt: 1, alignItems: 'center' }}>
                      <Typography variant="body2" color="text.secondary">Чтобы забронировать книгу, пожалуйста, войдите.</Typography>
                      <Button variant="text" size="small" onClick={() => navigate('/login')}>Войти</Button>
                    </Box>
                  )}
                </Box>
              </Box>
            </Container>
          ) : (
            // Mobile layout: hero + stacked blocks
            <>
              {/* Hero with background image */}
              <Box
                sx={{
                  // Default hero height on mobile; we'll cap it dynamically based on
                  // title+blocks height so the image never exceeds the visible content.
                  height: { xs: 440, md: 520 },
                  maxHeight: heroMaxHeight ? `${heroMaxHeight}px` : { xs: 'min(440px, calc(100vh - 220px))', md: 520 },
                  overflow: 'hidden',
                  // Apply a gradient that is transparent at the top and strongly darkens toward the bottom
                  // so the text on the image is highly legible on mobile.
                  backgroundImage: `linear-gradient(to bottom, rgba(0,0,0,0) 55%, rgba(0,0,0,0.8) 100%), url(${resolveImageUrl(book.imageUrl)})`,
                  backgroundSize: 'cover',
                  backgroundPosition: 'center',
                  display: 'flex',
                  alignItems: 'flex-end',
                  color: '#fff',
                }}
              >
                <Container ref={(el) => { titleRef.current = el; }} maxWidth="md" sx={{ pb: 3 }}>
                  <Typography
                    variant="h3"
                    component="h1"
                    sx={{
                      fontWeight: 700,
                      lineHeight: 1.05,
                      fontSize: { xs: '2.2rem', md: '3rem' },
                      textShadow: '0 6px 18px rgba(0,0,0,0.65)'
                    }}
                  >
                    {book.title}
                  </Typography>
                  <Box sx={{ width: 60, height: 4, bgcolor: 'rgba(255,255,255,0.85)', my: 1, borderRadius: 2 }} />
                  <Typography variant="h6" sx={{ color: 'rgba(255,255,255,0.95)', textShadow: '0 4px 12px rgba(0,0,0,0.55)' }}>
                    {book.author}
                  </Typography>
                </Container>
              </Box>

              {/* Content blocks */}
              {/*
                Previously we used a negative top margin to pull the content up over the hero.
                That caused overlap with the hero on narrow screens and hid the author line.
                Use a positive top margin so blocks appear below the hero instead.
              */}
              <Container ref={(el) => { blocksRef.current = el; }} maxWidth="md" sx={{ mt: 3, mb: 6 }}>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <Paper sx={{ px: { xs: 2, md: 4 }, py: { xs: 0.5, md: 1 }, borderRadius: 0.75, width: '100%', maxWidth: '100%' }} elevation={1}>
                    <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2 }}>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Состояние</Typography>
                      <Typography
                        variant="body2"
                        sx={(theme) => ({
                          textTransform: 'lowercase',
                          fontWeight: 600,
                          color: cond ? theme.palette[cond.color].main : theme.palette.text.secondary,
                        })}
                      >
                        {cond ? cond.label.toLowerCase() : '—'}
                      </Typography>
                    </Box>
                  </Paper>

                  <Paper sx={{ px: { xs: 2, md: 4 }, py: { xs: 0.5, md: 1 }, borderRadius: 0.75 }} elevation={1}>
                    <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>Описание</Typography>
                    <Typography variant="body1">{book.description || 'Описание отсутствует.'}</Typography>
                    <Box sx={{ display: 'flex', justifyContent: 'flex-end', mt: 1 }}>
                      {/* optional expand link */}
                    </Box>
                  </Paper>

                  <Paper
                    onClick={() => openLocationInYandexMaps()}
                    sx={{
                      px: { xs: 2, md: 4 },
                      py: { xs: 0.5, md: 1 },
                      borderRadius: 0.75,
                      display: 'flex',
                      alignItems: 'center',
                      gap: 1,
                      cursor: 'pointer',
                      '&:hover': { bgcolor: 'action.hover' },
                    }}
                    elevation={1}
                  >
                    <RoomIcon color="action" />
                    <Box sx={{ flex: 1 }}>
                      <Typography variant="subtitle1" sx={{ fontWeight: 600 }}>Локация</Typography>
                      <Typography variant="body2" color="text.secondary">
                        {(() => {
                          const loc = book.location;
                          if (!loc) return 'Не указана';
                          if (loc.address && loc.name) return `${loc.address} (${loc.name})`;
                          return loc.address || loc.name || 'Не указана';
                        })()}
                      </Typography>
                    </Box>
                    <ChevronRightIcon color="action" />
                  </Paper>

                  <Button
                    variant="contained"
                    color="primary"
                    size="large"
                    fullWidth
                    sx={{ mt: 1, borderRadius: 3, py: 2 }}
                    onClick={handleBooking}
                    disabled={requesting || !book.location?.id || !isAuthenticated}
                  >
                    {requesting ? <ButtonSpinner /> : 'Забронировать'}
                  </Button>

                  {!isAuthenticated && (
                    <Box sx={{ display: 'flex', gap: 1, mt: 1, alignItems: 'center' }}>
                      <Typography variant="body2" color="text.secondary">Чтобы забронировать книгу, пожалуйста, войдите.</Typography>
                      <Button variant="text" size="small" onClick={() => navigate('/login')}>Войти</Button>
                    </Box>
                  )}

                  {isAuthenticated && (
                    <Button
                      variant="text"
                      color="error"
                      size="small"
                      startIcon={<FlagIcon />}
                      sx={{ mt: 0.5 }}
                      onClick={() => setReportOpen(true)}
                    >
                      Пожаловаться
                    </Button>
                  )}

                  {/* Диалог жалобы */}
                  <Dialog open={reportOpen} onClose={() => setReportOpen(false)} maxWidth="sm" fullWidth>
                    <DialogTitle>Пожаловаться на объявление</DialogTitle>
                    <DialogContent>
                      <TextField
                        autoFocus
                        multiline
                        rows={3}
                        fullWidth
                        label="Причина жалобы"
                        value={reportReason}
                        onChange={e => setReportReason(e.target.value)}
                        sx={{ mt: 1 }}
                        placeholder="Опишите, что не так с этим объявлением..."
                      />
                    </DialogContent>
                    <DialogActions>
                      <Button onClick={() => setReportOpen(false)}>Отмена</Button>
                      <Button
                        variant="contained"
                        color="error"
                        onClick={handleReport}
                        disabled={!reportReason.trim() || reporting}
                      >
                        {reporting ? <ButtonSpinner /> : 'Отправить'}
                      </Button>
                    </DialogActions>
                  </Dialog>
                </Box>
              </Container>
            </>
          )}
    </Box>
  );
};

export default BookDetailPage;