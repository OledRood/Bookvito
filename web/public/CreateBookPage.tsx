import React, { useState, useEffect } from 'react';
import { useTheme } from '@mui/material/styles';
import useMediaQuery from '@mui/material/useMediaQuery';
import {
  Box,
  Typography,
  TextField,
  Button,
  Container,
  Select,
  MenuItem,
  FormControl,
  InputLabel,
  FormHelperText,
  CircularProgress,
  ToggleButton,
  ToggleButtonGroup,
} from '@mui/material';
import { useForm, Controller, SubmitHandler } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useNavigate } from 'react-router-dom';
import { useNotification } from '../contexts/NotificationContext';
import useCreateBook from '../src/hooks/useCreateBook';
import useLocationsList from '../src/hooks/useLocationsList';
import api from '../src/services/api';
import PhotoCamera from '@mui/icons-material/PhotoCamera';

interface CreateBookFormInputs {
  title: string;
  author: string;
  description: string;
  condition: 'excellent' | 'good' | 'bad';
  locationId: string;
  image: FileList | undefined;
}

const schema = yup.object().shape({
  title: yup.string().required('Название обязательно'),
  author: yup.string().required('Автор обязателен'),
  description: yup.string().max(1000, 'Описание не должно превышать 1000 символов'),
  condition: yup.string().oneOf(['excellent', 'good', 'bad']).required(),
  locationId: yup.string().required('Необходимо выбрать пункт выдачи').nonNullable(),
  image: yup
    .mixed()
    .test('required-image', 'Фото обязательно', (value) => {
      if (!value) return false;
      try {
        return (value as FileList).length > 0;
      } catch (e) {
        return false;
      }
    }),
});

const CreateBookPage: React.FC = () => {
  const navigate = useNavigate();
  const { showNotification } = useNotification();
  const { createBook, loading } = useCreateBook();
  const { locations, loading: loadingLocations } = useLocationsList();

  const {
    control,
    register,
    handleSubmit,
    formState: { errors, isSubmitted },
    setValue,
  } = useForm<any>({
    resolver: yupResolver(schema),
    defaultValues: {
      condition: 'good',
      locationId: '',
    },
  });

  const [file, setFile] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  useEffect(() => {
    if (file) {
      const dt = new DataTransfer();
      dt.items.add(file);
      setValue('image', dt.files, { shouldValidate: true });
    } else {
      setValue('image', undefined, { shouldValidate: true });
    }
  }, [file, setValue]);

  useEffect(() => {
    if (!file) {
      setPreviewUrl(null);
      return;
    }
    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
    return () => {
      URL.revokeObjectURL(url);
      setPreviewUrl((current) => (current === url ? null : current));
    };
  }, [file]);

  const fileInputRef = React.useRef<HTMLInputElement | null>(null);
  const theme = useTheme();
  const isMdUp = useMediaQuery(theme.breakpoints.up('md'));
  const showImageError = !!errors.image && !file && !!isSubmitted;

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    const newFiles = Array.from(e.dataTransfer.files || []);
    if (newFiles.length > 0) {
      setFile(newFiles[0]);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
  };

  const onInputFiles = (fileList: FileList | null) => {
    if (!fileList) return;
    const incoming = Array.from(fileList);
    if (incoming.length === 0) return;
    setFile(incoming[0]);
  };

  const removeFile = () => setFile(null);

  const onSubmit: SubmitHandler<CreateBookFormInputs> = async (data) => {
    try {
      const payload: any = {
        title: data.title,
        author: data.author,
        description: data.description,
        condition: data.condition,
        image_url: '',
        current_location_id: data.locationId || null,
      };
      if (data.image && data.image.length > 0) {
        const form = new FormData();
        form.append('image', data.image[0]);
        try {
          const resp = await api.post('books/upload', form, {
            headers: { 'Content-Type': 'multipart/form-data' },
          });
          payload.image_url = resp.data?.url || '';
        } catch (e) {
          showNotification('Не удалось загрузить изображение', 'error');
          return;
        }
      }
      await createBook(payload);
      navigate('/books');
    } catch (error) {
      // createBook already shows notification
    }
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ p: 0, mt: 4 }}>
        <label htmlFor="create-book-file-input" style={{ display: 'none' }}>
          <input
            id="create-book-file-input"
            ref={fileInputRef}
            type="file"
            accept="image/*"
            onChange={(e) => onInputFiles(e.target.files)}
          />
        </label>

        <Typography variant="h4" component="h1" gutterBottom>
          Создать объявление
        </Typography>

        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate sx={{ mt: 1 }}>
          <Box sx={{ display: 'flex', flexDirection: { xs: 'column', md: 'row' }, gap: 3 }}>
            <Box sx={{ width: { xs: '100%', md: '50%' } }}>
              <Box
                onDrop={handleDrop}
                onDragOver={handleDragOver}
                sx={{
                  width: '100%',
                  height: { xs: 320, md: 560 },
                  bgcolor: file ? 'transparent' : (theme: any) => theme.palette.action.hover,
                  borderRadius: 2,
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  overflow: 'hidden',
                  flexDirection: 'column',
                }}
              >
                {file ? (
                  <Box sx={{ width: '100%', height: '100%', display: 'flex', flexDirection: 'column', gap: 1 }}>
                    <Box sx={{ width: '100%', height: '100%', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
                      <Box sx={{ width: '100%', height: '100%', overflow: 'hidden', borderRadius: 1 }}>
                        {previewUrl && (
                          <img
                            src={previewUrl}
                            alt={file.name}
                            style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                          />
                        )}
                      </Box>
                    </Box>

                    <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center', py: 1 }}>
                      <Button size="small" onClick={() => removeFile()} sx={{ minWidth: 0, p: 0.5 }}>
                        ✕
                      </Button>
                    </Box>

                    <Button
                      component="label"
                      htmlFor="create-book-file-input"
                      variant="outlined"
                      startIcon={<PhotoCamera />}
                      sx={{ alignSelf: 'center', mt: 1 }}
                    >
                      Заменить фото
                    </Button>
                  </Box>
                ) : (
                  <Box sx={{ textAlign: 'center', p: 2 }}>
                    <Box sx={{ mb: 1 }}>{/* placeholder */}</Box>
                    <Button
                      component="label"
                      htmlFor="create-book-file-input"
                      variant="outlined"
                      startIcon={<PhotoCamera />}
                      sx={(theme: any) => ({
                        p: 2,
                        border: showImageError ? `2px solid ${theme.palette.error.main}` : undefined,
                        borderRadius: theme.shape.borderRadius,
                        cursor: 'pointer',
                      })}
                    >
                      Загрузить фото
                    </Button>
                    <Typography variant="caption" display="block" sx={{ mt: 1 }}>
                      Перетащите фото сюда или нажмите кнопку
                    </Typography>
                  </Box>
                )}
              </Box>
            </Box>

            <Box sx={{ width: { xs: '100%', md: '50%' } }}>
              <TextField
                margin="normal"
                required
                fullWidth
                label="Название книги"
                {...register('title')}
                error={!!errors.title}
                helperText={String(errors.title?.message || '')}
              />

              <TextField
                margin="normal"
                required
                fullWidth
                label="Автор(ы)"
                {...register('author')}
                error={!!errors.author}
                helperText={String(errors.author?.message || '')}
              />

              <TextField
                margin="normal"
                fullWidth
                label="Описание"
                multiline
                rows={4}
                {...register('description')}
                error={!!errors.description}
                helperText={String(errors.description?.message || '')}
              />

              <FormControl fullWidth margin="normal">
                <Controller
                  name="condition"
                  control={control}
                  render={({ field }) => (
                    <>
                      <Typography variant="subtitle1" sx={{ mb: 1, fontWeight: 500 }}>
                        Состояние
                      </Typography>
                      <ToggleButtonGroup
                        exclusive
                        value={field.value}
                        onChange={(_, val) => {
                          if (val !== null) field.onChange(val);
                        }}
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
                    </>
                  )}
                />
              </FormControl>

              <FormControl fullWidth margin="normal" error={!!errors.locationId}>
                <InputLabel id="location-select-label">Пункт выдачи</InputLabel>
                <Controller
                  name="locationId"
                  control={control}
                  render={({ field }) => (
                    <Select
                      labelId="location-select-label"
                      label="Пункт выдачи"
                      {...field}
                      renderValue={(selected: any) => {
                        if (!selected) return '';
                        const found = locations?.find((l) => l.id === selected as string);
                        return found ? found.name : '';
                      }}
                    >
                      {loadingLocations && <MenuItem value="">Загрузка...</MenuItem>}
                      {!loadingLocations && locations && locations.map((loc) => (
                        <MenuItem key={loc.id} value={loc.id}>
                          {loc.address ? loc.address : loc.name}
                        </MenuItem>
                      ))}
                    </Select>
                  )}
                />
                {errors.locationId && <FormHelperText>{String(errors.locationId.message)}</FormHelperText>}
              </FormControl>

              <Box sx={{ display: 'flex', gap: 2, mt: 3, flexDirection: { xs: 'column', md: 'row' } }}>
                <Button type="submit" variant="contained" color="primary" size="large" sx={{ flex: 1 }} disabled={loading}>
                  {loading ? <CircularProgress size={26} color="inherit" /> : 'Создать'}
                </Button>
                {/* Кнопка отмены здесь была удалена */}
              </Box>
            </Box>
          </Box>
        </Box>
      </Box>
    </Container>
  );
};

export default CreateBookPage;
