import React, { useEffect, useState, PropsWithChildren } from 'react';
import { Box, Button, CircularProgress, useTheme } from '@mui/material';
import { keyframes } from '@emotion/react';
// bundled fallback so background lines show even when origin/public fetch fails
import staticDostoevsky from '../src/assets/dostoevsky.txt?raw';
import api from '../src/services/api';

const float = keyframes`
  0% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
  100% { transform: translateY(0); }
`;

const gradientShift = keyframes`
  0% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
  100% { background-position: 0% 50%; }
`;

const glow = keyframes`
  0% { filter: drop-shadow(0 0 0px rgba(0,0,0,0.0)); }
  50% { filter: drop-shadow(0 12px 40px rgba(0,0,0,0.18)); }
  100% { filter: drop-shadow(0 0 0px rgba(0,0,0,0.0)); }
`;

const buttonBeat = keyframes`
  0% { transform: translateY(0) scale(1); box-shadow: 0 8px 24px rgba(16,24,40,0.12); }
  40% { transform: translateY(-8px) scale(1.04); box-shadow: 0 20px 40px rgba(16,24,40,0.18); }
  60% { transform: translateY(-6px) scale(1.03); box-shadow: 0 16px 36px rgba(16,24,40,0.16); }
  100% { transform: translateY(0) scale(1); box-shadow: 0 8px 24px rgba(16,24,40,0.12); }
`;

const ServerCheck: React.FC<PropsWithChildren<{}>> = ({ children }) => {
  const [checking, setChecking] = useState(true);
  const [available, setAvailable] = useState<boolean | null>(null);
  const theme = useTheme();

  const check = async () => {
    setChecking(true);
    try {
      await api.get('');
      setAvailable(true);
    } catch (err: any) {
      if (!err || !err.response) setAvailable(false);
      else setAvailable(true);
    } finally {
      setChecking(false);
    }
  };

  useEffect(() => {
    check();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  if (checking) {
    return (
      <Box sx={{ minHeight: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <CircularProgress />
      </Box>
    );
  }

  if (available === false) {
    return (
      <Box
        sx={{
          minHeight: '100vh',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          bgcolor: 'background.default',
          position: 'relative',
          overflow: 'hidden',
          p: 2,
        }}
      >
        {/* decorative background text (Dostoevsky) - bundled fallback */}
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
            color: theme.palette.mode === 'dark' ? 'rgba(255,255,255,0.12)' : 'rgba(0,0,0,0.12)',
            textAlign: 'left',
            opacity: 1,
          }}
        >
          {Array(30).fill((staticDostoevsky || '«Бесы» — Ф. М. Достоевский').trim()).join('\n\n')}
        </Box>
        <Box
          sx={{
            position: 'absolute',
            width: 420,
            height: 420,
            borderRadius: '50%',
            background: `radial-gradient(circle at 30% 30%, color-mix(in srgb, var(--md-sys-color-primary, #7c3aed) 18%, transparent), color-mix(in srgb, var(--md-sys-color-primary, #7c3aed) 6% , transparent) 40%, transparent 60%)`,
            top: -120,
            left: -120,
            transform: 'rotate(20deg)',
            animation: `${float} 6s ease-in-out infinite`,
            filter: 'blur(36px)',
            opacity: 0.9,
          }}
        />
        <Box
          sx={{
            position: 'absolute',
            width: 360,
            height: 360,
            borderRadius: '50%',
            background: `radial-gradient(circle at 70% 70%, color-mix(in srgb, var(--md-sys-color-secondary, #10b981) 12%, transparent), color-mix(in srgb, var(--md-sys-color-secondary, #10b981) 4% , transparent) 40%, transparent 60%)`,
            bottom: -100,
            right: -100,
            transform: 'rotate(-10deg)',
            animation: `${float} 7s ease-in-out infinite`,
            filter: 'blur(36px)',
            opacity: 0.95,
          }}
        />

        <Box sx={{ textAlign: 'center', zIndex: 2, width: '100%' }}>
          <Box
            component="h1"
            sx={{
              m: 0,
              fontWeight: 800,
              letterSpacing: '-0.02em',
              lineHeight: 1,
              fontSize: 'clamp(96px, 24vw, 220px)',
              background: 'linear-gradient(90deg, var(--md-sys-color-primary, #7c3aed) 0%, var(--md-sys-color-surface-tint, #06b6d4) 50%, var(--md-sys-color-secondary, #34d399) 100%)',
              WebkitBackgroundClip: 'text',
              backgroundClip: 'text',
              color: 'transparent',
              backgroundSize: '200% 200%',
              animation: `${gradientShift} 6s ease infinite, ${glow} 3s ease-in-out infinite`,
              textShadow: '0 8px 24px rgba(15,23,42,0.18)',
              display: 'inline-block',
            }}
          >
            404
          </Box>

          <Box sx={{ mt: 6 }}>
            <Button
              variant="contained"
              color="primary"
              onClick={check}
              sx={{
                px: 5,
                py: 1.8,
                fontSize: '16px',
                textTransform: 'uppercase',
                borderRadius: 3,
                // Sync button animation with header glow/gradient (use CSS tokens)
                background: 'linear-gradient(90deg, var(--md-sys-color-primary, #7c3aed) 0%, var(--md-sys-color-surface-tint, #06b6d4) 50%, var(--md-sys-color-secondary, #34d399) 100%)',
                color: 'white',
                backgroundSize: '200% 200%',
                animation: `${buttonBeat} 3s ease-in-out infinite, ${gradientShift} 6s linear infinite`,
                boxShadow: '0 8px 24px rgba(16,24,40,0.12)',
                transformOrigin: 'center',
                transition: 'transform 220ms ease, box-shadow 220ms ease',
                '&:hover': {
                  transform: 'translateY(-6px) scale(1.04)',
                },
              }}
            >
              Повторить
            </Button>
          </Box>
        </Box>
      </Box>
    );
  }

  return <>{children}</>;
};

export default ServerCheck;
