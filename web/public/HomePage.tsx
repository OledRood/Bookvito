import React, { useState, useRef, useEffect } from 'react';
import { Box, Typography, TextField, InputAdornment, Chip, Container, IconButton } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import CloseIcon from '@mui/icons-material/Close';
import BookCard from '../shared/components/BookCard';
import useBooksList, { BookSummary } from '../src/hooks/useBooksList';
import SkeletonCard from '../src/components/SkeletonCard';
import api from '../src/services/api';

const HomePage: React.FC = () => {
  const { books, loading } = useBooksList();
  // Filters are removed for now (kept as a commented block for future use)
  const [query, setQuery] = useState('');
  const [searchResults, setSearchResults] = useState<BookSummary[] | null>(null);
  const [searchLoading, setSearchLoading] = useState(false);
  const inputRef = useRef<HTMLInputElement | null>(null);
  const debounceRef = useRef<number | null>(null);

  const doSearch = async () => {
    const q = query.trim();
    if (!q) {
      setSearchResults(null);
      return;
    }
    setSearchLoading(true);
    try {
      const resp = await api.get(`books/search?q=${encodeURIComponent(q)}&limit=50`);
      const data = resp.data || [];
      const normalized = data.map((b: any) => ({
        id: b.id,
        title: b.title,
        author: b.author,
        imageUrl: b.image_url || b.imageUrl || '',
      }));
      setSearchResults(normalized);
    } catch (e: any) {
      // swallow error for now, could show notification
      setSearchResults([]);
    } finally {
      setSearchLoading(false);
    }
  };

  // Debounced search: trigger doSearch 300ms after user stops typing
  useEffect(() => {
    // clear previous timer
    if (debounceRef.current) {
      window.clearTimeout(debounceRef.current);
      debounceRef.current = null;
    }

    if (!query || query.trim() === '') {
      setSearchResults(null);
      return;
    }

    // set new timer
    debounceRef.current = window.setTimeout(() => {
      doSearch();
    }, 300);

    return () => {
      if (debounceRef.current) {
        window.clearTimeout(debounceRef.current);
        debounceRef.current = null;
      }
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query]);

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
    <Container maxWidth="lg">
      {/* Поиск */}
      <Box sx={{ my: 2 }}>
        <TextField
          fullWidth
          variant="outlined"
          value={query}
          inputRef={inputRef}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => { if (e.key === 'Enter') { if (debounceRef.current) { window.clearTimeout(debounceRef.current); debounceRef.current = null; } doSearch(); } }}
          placeholder="Поиск по названию или автору..."
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <IconButton
                  size="small"
                  onClick={() => {
                    const q = query.trim();
                    if (q) doSearch();
                    else inputRef.current?.focus();
                  }}
                  aria-label="поиск"
                >
                  <SearchIcon />
                </IconButton>
              </InputAdornment>
            ),
            endAdornment: (
              query ? (
                <InputAdornment position="end">
                  <IconButton size="small" onClick={() => { setQuery(''); setSearchResults(null); inputRef.current?.focus(); }} aria-label="очистить поиск">
                    <CloseIcon fontSize="small" />
                  </IconButton>
                </InputAdornment>
              ) : null
            ),
          }}
        />
      </Box>

      {/*
        Фильтры (визуальные заглушки) — временно закомментированы.
        Если нужно вернуть фильтры, раскомментируйте блок ниже.

        <Box sx={{ my: 2, display: 'flex', gap: 1, flexWrap: 'wrap' }}>
          <Typography variant="subtitle1" sx={{ mr: 1, alignSelf: 'center' }}>
            Фильтры:
          </Typography>
          <Chip label="Жанр" onClick={() => {}} />
          <Chip label="Год издания" onClick={() => {}} />
          <Chip label="Состояние" variant="outlined" onClick={() => {}} />
          <Chip label="Бесплатно" variant="outlined" onClick={() => {}} />
        </Box>
      */}

      {/* Сетка книг: responsive grid 2..5 columns */}
      <Box
        sx={{
          display: 'grid',
          gap: 3,
          gridTemplateColumns: { xs: 'repeat(2, 1fr)', sm: 'repeat(2, 1fr)', md: 'repeat(3, 1fr)', lg: 'repeat(4, 1fr)', xl: 'repeat(5, 1fr)' },
        }}
      >
        {searchLoading && Array.from({ length: 8 }).map((_, i) => <SkeletonCard key={i} />)}
        {!searchLoading && (searchResults ? (
          searchResults.length === 0 ? (
            <Box sx={{ gridColumn: '1 / -1', p: 2 }}>Ничего не найдено.</Box>
          ) : (
            searchResults.map((b: BookSummary) => (
              <BookCard key={b.id} id={b.id as any} imageUrl={b.imageUrl || ''} title={b.title} author={b.author} />
            ))
          )
        ) : (
          // no active search — show default summary list
          loading && Array.from({ length: 8 }).map((_, i) => <SkeletonCard key={i} />)
        ))}
        {!searchResults && !loading && (!books || books.length === 0) && (
          <Box sx={{ gridColumn: '1 / -1', p: 2 }}>Нет книг для отображения.</Box>
        )}
        {!searchResults && !loading && books && books.map((b: BookSummary) => (
          <BookCard key={b.id} id={b.id as any} imageUrl={b.imageUrl || ''} title={b.title} author={b.author} />
        ))}
      </Box>
    </Container>
  );
};

export default HomePage;

// Component state and helpers (placed after export to keep top simple)