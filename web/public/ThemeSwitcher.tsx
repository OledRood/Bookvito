import React from 'react';
import { IconButton, useTheme } from '@mui/material';
import LightModeIcon from '@mui/icons-material/LightMode';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import { useThemeContext } from './ThemeContext';

interface ThemeSwitcherProps {
  isDark?: boolean;
  onClick?: () => void;
  isDarkMode?: boolean;
  toggleTheme?: () => void;
}

const ThemeSwitcher: React.FC<ThemeSwitcherProps> = (props) => {
  const { isDark, onClick, isDarkMode, toggleTheme } = props;
  const theme = useTheme();
  const ctx = useThemeContext();

  const effectiveIsDark =
    typeof isDark === 'boolean'
      ? isDark
      : typeof isDarkMode === 'boolean'
        ? isDarkMode
        : ctx.mode === 'dark';

  const handleClick = onClick ?? toggleTheme ?? ctx.toggleTheme;

  return (
    <IconButton
      onClick={handleClick}
      sx={{
        border: `2px solid ${theme.palette.divider}`,
        bgcolor: 'transparent',
        borderRadius: '50%',
        width: 40,
        height: 40,
        color: theme.palette.text.primary,
        overflow: 'hidden',
        position: 'relative',
        alignItems: 'center',
        justifyContent: 'center',
        display: 'flex',
        '&:hover': {
          bgcolor: 'transparent',
          borderColor: theme.palette.primary.main,
        },
      }}
    >
      {/* Dark Icon */}
      <span
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          position: 'absolute',
          left: 0,
          right: 0,
          top: 0,
          bottom: 0,
          transition: 'transform 0.35s cubic-bezier(.65,.05,.36,1), opacity 0.2s',
          transform: `translateY(${effectiveIsDark ? 0 : -32}px)`,
          opacity: effectiveIsDark ? 1 : 0,
          pointerEvents: effectiveIsDark ? 'auto' : 'none',
          width: '100%',
          height: '100%'
        }}
      >
        <DarkModeIcon sx={{ fontSize: 24 }} />
      </span>
      {/* Light Icon */}
      <span
        style={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          position: 'absolute',
          left: 0,
          right: 0,
          top: 0,
          bottom: 0,
          transition: 'transform 0.35s cubic-bezier(.65,.05,.36,1), opacity 0.2s',
          transform: `translateY(${effectiveIsDark ? 32 : 0}px)`,
          opacity: effectiveIsDark ? 0 : 1,
          pointerEvents: effectiveIsDark ? 'none' : 'auto',
          width: '100%',
          height: '100%'
        }}
      >
        <LightModeIcon sx={{ fontSize: 24 }} />
      </span>
    </IconButton>

  );
};

export default ThemeSwitcher;
