import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  Box,
  Typography,
  Container,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  ToggleButtonGroup,
  ToggleButton,
  CircularProgress,
  SelectChangeEvent,
} from '@mui/material';
import ButtonSpinner from '../src/components/ButtonSpinner';
import { useTheme } from '@mui/material/styles';
import useMediaQuery from '@mui/material/useMediaQuery';
import resolveImageUrl from '../src/utils/imageUrl';
import useBook from '../src/hooks/useBook';
import { useNotification } from '../contexts/NotificationContext';
import { useAuth } from './AuthContext';
import { useLocationsList } from '../src/hooks/useLocationsList';
import useReturnBook from '../src/hooks/useReturnBook';

const ReturnBookPage: React.FC = () => {
  const { bookId } = useParams<{ bookId: string }>();
  const navigate = useNavigate();
  const theme = useTheme();
  const isMdUp = useMediaQuery(theme.breakpoints.up('md'));
  const { showNotification } = useNotification();
  const { isAuthenticated } = useAuth();
  const { locations } = useLocationsList();
  const { book, loading } = useBook(bookId || undefined);
  const { returnBook, loading: returning } = useReturnBook();

  const [condition, setCondition] = useState<'excellent' | 'good' | 'bad'>('good');
  const [locationId, setLocationId] = useState<string | undefined>(undefined);
  const [description, setDescription] = useState<string>('');

  useEffect(() => {
    if (book) {
      const map: any = { excellent: 'excellent', new: 'excellent', ok: 'good', normal: 'good', worn: 'bad', bad: 'bad' };
      setCondition((book.condition && (map[book.condition] || book.condition)) || 'good');
      setLocationId(book.location?.id || undefined);
      setDescription(book.description || '');
    }
  }, [book]);

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
        <Typography variant="h6">Книга не найдена</Typography>
        <Button onClick={() => navigate(-1)} sx={{ mt: 2 }}>Назад</Button>
      </Container>
    );
  }

  const handleLocationChange = (e: SelectChangeEvent<string>) => setLocationId(e.target.value as string);

  const handleSubmit = async () => {
    if (!isAuthenticated) {
      showNotification('Вам нужно войти, чтобы вернуть книгу', 'warning');
      return;
    }
    try {
      await returnBook({
        bookId: book.id,
        title: book.title,
        author: book.author,
        condition,
        description,
        currentLocationId: locationId,
      });
      try { window.dispatchEvent(new Event('books:update')); } catch (e) {}
      showNotification('Книга успешно возвращена', 'success');
      navigate('/books/shelf');
    } catch (err) {
      // handled in hook
    }
  };

  return (
    <Box>
      {isMdUp ? (
        // Desktop/tablet layout: image left, form right
        <Container maxWidth="lg" sx={{ mt: 4, mb: 6 }}>
          <Box sx={{ display: 'flex', gap: 4, alignItems: 'flex-start' }}>
            <Box sx={{ width: '40%', borderRadius: 2, overflow: 'hidden' }}>
              <Box
                component="img"
                src={resolveImageUrl((book as any).imageUrl || (book as any).image_url || (book as any).image || '')}
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

              <Box sx={{ px: { xs: 2, md: 4 }, py: { xs: 1, md: 2 }, width: '100%' }}>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>

                  <Box>
                    <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>Состояние</Typography>
                    <FormControl fullWidth>
                      <ToggleButtonGroup
                        exclusive
                        value={condition}
                        onChange={(_, val) => { if (val !== null) setCondition(val as any); }}
                        sx={(theme: any) => ({
                          display: 'flex',
                          width: '100%',
                          borderRadius: 12,
                          overflow: 'hidden',
                          border: 'none',
                          boxShadow: 'none',
                          bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`,
                        })}
                        aria-label="Состояние книги"
                      >
                        <ToggleButton value="bad" aria-label="Потрепанное" sx={(theme: any) => ({ flex: 1, borderRadius: 0, textTransform: 'none', fontWeight: 500, border: 'none', boxShadow: 'none', outline: 'none', bgcolor: `var(--md-sys-color-primary-container, transparent)`, color: theme.palette.text.primary, '&.Mui-selected': { bgcolor: `var(--md-sys-color-primary, ${theme.palette.primary.main})`, color: `var(--md-sys-color-on-primary, ${theme.palette.getContrastText(theme.palette.primary.main)})`, borderRadius: 12, }, '&:not(.Mui-selected):hover': { bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`, }, })}>
                          Потрепанное
                        </ToggleButton>
                        <ToggleButton value="good" aria-label="Нормальное" sx={(theme: any) => ({ flex: 1, borderRadius: 0, textTransform: 'none', fontWeight: 500, border: 'none', boxShadow: 'none', outline: 'none', bgcolor: `var(--md-sys-color-primary-container, transparent)`, color: theme.palette.text.primary, '&.Mui-selected': { bgcolor: `var(--md-sys-color-primary, ${theme.palette.primary.main})`, color: `var(--md-sys-color-on-primary, ${theme.palette.getContrastText(theme.palette.primary.main)})`, borderRadius: 12, }, '&:not(.Mui-selected):hover': { bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`, }, })}>
                          Нормальное
                        </ToggleButton>
                        <ToggleButton value="excellent" aria-label="Отличное" sx={(theme: any) => ({ flex: 1, borderRadius: 0, textTransform: 'none', fontWeight: 500, border: 'none', boxShadow: 'none', outline: 'none', bgcolor: `var(--md-sys-color-primary-container, transparent)`, color: theme.palette.text.primary, '&.Mui-selected': { bgcolor: `var(--md-sys-color-primary, ${theme.palette.primary.main})`, color: `var(--md-sys-color-on-primary, ${theme.palette.getContrastText(theme.palette.primary.main)})`, borderRadius: 12, }, '&:not(.Mui-selected):hover': { bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`, }, })}>
                          Отличное
                        </ToggleButton>
                      </ToggleButtonGroup>
                    </FormControl>
                  </Box>

                  <TextField label="Описание" multiline rows={4} value={description} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDescription(e.target.value)} fullWidth sx={{ mt: 2, mb: 2 }} />

                  <FormControl fullWidth>
                    <InputLabel id="location-select-label">Пункт приёма</InputLabel>
                    <Select
                      labelId="location-select-label"
                      value={locationId || ''}
                      label="Пункт приёма"
                      onChange={handleLocationChange}
                      renderValue={(selected) => {
                        if (!selected) return '';
                        const loc = (locations || []).find((l) => l.id === selected as string);
                        return loc ? loc.name : '';
                      }}
                    >
                      {(locations || []).map((loc) => (
                        <MenuItem key={loc.id} value={loc.id} sx={{ whiteSpace: 'normal' }}>{loc.address || loc.name}</MenuItem>
                      ))}
                    </Select>
                  </FormControl>

                  <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
                    <Button
                        variant="contained"
                        color="primary"
                        size="large"
                        sx={{ borderRadius: 3, py: 1.5, width: 220 }}
                        onClick={handleSubmit}
                        disabled={returning || !isAuthenticated}
                      >
                        {returning ? <ButtonSpinner /> : 'Вернуть'}
                      </Button>
                  </Box>
                </Box>
              </Box>
            </Box>
          </Box>
        </Container>
      ) : (
        // Mobile layout: hero + stacked blocks
        <>
          <Box sx={{ height: 360, backgroundImage: `linear-gradient(to bottom, rgba(0,0,0,0) 55%, rgba(0,0,0,0.8) 100%), url(${resolveImageUrl((book as any).imageUrl || (book as any).image_url || (book as any).image || '')})`, backgroundSize: 'cover', backgroundPosition: 'center', display: 'flex', alignItems: 'flex-end', color: '#fff', borderRadius: 2, overflow: 'hidden' }}>
            <Container sx={{ pb: 3 }}>
              <Typography variant="h4" sx={{ fontWeight: 700 }}>{book.title}</Typography>
              <Typography variant="subtitle1">{book.author}</Typography>
            </Container>
          </Box>

          <Container sx={{ mt: 3 }}>
            <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, mb: 6 }}>
              <Box sx={{ px: 2, py: 2 }}>
                <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2 }}>
                  <Typography variant="h6" sx={{ fontWeight: 700 }}>{book.title}</Typography>
                  <Typography variant="subtitle2" color="text.secondary">{book.author}</Typography>

                  <Box sx={{ mb: 2 }}>
                    <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>Состояние</Typography>
                    <FormControl fullWidth>
                      <ToggleButtonGroup
                        exclusive
                        value={condition}
                        onChange={(_, val) => { if (val !== null) setCondition(val as any); }}
                        sx={(theme: any) => ({
                          display: 'flex',
                          width: '100%',
                          borderRadius: 12,
                          overflow: 'hidden',
                          border: 'none',
                          boxShadow: 'none',
                          bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`,
                        })}
                        aria-label="Состояние книги"
                      >
                        <ToggleButton value="bad" aria-label="Потрепанное" sx={(theme: any) => ({ flex: 1, borderRadius: 0, textTransform: 'none', fontWeight: 500, border: 'none', boxShadow: 'none', outline: 'none', bgcolor: `var(--md-sys-color-primary-container, transparent)`, color: theme.palette.text.primary, '&.Mui-selected': { bgcolor: `var(--md-sys-color-primary, ${theme.palette.primary.main})`, color: `var(--md-sys-color-on-primary, ${theme.palette.getContrastText(theme.palette.primary.main)})`, borderRadius: 12, }, '&:not(.Mui-selected):hover': { bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`, }, })}>Потрепанное</ToggleButton>
                        <ToggleButton value="good" aria-label="Нормальное" sx={(theme: any) => ({ flex: 1, borderRadius: 0, textTransform: 'none', fontWeight: 500, border: 'none', boxShadow: 'none', outline: 'none', bgcolor: `var(--md-sys-color-primary-container, transparent)`, color: theme.palette.text.primary, '&.Mui-selected': { bgcolor: `var(--md-sys-color-primary, ${theme.palette.primary.main})`, color: `var(--md-sys-color-on-primary, ${theme.palette.getContrastText(theme.palette.primary.main)})`, borderRadius: 12, }, '&:not(.Mui-selected):hover': { bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`, }, })}>Нормальное</ToggleButton>
                        <ToggleButton value="excellent" aria-label="Отличное" sx={(theme: any) => ({ flex: 1, borderRadius: 0, textTransform: 'none', fontWeight: 500, border: 'none', boxShadow: 'none', outline: 'none', bgcolor: `var(--md-sys-color-primary-container, transparent)`, color: theme.palette.text.primary, '&.Mui-selected': { bgcolor: `var(--md-sys-color-primary, ${theme.palette.primary.main})`, color: `var(--md-sys-color-on-primary, ${theme.palette.getContrastText(theme.palette.primary.main)})`, borderRadius: 12, }, '&:not(.Mui-selected):hover': { bgcolor: `var(--md-sys-color-primary-container, ${theme.palette.action.hover})`, }, })}>Отличное</ToggleButton>
                      </ToggleButtonGroup>
                    </FormControl>
                  </Box>

                  <TextField label="Описание" multiline rows={4} value={description} onChange={(e: React.ChangeEvent<HTMLInputElement>) => setDescription(e.target.value)} fullWidth sx={{ mt: 2, mb: 2 }} />

                  <FormControl fullWidth>
                    <InputLabel id="location-select-label-mobile">Пункт приёма</InputLabel>
                    <Select
                      labelId="location-select-label-mobile"
                      value={locationId || ''}
                      label="Пункт приёма"
                      onChange={handleLocationChange}
                      renderValue={(selected) => {
                        if (!selected) return '';
                        const loc = (locations || []).find((l) => l.id === selected as string);
                        return loc ? loc.name : '';
                      }}
                    >
                      {(locations || []).map((loc) => (
                        <MenuItem key={loc.id} value={loc.id} sx={{ whiteSpace: 'normal' }}>{loc.address || loc.name}</MenuItem>
                      ))}
                    </Select>
                  </FormControl>

                  <Box sx={{ display: 'flex', gap: 2, mt: 2 }}>
                    <Button variant="contained" color="primary" size="large" fullWidth sx={{ borderRadius: 3, py: 2 }} onClick={handleSubmit} disabled={returning || !isAuthenticated}>
                      {returning ? <ButtonSpinner /> : 'Вернуть'}
                    </Button>
                  </Box>

                  {!isAuthenticated && (
                    <Box sx={{ display: 'flex', gap: 1, mt: 1, alignItems: 'center' }}>
                      <Typography variant="body2" color="text.secondary">Чтобы вернуть книгу, пожалуйста, войдите.</Typography>
                      <Button variant="text" size="small" onClick={() => navigate('/login')}>Войти</Button>
                    </Box>
                  )}
                </Box>
              </Box>
            </Box>
          </Container>
        </>
      )}
    </Box>
  );
};

export default ReturnBookPage;