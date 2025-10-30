import { createTheme } from '@mui/material/styles';
import { red } from '@mui/material/colors';

// Small helper to expose common design tokens (radii, durations, shadows)
const tokens = {
  radii: {
    small: 8,
    medium: 12,
    large: 16,
  },
  motion: {
    short: 200,
    normal: 240,
  },
  shadow: {
    dp1: '0 1px 2px rgba(0,0,0,0.08)',
    dp2: '0 2px 6px rgba(0,0,0,0.12)',
    dp4: '0 8px 16px rgba(0,0,0,0.12)',
  },
  fontFamily: `Inter, Roboto, Helvetica, Arial, sans-serif`,
  // prefer Roboto as primary font
  fontFamilyPrimary: `Roboto, Helvetica, Arial, sans-serif`,
};

// Factory to create a light theme with optional primary seed
export const createLightTheme = (primarySeed = '#556cd6', sysColors: Record<string, string> = {}) =>
  createTheme({
    // Use concrete color tokens in the MUI palette so MUI color helpers can
    // compute tones (lighten/darken) correctly. Component backgrounds and
    // other places still use CSS variables (via `var(...)`) where appropriate.
    palette: ({
      mode: 'light',
      primary: { main: primarySeed },
      secondary: { main: '#19857b' },
      error: { main: '#b91c1c' },
      background: {
        default: '#fafafa',
        paper: '#ffffff',
      },
      // expose a neutral/info container color for fallbacks
      info: { main: '#e6f4f1' } as any,
      // attach runtime system CSS colors (e.g. md-sys-color-*) for consuming
      // components. Stored under `palette.sys` to avoid colliding with MUI
      // standard keys. This is intentionally `as any` to keep typing simple.
      sys: sysColors as any,
    } as any),
    typography: {
      fontFamily: tokens.fontFamilyPrimary,
      // normalize button typography so sizes don't jump between light/dark
      button: { textTransform: 'none', fontWeight: 600, fontSize: '0.875rem', lineHeight: '20px' },
    },
    shape: { borderRadius: tokens.radii.medium },
    components: {
      MuiButton: {
        defaultProps: { disableElevation: false },
        styleOverrides: {
          // enforce consistent button sizing and padding so theme switches don't change dimensions
          root: {
            borderRadius: tokens.radii.medium,
            transition: `all ${tokens.motion.short}ms ease`,
            minHeight: 34,
            padding: '6px 12px',
            fontSize: '0.875rem',
            lineHeight: '20px',
          },
        },
      },
      MuiCard: {
        styleOverrides: {
          root: { borderRadius: tokens.radii.medium, transition: `transform ${tokens.motion.short}ms ease, box-shadow ${tokens.motion.short}ms ease` },
        },
      },
    },
  });

export const createDarkTheme = (primarySeed = '#bb86fc', sysColors: Record<string, string> = {}) =>
  createTheme({
    palette: ({
      mode: 'dark',
      primary: { main: primarySeed },
      secondary: { main: '#03dac6' },
      error: { main: '#cf6679' },
      background: {
        default: '#121212',
        paper: '#1e1e1e',
      },
      info: { main: '#003731' } as any,
      sys: sysColors as any,
    } as any),
  typography: { fontFamily: tokens.fontFamilyPrimary, button: { textTransform: 'none', fontWeight: 600, fontSize: '0.875rem', lineHeight: '20px' } },
    shape: { borderRadius: tokens.radii.medium },
    components: {
      MuiButton: {
        styleOverrides: { root: { borderRadius: tokens.radii.medium, transition: `all ${tokens.motion.short}ms ease`, minHeight: 34, padding: '6px 12px', fontSize: '0.875rem', lineHeight: '20px' } },
      },
      MuiCard: {
        styleOverrides: { root: { borderRadius: tokens.radii.medium, transition: `transform ${tokens.motion.short}ms ease, box-shadow ${tokens.motion.short}ms ease` } },
      },
    },
  });

export const designTokens = tokens;