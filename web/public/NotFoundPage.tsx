import React, { useEffect, useState } from 'react';
import { Box, Typography, Button, useTheme } from '@mui/material';
import { keyframes } from '@emotion/react';
// static bundled fallback so background text works even when origin/static server is down
import staticDostoevsky from '../src/assets/dostoevsky.txt?raw';
import { Link } from 'react-router-dom';

const NotFoundPage: React.FC = () => {
  const theme = useTheme();
  // start with bundled text so page shows text even if fetch fails
  const [bgText, setBgText] = useState<string>(staticDostoevsky?.trim() ?? '');

  useEffect(() => {
    let cancelled = false;
    // try to fetch an external /public override (optional). If not present, we keep the bundled text.
    fetch('/dostoevsky.txt')
      .then((res) => {
        if (!res.ok) throw new Error('not found');
        const contentType = res.headers.get('content-type') || '';
        return res.text().then((txt) => ({ txt, contentType }));
      })
      .then(({ txt, contentType }) => {
        if (cancelled) return;
        const trimmed = txt.trim();
        // simple heuristic: ignore responses that look like HTML (index.html, error pages, etc.)
        const looksLikeHtml = /<!doctype|<html|<script|<body|<div/i.test(trimmed);
        const isPlainText = contentType.includes('text/plain') || (!looksLikeHtml && trimmed.length > 0);
        if (isPlainText) {
          setBgText(trimmed);
        } else {
          // Don't overwrite the bundled text with HTML page content served by proxies/static servers
          // Helpful during dev when origin may return index.html for unknown paths.
          // eslint-disable-next-line no-console
          console.warn('[NotFoundPage] fetched /dostoevsky.txt looks like HTML or non-plain content — ignoring.');
        }
      })
      .catch(() => {
        // ignore - fallback (bundled text) will be used
      });
    return () => {
      cancelled = true;
    };
  }, []);

  const repeatedText = React.useMemo(() => {
    const base = bgText || '«Бесы» — Ф. М. Достоевский';
    const count = 30;
    return Array(count).fill(base).join('\n\n');
  }, [bgText]);

  const gradientShift = keyframes`
    0% { background-position: 0% 50%; }
    50% { background-position: 100% 50%; }
    100% { background-position: 0% 50%; }
  `;

  return (
    <Box
      sx={{
        p: 3,
        textAlign: 'center',
  mt: 0,
  minHeight: '100vh',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        bgcolor: 'background.default',
        position: 'relative',
        overflow: 'hidden',
      }}
    >
      <Box
        component="pre"
        aria-hidden
        sx={{
          position: 'absolute',
          top: 0,
          left: 0,
          width: '100%',
          height: '100%',
          m: 0,
          p: 4,
          overflow: 'hidden',
          whiteSpace: 'pre-wrap',
          fontSize: '14px',
          lineHeight: 1.6,
          pointerEvents: 'none',
          zIndex: 0,
          // increase opacity so background lines are visible on screenshots/low-contrast displays
          color: theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.12)' : 'rgba(0,0,0,0.12)',
          textAlign: 'left',
        }}
      >
        {repeatedText}
      </Box>

      <Box sx={{ position: 'relative', zIndex: 1, width: '100%' }}>
        <Typography
          variant="h2"
          component="h1"
          gutterBottom
          sx={{
            fontWeight: 800,
            lineHeight: 1,
            fontSize: 'clamp(96px, 24vw, 220px)',
            // animated gradient text for the numeric code
            background: 'linear-gradient(90deg, var(--md-sys-color-primary, #7c3aed) 0%, var(--md-sys-color-surface-tint, #06b6d4) 50%, var(--md-sys-color-secondary, #34d399) 100%)',
            WebkitBackgroundClip: 'text',
            backgroundClip: 'text',
            color: 'transparent',
            backgroundSize: '200% 200%',
            animation: `${gradientShift} 6s linear infinite`,
          }}
        >
          404
        </Typography>

        <Typography variant="h5" gutterBottom>
          Страница не найдена
        </Typography>

        <Button
          component={Link}
          to="/"
          variant="contained"
          color="primary"
          sx={{
            mt: 4,
            // match size used on ServerCheck's button (smaller)
            px: 5,
            py: 1.8,
            fontSize: '16px',
            textTransform: 'uppercase',
            borderRadius: 3,
            // allow natural sizing similar to ServerCheck
            width: 'auto',
            minWidth: 140,
          }}
        >
          На главную
        </Button>
      </Box>
    </Box>
  );
};

export default NotFoundPage;