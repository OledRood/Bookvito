import React, { useEffect, useMemo, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
  Box,
  Button,
  CircularProgress,
  Container,
  FormControl,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  SelectChangeEvent,
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material';
import PhotoCamera from '@mui/icons-material/PhotoCamera';
import DeleteOutline from '@mui/icons-material/DeleteOutline';
import ButtonSpinner from '../src/components/ButtonSpinner';
import useBook from '../src/hooks/useBook';
import useReturnBook from '../src/hooks/useReturnBook';
import { useLocationsList } from '../src/hooks/useLocationsList';
import { deleteBookImage, setBookImage, updateBook, uploadBookImage } from '../src/services/bookService';
import resolveImageUrl from '../src/utils/imageUrl';
import { useNotification } from '../contexts/NotificationContext';
import { useAuth } from './AuthContext';

type ConditionValue = 'excellent' | 'good' | 'bad';

const maxImageSize = 20 * 1024 * 1024;
const allowedImageExtensions = new Set(['jpg', 'jpeg', 'png', 'gif', 'webp']);

const conditionMap: Record<string, ConditionValue> = {
  excellent: 'excellent',
  new: 'excellent',
  ok: 'good',
  normal: 'good',
  good: 'good',
  worn: 'bad',
  bad: 'bad',
};

const ReturnBookPage: React.FC = () => {
  const { bookId } = useParams<{ bookId: string }>();
  const navigate = useNavigate();
  const fileInputRef = useRef<HTMLInputElement | null>(null);
  const { showNotification } = useNotification();
  const { isAuthenticated } = useAuth();
  const { locations } = useLocationsList();
  const { book, loading, refresh } = useBook(bookId || undefined);
  const { returnBook, loading: returning } = useReturnBook();

  const [title, setTitle] = useState('');
  const [author, setAuthor] = useState('');
  const [description, setDescription] = useState('');
  const [condition, setCondition] = useState<ConditionValue>('good');
  const [locationId, setLocationId] = useState('');
  const [currentImageUrl, setCurrentImageUrl] = useState('');
  const [removeCurrentImage, setRemoveCurrentImage] = useState(false);
  const [pendingFile, setPendingFile] = useState<File | null>(null);
  const [pendingPreviewUrl, setPendingPreviewUrl] = useState('');
  const [actionInProgress, setActionInProgress] = useState<'save' | 'return' | null>(null);

  useEffect(() => {
    if (!book) {
      return;
    }
    setTitle(book.title || '');
    setAuthor(book.author || '');
    setDescription(book.description || '');
    setCondition(conditionMap[book.condition || ''] || 'good');
    setLocationId(book.location?.id || '');
    setCurrentImageUrl(book.imageUrl || '');
    setRemoveCurrentImage(false);
    setPendingFile(null);
    setPendingPreviewUrl('');
  }, [book]);

  useEffect(() => {
    if (!pendingFile) {
      setPendingPreviewUrl('');
      return;
    }
    const nextUrl = URL.createObjectURL(pendingFile);
    setPendingPreviewUrl(nextUrl);
    return () => {
      URL.revokeObjectURL(nextUrl);
    };
  }, [pendingFile]);

  const displayImageSrc = useMemo(() => {
    if (pendingPreviewUrl) {
      return pendingPreviewUrl;
    }
    if (!removeCurrentImage && currentImageUrl) {
      return resolveImageUrl(currentImageUrl);
    }
    return '';
  }, [currentImageUrl, pendingPreviewUrl, removeCurrentImage]);

  const isBusy = actionInProgress !== null || returning;

  const getErrorMessage = (error: any, fallback: string) => (
    error?.response?.data?.message ||
    error?.response?.data?.error ||
    error?.message ||
    fallback
  );

  const handleLocationChange = (event: SelectChangeEvent<string>) => {
    setLocationId(event.target.value);
  };

  const validateForm = () => {
    if (!isAuthenticated) {
      showNotification('Вам нужно войти, чтобы изменить книгу', 'warning');
      return false;
    }
    if (!title.trim()) {
      showNotification('Название книги не может быть пустым', 'warning');
      return false;
    }
    if (!author.trim()) {
      showNotification('Автор книги не может быть пустым', 'warning');
      return false;
    }
    return true;
  };

  const handleSelectImage = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) {
      return;
    }

    const ext = file.name.split('.').pop()?.toLowerCase() || '';
    if (!allowedImageExtensions.has(ext)) {
      showNotification('Поддерживаются только JPG, JPEG, PNG, GIF и WEBP', 'error');
      event.target.value = '';
      return;
    }
    if (file.size <= 0) {
      showNotification('Файл изображения пустой', 'error');
      event.target.value = '';
      return;
    }
    if (file.size > maxImageSize) {
      showNotification('Размер изображения не может превышать 20 MB', 'error');
      event.target.value = '';
      return;
    }

    setPendingFile(file);
    setRemoveCurrentImage(false);
    event.target.value = '';
  };

  const handleRemoveImage = () => {
    if (pendingFile) {
      setPendingFile(null);
      return;
    }
    if (currentImageUrl) {
      setRemoveCurrentImage(true);
    }
  };

  const syncImage = async () => {
    if (!book) {
      throw new Error('Книга не найдена');
    }

    let nextImageUrl = removeCurrentImage ? '' : currentImageUrl;

    if (removeCurrentImage && currentImageUrl) {
      await deleteBookImage(book.id);
      setCurrentImageUrl('');
      nextImageUrl = '';
      setRemoveCurrentImage(false);
    }

    if (pendingFile) {
      const uploadResponse = await uploadBookImage(pendingFile);
      const uploadedUrl = uploadResponse?.url || '';
      if (!uploadedUrl) {
        throw new Error('Не удалось получить URL загруженного изображения');
      }
      await setBookImage(book.id, uploadedUrl);
      setCurrentImageUrl(uploadedUrl);
      setPendingFile(null);
      setRemoveCurrentImage(false);
      nextImageUrl = uploadedUrl;
    }

    return nextImageUrl;
  };

  const buildUpdatePayload = (imageUrl: string) => ({
    title: title.trim(),
    author: author.trim(),
    description,
    condition,
    image_url: imageUrl,
    current_location_id: locationId || null,
  });

  const handleSave = async () => {
    if (!book || !validateForm()) {
      return;
    }

    setActionInProgress('save');
    try {
      const finalImageUrl = await syncImage();
      await updateBook(book.id, buildUpdatePayload(finalImageUrl));
      await refresh();
      try { window.dispatchEvent(new Event('books:update')); } catch (error) {}
      showNotification('Изменения книги сохранены', 'success');
    } catch (error: any) {
      showNotification(getErrorMessage(error, 'Не удалось сохранить изменения'), 'error');
    } finally {
      setActionInProgress(null);
    }
  };

  const handleReturn = async () => {
    if (!book || !validateForm()) {
      return;
    }

    setActionInProgress('return');
    try {
      const finalImageUrl = await syncImage();
      await returnBook({
        bookId: book.id,
        title: title.trim(),
        author: author.trim(),
        description,
        condition,
        currentLocationId: locationId || undefined,
        imageUrl: finalImageUrl,
      }, { silent: true });
      showNotification('Книга успешно возвращена', 'success');
      try { window.dispatchEvent(new Event('books:update')); } catch (error) {}
      navigate('/books/shelf');
    } catch (error: any) {
      showNotification(getErrorMessage(error, 'Не удалось вернуть книгу'), 'error');
    } finally {
      setActionInProgress(null);
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
        <Typography variant="h6">Книга не найдена</Typography>
        <Button onClick={() => navigate(-1)} sx={{ mt: 2 }}>Назад</Button>
      </Container>
    );
  }

  return (
    <Container maxWidth="lg" sx={{ py: { xs: 3, md: 4 } }}>
      <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 3, alignItems: 'flex-start' }}>
        <Paper
          variant="outlined"
          sx={{
            width: { xs: '100%', md: 380 },
            p: 2,
            borderRadius: 3,
            display: 'flex',
            flexDirection: 'column',
            gap: 2,
            position: { md: 'sticky' },
            top: { md: 24 },
          }}
        >
          <Box
            sx={{
              aspectRatio: '3 / 4',
              borderRadius: 2,
              overflow: 'hidden',
              bgcolor: 'action.hover',
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            {displayImageSrc ? (
              <Box
                component="img"
                src={displayImageSrc}
                alt={title || 'Обложка книги'}
                sx={{ width: '100%', height: '100%', objectFit: 'cover', display: 'block' }}
              />
            ) : (
              <Typography variant="body2" color="text.secondary">
                Фото не загружено
              </Typography>
            )}
          </Box>

          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            hidden
            onChange={handleSelectImage}
          />

          <Stack direction={{ xs: 'column', sm: 'row', md: 'column' }} spacing={1.5}>
            <Button
              variant="outlined"
              startIcon={<PhotoCamera />}
              onClick={() => fileInputRef.current?.click()}
              disabled={isBusy}
            >
              {currentImageUrl || pendingFile ? 'Заменить фото' : 'Загрузить фото'}
            </Button>
            <Button
              variant="text"
              color="inherit"
              startIcon={<DeleteOutline />}
              onClick={handleRemoveImage}
              disabled={isBusy || (!currentImageUrl && !pendingFile) || removeCurrentImage}
            >
              Удалить фото
            </Button>
          </Stack>

          <Typography variant="body2" color="text.secondary">
            Фото применится вместе с действием «Сохранить» или «Вернуть».
          </Typography>
        </Paper>

        <Paper
          variant="outlined"
          sx={{
            flex: 1,
            width: '100%',
            p: { xs: 2, md: 3 },
            borderRadius: 3,
          }}
        >
          <Stack spacing={2.5}>
            <Box>
              <Typography variant="h4" component="h1" sx={{ fontWeight: 700 }}>
                Изменение книги
              </Typography>
            </Box>

            <TextField
              label="Название"
              value={title}
              onChange={(event) => setTitle(event.target.value)}
              fullWidth
              required
              inputProps={{ maxLength: 255 }}
            />

            <TextField
              label="Автор"
              value={author}
              onChange={(event) => setAuthor(event.target.value)}
              fullWidth
              required
              inputProps={{ maxLength: 255 }}
            />

            <Box>
              <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 600 }}>
                Состояние
              </Typography>
              <ToggleButtonGroup
                exclusive
                value={condition}
                onChange={(_, value) => {
                  if (value !== null) {
                    setCondition(value);
                  }
                }}
                fullWidth
                sx={{ display: 'grid', gridTemplateColumns: 'repeat(3, minmax(0, 1fr))' }}
              >
                <ToggleButton value="bad" sx={{ textTransform: 'none' }}>Потрепанное</ToggleButton>
                <ToggleButton value="good" sx={{ textTransform: 'none' }}>Нормальное</ToggleButton>
                <ToggleButton value="excellent" sx={{ textTransform: 'none' }}>Отличное</ToggleButton>
              </ToggleButtonGroup>
            </Box>

            <TextField
              label="Описание"
              value={description}
              onChange={(event) => setDescription(event.target.value)}
              multiline
              minRows={4}
              fullWidth
              inputProps={{ maxLength: 2000 }}
            />

            <FormControl fullWidth>
              <InputLabel id="return-book-location-label">Пункт приёма</InputLabel>
              <Select
                labelId="return-book-location-label"
                label="Пункт приёма"
                value={locationId}
                onChange={handleLocationChange}
              >
                <MenuItem value="">
                  <em>Не выбрано</em>
                </MenuItem>
                {(locations || []).map((location) => (
                  <MenuItem key={location.id} value={location.id}>
                    {location.address || location.name}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>

            <Stack direction={{ xs: 'column', sm: 'row' }} spacing={1.5} sx={{ pt: 1 }}>
              <Button
                variant="outlined"
                size="large"
                onClick={handleSave}
                disabled={isBusy || !isAuthenticated}
                sx={{ minWidth: 180 }}
              >
                {actionInProgress === 'save' ? <ButtonSpinner /> : 'Сохранить'}
              </Button>
              <Button
                variant="contained"
                size="large"
                onClick={handleReturn}
                disabled={isBusy || !isAuthenticated}
                sx={{ minWidth: 180 }}
              >
                {actionInProgress === 'return' || returning ? <ButtonSpinner /> : 'Вернуть'}
              </Button>
            </Stack>

            {!isAuthenticated && (
              <Typography variant="body2" color="text.secondary">
                Чтобы сохранить изменения или вернуть книгу, нужно войти в аккаунт.
              </Typography>
            )}
          </Stack>
        </Paper>
      </Box>
    </Container>
  );
};

export default ReturnBookPage;
