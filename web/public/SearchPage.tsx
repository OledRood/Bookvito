import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import {
  Alert,
  Box,
  Button,
  Chip,
  CircularProgress,
  Container,
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  Select,
  Skeleton,
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import ViewAgendaIcon from '@mui/icons-material/ViewAgenda';
import ViewStreamIcon from '@mui/icons-material/ViewStream';
import BookCard from './BookCard';
import api from '../src/services/api';

type ScrollMode = 'continuous' | 'paged';

type BookItem = {
  id: string;
  title: string;
  author: string;
  image_url?: string;
  imageUrl?: string;
};

type BookListResponse = {
  items?: BookItem[];
  has_more?: boolean;
};

type LocationItem = {
  id: string;
  name?: string;
  address?: string;
};

const DEFAULT_LIMIT = 16;
const LIMIT_OPTIONS = [8, 12, 16, 24, 32];

const SORT_OPTIONS = [
  { value: 'created_at', label: 'Сначала новые' },
  { value: 'updated_at', label: 'Недавно обновленные' },
  { value: 'title', label: 'По названию' },
  { value: 'author', label: 'По автору' },
  { value: 'status', label: 'По статусу' },
  { value: 'current_location_id', label: 'По локации' },
];

const STATUS_OPTIONS = [
  { value: '', label: 'Все статусы' },
  { value: 'available', label: 'Доступна' },
];

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');
  const [mode, setMode] = useState<ScrollMode>('continuous');
  const [sortBy, setSortBy] = useState('created_at');
  const [order, setOrder] = useState<'asc' | 'desc'>('desc');
  const [statusFilter, setStatusFilter] = useState('');
  const [locationFilter, setLocationFilter] = useState('');
  const [limit, setLimit] = useState(DEFAULT_LIMIT);

  const [locations, setLocations] = useState<LocationItem[]>([]);

  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<BookItem[]>([]);
  const [hasMore, setHasMore] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [error, setError] = useState<string | null>(null);
  const [hasSearched, setHasSearched] = useState(false);
  const autoLoadRef = useRef<HTMLDivElement | null>(null);

  const hasAnyFilter = useMemo(() => {
    return Boolean(statusFilter || locationFilter || sortBy !== 'created_at' || order !== 'desc' || limit !== DEFAULT_LIMIT);
  }, [statusFilter, locationFilter, sortBy, order, limit]);

  const showCatalogResults = useMemo(() => {
    return hasSearched || hasAnyFilter || Boolean(submittedQuery);
  }, [hasSearched, hasAnyFilter, submittedQuery]);

  useEffect(() => {
    let active = true;

    const loadLocations = async () => {
      try {
        const resp = await api.get('locations/getAll');
        if (!active) return;
        const list = Array.isArray(resp.data) ? resp.data : [];
        setLocations(list);
      } catch {
        if (!active) return;
        setLocations([]);
      }
    };

    loadLocations();

    return () => {
      active = false;
    };
  }, []);

  const fetchBooks = async (pageToLoad: number, append: boolean, searchValue?: string) => {
    const cleanSearch = (searchValue ?? submittedQuery).trim();
    const params = new URLSearchParams();
    params.set('limit', String(limit));
    params.set('offset', String((pageToLoad - 1) * limit));
    params.set('sort_by', sortBy);
    params.set('order', order);
    if (cleanSearch) params.set('search', cleanSearch);
    if (statusFilter) params.set('status', statusFilter);
    if (locationFilter) params.set('location_id', locationFilter);

    const resp = await api.get(`books/list?${params.toString()}`);
    const payload: BookListResponse = resp.data || {};
    const items = Array.isArray(payload.items) ? payload.items : [];
    const nextHasMore = Boolean(payload.has_more);

    setHasMore(nextHasMore);
    setResults((prev) => (append ? [...prev, ...items] : items));
    setCurrentPage(pageToLoad);
  };

  const runFreshSearch = async () => {
    const cleanSearch = query.trim();

    setSubmittedQuery(cleanSearch);
    setHasSearched(true);
    setLoading(true);
    setError(null);

    try {
      await fetchBooks(1, false, cleanSearch);
    } catch (e: any) {
      console.error('Ошибка поиска', e);
      setResults([]);
      setHasMore(false);
      setError(e?.response?.data?.error || e?.message || 'Ошибка поиска');
    } finally {
      setLoading(false);
    }
  };

  const loadNext = useCallback(async () => {
    if (loading || !hasMore) return;

    setLoading(true);
    setError(null);
    try {
      await fetchBooks(currentPage + 1, mode === 'continuous');
    } catch (e: any) {
      setError(e?.response?.data?.error || e?.message || 'Не удалось загрузить следующую страницу');
    } finally {
      setLoading(false);
    }
  }, [loading, hasMore, currentPage, mode]);

  const goToPage = useCallback(async (page: number) => {
    if (loading || page < 1) return;
    setLoading(true);
    setError(null);
    try {
      await fetchBooks(page, false);
    } catch (e: any) {
      setError(e?.response?.data?.error || e?.message || 'Ошибка переключения страницы');
    } finally {
      setLoading(false);
    }
  }, [loading]);

  const resetFilters = () => {
    setSortBy('created_at');
    setOrder('desc');
    setStatusFilter('');
    setLocationFilter('');
    setLimit(DEFAULT_LIMIT);
  };

  useEffect(() => {
    if (!showCatalogResults) return;

    setLoading(true);
    setError(null);
    fetchBooks(1, false)
      .catch((e: any) => {
        setResults([]);
        setHasMore(false);
        setError(e?.response?.data?.error || e?.message || 'Ошибка обновления выдачи');
      })
      .finally(() => setLoading(false));
  }, [showCatalogResults, mode, sortBy, order, statusFilter, locationFilter, limit]);

  useEffect(() => {
    if (mode !== 'continuous' || !hasSearched || loading || !hasMore || results.length === 0) return;
    const node = autoLoadRef.current;
    if (!node) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting) {
          loadNext();
        }
      },
      { root: null, rootMargin: '240px 0px', threshold: 0 }
    );

    observer.observe(node);
    return () => observer.disconnect();
  }, [mode, hasSearched, loading, hasMore, results.length, loadNext]);

  const onToggleMode = (_: React.MouseEvent<HTMLElement>, nextMode: ScrollMode | null) => {
    if (!nextMode || nextMode === mode) return;
    setMode(nextMode);
  };

  return (
    <Container maxWidth="xl" sx={{ py: { xs: 2, md: 4 } }}>
      <Box
        sx={{
          mb: 3,
          p: { xs: 2, md: 3 },
          borderRadius: 4,
          border: '1px solid',
          borderColor: 'divider',
          background: 'linear-gradient(135deg, rgba(20,83,45,0.08) 0%, rgba(255,255,255,0.96) 58%, rgba(14,165,233,0.08) 100%)',
        }}
      >
        <Typography variant="overline" sx={{ color: 'text.secondary', letterSpacing: '0.1em' }}>
          Catalog View
        </Typography>
        <Typography variant="h4" component="h1" sx={{ mt: 0.5, mb: 1, fontWeight: 700 }}>
          Страница поиска
        </Typography>
        <Typography color="text.secondary" sx={{ maxWidth: 760, mb: 2.5 }}>
          Основной поиск теперь живёт на главной странице, а здесь оставлен расширенный режим каталога на том же backend endpoint `books/list`.
        </Typography>

        <Box sx={{ display: 'flex', gap: 1 }}>
          <TextField
            fullWidth
            placeholder="Введите название или автора"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === 'Enter') runFreshSearch();
            }}
          />
          <IconButton color="primary" onClick={runFreshSearch} disabled={loading} aria-label="поиск">
            {loading ? <CircularProgress size={24} /> : <SearchIcon />}
          </IconButton>
        </Box>
      </Box>

      <Box
        sx={{
          display: 'grid',
          gap: 3,
          gridTemplateColumns: { xs: '1fr', lg: 'minmax(0, 1fr) 320px' },
          alignItems: 'start',
        }}
      >
        <Box>
          {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

          {!showCatalogResults && !loading && (
            <Box
              sx={{
                mb: 2.5,
                p: 2,
                borderRadius: 3,
                bgcolor: 'background.paper',
                border: '1px solid',
                borderColor: 'divider',
              }}
            >
              <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 0.5 }}>
                Расширенный каталог
              </Typography>
              <Typography color="text.secondary">
                Здесь можно тестировать фильтры и пагинацию, но основной пользовательский поиск стоит вести с HomeScreen.
              </Typography>
            </Box>
          )}

          {(submittedQuery || hasAnyFilter) && (
            <Stack direction="row" spacing={1} sx={{ mb: 2, flexWrap: 'wrap', gap: 1 }}>
              {submittedQuery && <Chip label={`Запрос: ${submittedQuery}`} size="small" />}
              {statusFilter && <Chip label={`Статус: ${STATUS_OPTIONS.find((s) => s.value === statusFilter)?.label || statusFilter}`} size="small" />}
              {locationFilter && <Chip label={`Локация: ${locations.find((l) => l.id === locationFilter)?.name || 'выбрана'}`} size="small" />}
              <Chip label={`Сортировка: ${SORT_OPTIONS.find((s) => s.value === sortBy)?.label || sortBy}`} size="small" />
              <Chip label={`Порядок: ${order === 'asc' ? 'возрастание' : 'убывание'}`} size="small" />
              <Chip label={`На странице: ${limit}`} size="small" />
            </Stack>
          )}

          <Box
            sx={{
              display: 'grid',
              gap: 2,
              gridTemplateColumns: {
                xs: 'repeat(2, 1fr)',
                sm: 'repeat(3, 1fr)',
                md: 'repeat(4, 1fr)'
              }
            }}
          >
            {showCatalogResults && loading && results.length === 0 && Array.from({ length: 8 }).map((_, idx) => (
              <Box key={idx}>
                <Skeleton variant="rounded" height={220} sx={{ mb: 1 }} />
                <Skeleton width="80%" />
                <Skeleton width="60%" />
              </Box>
            ))}

            {showCatalogResults && !loading && results.length === 0 && (
              <Box
                sx={{
                  gridColumn: '1 / -1',
                  p: 4,
                  textAlign: 'center',
                  borderRadius: 3,
                  border: '1px dashed',
                  borderColor: 'divider',
                  bgcolor: 'background.paper',
                }}
              >
                <Typography variant="h6" sx={{ mb: 1 }}>Ничего не найдено</Typography>
                <Typography color="text.secondary">
                  Попробуйте сменить запрос или ослабить фильтры.
                </Typography>
              </Box>
            )}

            {results.map((b) => (
              <Box key={b.id}>
                <BookCard id={b.id} imageUrl={b.image_url || b.imageUrl || ''} title={b.title} author={b.author} />
              </Box>
            ))}
          </Box>

          {showCatalogResults && mode === 'continuous' && results.length > 0 && (
            <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center' }}>
              <Button variant="outlined" onClick={loadNext} disabled={loading || !hasMore}>
                {loading ? 'Загрузка...' : hasMore ? 'Показать еще' : 'Больше результатов нет'}
              </Button>
            </Box>
          )}

          {showCatalogResults && mode === 'continuous' && hasMore && (
            <Box ref={autoLoadRef} sx={{ height: 1, mt: 2 }} />
          )}

          {showCatalogResults && mode === 'paged' && results.length > 0 && (
            <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center', gap: 1.5, alignItems: 'center', flexWrap: 'wrap' }}>
              <Button variant="outlined" onClick={() => goToPage(currentPage - 1)} disabled={loading || currentPage <= 1}>
                Предыдущая
              </Button>
              <Typography variant="body2" sx={{ minWidth: 110, textAlign: 'center', fontWeight: 600 }}>
                Страница {currentPage}
              </Typography>
              <Button variant="contained" onClick={() => goToPage(currentPage + 1)} disabled={loading || !hasMore}>
                Следующая
              </Button>
            </Box>
          )}
        </Box>

        <Box sx={{ position: { lg: 'sticky' }, top: { lg: 96 } }}>
          <Box
            sx={{
              p: 2.5,
              borderRadius: 4,
              bgcolor: 'background.paper',
              border: '1px solid',
              borderColor: 'divider',
              boxShadow: '0 16px 40px rgba(15, 23, 42, 0.06)',
            }}
          >
            <Typography variant="overline" sx={{ color: 'text.secondary', letterSpacing: '0.1em' }}>
              Инструменты
            </Typography>
            <Typography variant="h6" sx={{ mb: 2, fontWeight: 700 }}>
              Фильтры каталога
            </Typography>

            <Stack spacing={1.5}>
              <FormControl size="small" fullWidth>
                <InputLabel id="search-sort-by-label">Сортировать</InputLabel>
                <Select
                  labelId="search-sort-by-label"
                  value={sortBy}
                  label="Сортировать"
                  onChange={(e) => setSortBy(e.target.value)}
                >
                  {SORT_OPTIONS.map((option) => (
                    <MenuItem key={option.value} value={option.value}>{option.label}</MenuItem>
                  ))}
                </Select>
              </FormControl>

              <FormControl size="small" fullWidth>
                <InputLabel id="search-order-label">Порядок</InputLabel>
                <Select
                  labelId="search-order-label"
                  value={order}
                  label="Порядок"
                  onChange={(e) => setOrder(e.target.value as 'asc' | 'desc')}
                >
                  <MenuItem value="desc">По убыванию</MenuItem>
                  <MenuItem value="asc">По возрастанию</MenuItem>
                </Select>
              </FormControl>

              <FormControl size="small" fullWidth>
                <InputLabel id="search-status-label">Статус</InputLabel>
                <Select
                  labelId="search-status-label"
                  value={statusFilter}
                  label="Статус"
                  onChange={(e) => setStatusFilter(e.target.value)}
                >
                  {STATUS_OPTIONS.map((option) => (
                    <MenuItem key={option.value || 'all'} value={option.value}>{option.label}</MenuItem>
                  ))}
                </Select>
              </FormControl>

              <FormControl size="small" fullWidth>
                <InputLabel id="search-location-label">Локация</InputLabel>
                <Select
                  labelId="search-location-label"
                  value={locationFilter}
                  label="Локация"
                  onChange={(e) => setLocationFilter(e.target.value)}
                >
                  <MenuItem value="">Все локации</MenuItem>
                  {locations.map((location) => (
                    <MenuItem key={location.id} value={location.id}>
                      {(location.name || 'Локация')} {location.address ? `(${location.address})` : ''}
                    </MenuItem>
                  ))}
                </Select>
              </FormControl>

              <FormControl size="small" fullWidth>
                <InputLabel id="search-limit-label">На странице</InputLabel>
                <Select
                  labelId="search-limit-label"
                  value={String(limit)}
                  label="На странице"
                  onChange={(e) => setLimit(Number(e.target.value))}
                >
                  {LIMIT_OPTIONS.map((value) => (
                    <MenuItem key={value} value={String(value)}>{value}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            </Stack>

            <Typography variant="subtitle2" sx={{ mt: 2.5, mb: 1, fontWeight: 700 }}>
              Режим просмотра
            </Typography>

            <ToggleButtonGroup
              size="small"
              color="primary"
              exclusive
              value={mode}
              onChange={onToggleMode}
              aria-label="режим прокрутки"
              fullWidth
            >
              <ToggleButton value="continuous" aria-label="непрерывный режим">
                <ViewStreamIcon sx={{ mr: 0.75 }} fontSize="small" />
                Лента
              </ToggleButton>
              <ToggleButton value="paged" aria-label="страничный режим">
                <ViewAgendaIcon sx={{ mr: 0.75 }} fontSize="small" />
                Страницы
              </ToggleButton>
            </ToggleButtonGroup>

            <Button variant="text" onClick={resetFilters} disabled={!hasAnyFilter} sx={{ mt: 2, width: '100%' }}>
              Сбросить фильтры
            </Button>
          </Box>
        </Box>
      </Box>
    </Container>
  );
};

export default SearchPage;
