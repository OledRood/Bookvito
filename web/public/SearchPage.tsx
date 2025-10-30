import React, { useState } from 'react';
import { Container, TextField, IconButton, CircularProgress, Typography, Box } from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import BookCard from './BookCard';

const SearchPage: React.FC = () => {
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [results, setResults] = useState<any[]>([]);
  const [error, setError] = useState<string | null>(null);

  const doSearch = async () => {
    if (!query || query.trim() === '') return;
    setLoading(true);
    setError(null);
    try {
      const params = new URLSearchParams({ q: query, limit: '50' });
      const token = localStorage.getItem('accessToken') || localStorage.getItem('token');
      const headers: Record<string, string> = {};
      if (token) headers['Authorization'] = `Bearer ${token}`;
      const resp = await fetch(`/api/v1/books/search?${params.toString()}`, { headers });
      if (!resp.ok) {
        const txt = await resp.text();
        throw new Error(txt || 'Ошибка поиска');
      }
      const data = await resp.json();
      setResults(Array.isArray(data) ? data : []);
    } catch (e: any) {
      console.error('Ошибка поиска', e);
      setError(e.message || 'Ошибка поиска');
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="lg" sx={{ pt: 3 }}>
      <Typography variant="h4" component="h1" sx={{ mb: 2 }}>
        Поиск книг
      </Typography>

      <Box sx={{ display: 'flex', gap: 1, mb: 3 }}>
        <TextField
          fullWidth
          placeholder="Введите название, автора или описание"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') doSearch();
          }}
        />
  <IconButton color="primary" onClick={doSearch} disabled={loading} aria-label="поиск">
          {loading ? <CircularProgress size={24} /> : <SearchIcon />}
        </IconButton>
      </Box>

      {error && <Typography color="error">{error}</Typography>}

      {results.length === 0 && !loading ? (
        <Typography color="text.secondary">Ничего не найдено</Typography>
      ) : (
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
          {results.map((b) => (
            <Box key={b.id}>
              <BookCard id={b.id} imageUrl={b.image_url || b.imageUrl || ''} title={b.title} author={b.author} />
            </Box>
          ))}
        </Box>
      )}
    </Container>
  );
};

export default SearchPage;
