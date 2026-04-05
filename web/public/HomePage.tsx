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
  InputAdornment,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  TextField,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';
import ViewAgendaIcon from '@mui/icons-material/ViewAgenda';
import ViewStreamIcon from '@mui/icons-material/ViewStream';
import BookCard from '../shared/components/BookCard';
import useBooksList, { BookSummary } from '../src/hooks/useBooksList';
import SkeletonCard from '../src/components/SkeletonCard';
import api from '../src/services/api';

type ScrollMode = 'continuous' | 'paged';

type BookListResponse = {
  items?: any[];
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

const HomePage: React.FC = () => {
  const { books, loading } = useBooksList();

  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');

  const [mode, setMode] = useState<ScrollMode>('continuous');
  const [sortBy, setSortBy] = useState('');
  const [order, setOrder] = useState<'' | 'asc' | 'desc'>('');
  const [statusFilter, setStatusFilter] = useState('');
  const [locationFilter, setLocationFilter] = useState('');
  const [limit, setLimit] = useState(DEFAULT_LIMIT);

  const [locations, setLocations] = useState<LocationItem[]>([]);
  const [searchResults, setSearchResults] = useState<BookSummary[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState<string | null>(null);
  const [hasSearched, setHasSearched] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  const inputRef = useRef<HTMLInputElement | null>(null);
  const autoLoadRef = useRef<HTMLDivElement | null>(null);

  const hasAnyFilter = useMemo(() => {
    return Boolean(statusFilter || locationFilter || sortBy || order || (mode === 'paged' && limit !== DEFAULT_LIMIT));
  }, [statusFilter, locationFilter, sortBy, order, limit, mode]);

  const showCatalogResults = useMemo(() => {
    return hasSearched || hasAnyFilter || Boolean(submittedQuery);
  }, [hasSearched, hasAnyFilter, submittedQuery]);

  useEffect(() => {
    let active = true;

    const loadLocations = async () => {
      try {
        const resp = await api.get('locations/getAll');
        if (!active) return;
        setLocations(Array.isArray(resp.data) ? resp.data : []);
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

  const normalizeBooks = (items: any[]): BookSummary[] => {
    return items.map((b) => ({
      id: b.id,
      title: b.title,
      author: b.author,
      imageUrl: b.image_url || b.imageUrl || '',
      condition: b.condition,
      locationId: b.current_location_id || b.locationId || undefined,
    }));
  };

  const fetchBooks = async (pageToLoad: number, append: boolean, searchValue?: string) => {
    const cleanSearch = (searchValue ?? submittedQuery).trim();
    const params = new URLSearchParams();
    const pageSize = mode === 'paged' ? limit : DEFAULT_LIMIT;
    params.set('limit', String(pageSize));
    params.set('offset', String((pageToLoad - 1) * pageSize));
    if (sortBy) params.set('sort_by', sortBy);
    if (order) params.set('order', order);
    if (cleanSearch) params.set('search', cleanSearch);
    if (statusFilter) params.set('status', statusFilter);
    if (locationFilter) params.set('location_id', locationFilter);

    const resp = await api.get(`books/list?${params.toString()}`);
    const payload: BookListResponse = resp.data || {};
    const items = normalizeBooks(Array.isArray(payload.items) ? payload.items : []);

    setHasMore(Boolean(payload.has_more));
    setSearchResults((prev) => (append ? [...prev, ...items] : items));
    setCurrentPage(pageToLoad);
  };

  const runFreshSearch = async () => {
    const cleanSearch = query.trim();
    setSubmittedQuery(cleanSearch);
    setHasSearched(true);
    setSearchLoading(true);
    setSearchError(null);

    try {
      await fetchBooks(1, false, cleanSearch);
    } catch (e: any) {
      setSearchError(e?.response?.data?.error || e?.message || 'Ошибка поиска');
      setSearchResults([]);
      setHasMore(false);
    } finally {
      setSearchLoading(false);
    }
  };

  const loadNext = useCallback(async () => {
    if (searchLoading || !hasMore) return;

    setSearchLoading(true);
    setSearchError(null);

    try {
      await fetchBooks(currentPage + 1, mode === 'continuous');
    } catch (e: any) {
      setSearchError(e?.response?.data?.error || e?.message || 'Не удалось загрузить следующую страницу');
    } finally {
      setSearchLoading(false);
    }
  }, [searchLoading, hasMore, currentPage, mode, submittedQuery, sortBy, order, statusFilter, locationFilter, limit]);

  const goToPage = useCallback(async (page: number) => {
    if (searchLoading || page < 1) return;

    setSearchLoading(true);
    setSearchError(null);
    try {
      await fetchBooks(page, false);
    } catch (e: any) {
      setSearchError(e?.response?.data?.error || e?.message || 'Ошибка переключения страницы');
    } finally {
      setSearchLoading(false);
    }
  }, [searchLoading, submittedQuery, sortBy, order, statusFilter, locationFilter, limit]);

  const resetFilters = () => {
    setSortBy('');
    setOrder('');
    setStatusFilter('');
    setLocationFilter('');
    setLimit(DEFAULT_LIMIT);
  };

  useEffect(() => {
    if (!showCatalogResults) return;

    setSearchLoading(true);
    setSearchError(null);
    fetchBooks(1, false)
      .catch((e: any) => {
        setSearchResults([]);
        setHasMore(false);
        setSearchError(e?.response?.data?.error || e?.message || 'Ошибка обновления выдачи');
      })
      .finally(() => setSearchLoading(false));
  }, [showCatalogResults, mode, sortBy, order, statusFilter, locationFilter, limit]);

  useEffect(() => {
    if (mode !== 'continuous' || !hasSearched || searchLoading || !hasMore || searchResults.length === 0) return;

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
  }, [mode, hasSearched, searchLoading, hasMore, searchResults.length, loadNext]);

  // Focus input when user presses '/'
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      // ignore if focus is already in an input/textarea or user uses modifier
      const tag = (document.activeElement && document.activeElement.tagName) || '';
      if (tag === 'INPUT' || tag === 'TEXTAREA' || (e.ctrlKey || e.metaKey || e.altKey)) return;
      if (e.key === '/') {
        e.preventDefault();
        inputRef.current?.focus();
      }
    };
    window.addEventListener('keydown', handler);
    return () => window.removeEventListener('keydown', handler);
  }, []);

  return (
    <Container maxWidth="xl" sx={{ py: { xs: 2, md: 4 } }}>
      <Box sx={{ mb: 3 }}>
        <TextField
          fullWidth
          variant="outlined"
          value={query}
          inputRef={inputRef}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => { if (e.key === 'Enter') runFreshSearch(); }}
          placeholder="Ищите по названию или автору"
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <IconButton
                  size="small"
                  onClick={runFreshSearch}
                  aria-label="поиск"
                >
                  {searchLoading ? <CircularProgress size={18} /> : <SearchIcon />}
                </IconButton>
              </InputAdornment>
            ),
            endAdornment: query ? (
              <InputAdornment position="end">
                <IconButton size="small" onClick={() => { setQuery(''); inputRef.current?.focus(); }} aria-label="очистить поиск">
                  <CloseIcon fontSize="small" />
                </IconButton>
              </InputAdornment>
            ) : null,
          }}
        />
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
          {searchError && <Alert severity="error" sx={{ mb: 2 }}>{searchError}</Alert>}

          {(submittedQuery || hasAnyFilter) && (
            <Stack direction="row" spacing={1} sx={{ mb: 2, flexWrap: 'wrap', gap: 1 }}>
              {submittedQuery && <Chip label={`Запрос: ${submittedQuery}`} size="small" />}
              {statusFilter && <Chip label={`Статус: ${STATUS_OPTIONS.find((s) => s.value === statusFilter)?.label || statusFilter}`} size="small" />}
              {locationFilter && <Chip label={`Локация: ${locations.find((l) => l.id === locationFilter)?.name || 'выбрана'}`} size="small" />}
              {sortBy && <Chip label={`Сортировка: ${SORT_OPTIONS.find((s) => s.value === sortBy)?.label || sortBy}`} size="small" />}
              {order && <Chip label={`Порядок: ${order === 'asc' ? 'возрастание' : 'убывание'}`} size="small" />}
            </Stack>
          )}

          <Box
            sx={{
              display: 'grid',
              gap: 3,
              gridTemplateColumns: { xs: 'repeat(2, 1fr)', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)', xl: 'repeat(4, 1fr)' },
            }}
          >
            {showCatalogResults && searchLoading && Array.from({ length: 8 }).map((_, i) => <SkeletonCard key={i} />)}

            {showCatalogResults && !searchLoading && searchResults.length === 0 && (
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
                  Попробуйте убрать часть фильтров, изменить сортировку или сократить запрос.
                </Typography>
              </Box>
            )}

            {showCatalogResults && !searchLoading && searchResults.map((b: BookSummary) => (
              <BookCard key={b.id} id={b.id as any} imageUrl={b.imageUrl || ''} title={b.title} author={b.author} />
            ))}

            {!showCatalogResults && loading && Array.from({ length: 8 }).map((_, i) => <SkeletonCard key={i} />)}
            {!showCatalogResults && !loading && (!books || books.length === 0) && (
              <Box sx={{ gridColumn: '1 / -1', p: 2 }}>Нет книг для отображения.</Box>
            )}
            {!showCatalogResults && !loading && books && books.map((b: BookSummary) => (
              <BookCard key={b.id} id={b.id as any} imageUrl={b.imageUrl || ''} title={b.title} author={b.author} />
            ))}
          </Box>

          {showCatalogResults && mode === 'continuous' && hasMore && (
            <>
              <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center' }}>
                <Button variant="outlined" onClick={loadNext} disabled={searchLoading || !hasMore}>
                  {searchLoading ? 'Загрузка...' : 'Показать еще'}
                </Button>
              </Box>
              <Box ref={autoLoadRef} sx={{ height: 1, mt: 2 }} />
            </>
          )}

          {showCatalogResults && mode === 'paged' && searchResults.length > 0 && (
            <Box sx={{ mt: 3, display: 'flex', justifyContent: 'center', gap: 1.5, alignItems: 'center', flexWrap: 'wrap' }}>
              <Button variant="outlined" onClick={() => goToPage(currentPage - 1)} disabled={searchLoading || currentPage <= 1}>
                Предыдущая
              </Button>
              <Typography variant="body2" sx={{ minWidth: 110, textAlign: 'center', fontWeight: 600 }}>
                Страница {currentPage}
              </Typography>
              <Button variant="contained" onClick={() => goToPage(currentPage + 1)} disabled={searchLoading || !hasMore}>
                Следующая
              </Button>
            </Box>
          )}
        </Box>

        <Box sx={{ position: { lg: 'sticky' }, top: { lg: 96 } }}>
          <Box
            sx={{
              p: 2,
              borderRadius: 2,
              bgcolor: 'background.paper',
              border: '1px solid',
              borderColor: 'divider',
              boxShadow: '0 4px 12px rgba(15, 23, 42, 0.04)',
            }}
          >
            <Typography variant="h6" sx={{ mb: 2, fontWeight: 700 }}>
              Фильтры и сортировка
            </Typography>

            <Stack spacing={1.5}>
              <FormControl size="small" fullWidth>
                <InputLabel id="home-sort-by-label">Сортировать</InputLabel>
                <Select
                  labelId="home-sort-by-label"
                  value={sortBy}
                  label="Сортировать"
                  onChange={(e) => {
                    const nextValue = e.target.value;
                    setSortBy(nextValue);
                    if (!nextValue) setOrder('');
                  }}
                >
                  <MenuItem value="">Без сортировки</MenuItem>
                  {SORT_OPTIONS.map((option) => (
                    <MenuItem key={option.value} value={option.value}>{option.label}</MenuItem>
                  ))}
                </Select>
              </FormControl>

              {sortBy && (
                <FormControl size="small" fullWidth>
                  <InputLabel id="home-order-label">Порядок</InputLabel>
                  <Select
                    labelId="home-order-label"
                    value={order}
                    label="Порядок"
                    onChange={(e) => setOrder(e.target.value as '' | 'asc' | 'desc')}
                  >
                    <MenuItem value="">По умолчанию</MenuItem>
                    <MenuItem value="desc">По убыванию</MenuItem>
                    <MenuItem value="asc">По возрастанию</MenuItem>
                  </Select>
                </FormControl>
              )}

              <FormControl size="small" fullWidth>
                <InputLabel id="home-status-label">Статус</InputLabel>
                <Select
                  labelId="home-status-label"
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
                <InputLabel id="home-location-label">Локация</InputLabel>
                <Select
                  labelId="home-location-label"
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

            </Stack>

            <Typography variant="subtitle2" sx={{ mt: 2.5, mb: 1, fontWeight: 700 }}>
              Режим выдачи
            </Typography>

            <ToggleButtonGroup
              size="small"
              color="primary"
              exclusive
              value={mode}
              onChange={(_, nextMode: ScrollMode | null) => {
                if (nextMode) setMode(nextMode);
              }}
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

            {mode === 'paged' && (
              <FormControl size="small" fullWidth sx={{ mt: 2 }}>
                <InputLabel id="home-limit-label">На странице</InputLabel>
                <Select
                  labelId="home-limit-label"
                  value={String(limit)}
                  label="На странице"
                  onChange={(e) => setLimit(Number(e.target.value))}
                >
                  {LIMIT_OPTIONS.map((value) => (
                    <MenuItem key={value} value={String(value)}>{value}</MenuItem>
                  ))}
                </Select>
              </FormControl>
            )}

            <Stack direction={{ xs: 'column', sm: 'row', lg: 'column' }} spacing={1} sx={{ mt: 2.5 }}>
              <Button variant="text" onClick={resetFilters} disabled={!hasAnyFilter}>Сбросить фильтры</Button>
            </Stack>
          </Box>
        </Box>
      </Box>
    </Container>
  );
};

export default HomePage;
