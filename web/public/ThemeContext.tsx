import React, { createContext, useState, useMemo, useContext, ReactNode, useEffect } from 'react';
import { ThemeProvider as MuiThemeProvider, CssBaseline } from '@mui/material';
import { createLightTheme, createDarkTheme, designTokens } from './theme';
// Import CSS theme token files so they are available to the app build.
// These files define the CSS variables under class selectors (e.g. .light-mc, .dark-mc)
// We will toggle those classes on the document element so the vars become active.
import '../css/light.css';
import '../css/dark.css';

interface ThemeContextType {
  toggleTheme: () => void;
  mode: 'light' | 'dark';
  setPrimarySeed: (hex: string) => void;
}

const ThemeContext = createContext<ThemeContextType | undefined>(undefined);

export const useThemeContext = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useThemeContext must be used within a ThemeProvider');
  }
  return context;
};

export const ThemeProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  // initialize mode from localStorage if available so user's choice persists
  const [mode, setMode] = useState<'light' | 'dark'>(() => {
    try {
      const stored = localStorage.getItem('bookvito-theme-mode');
      return stored === 'dark' ? 'dark' : 'light';
    } catch (e) {
      return 'light';
    }
  });

  const [primarySeed, setPrimarySeedState] = useState<string>(() => {
    try {
      const stored = localStorage.getItem('bookvito-primary-seed');
      return stored && stored.startsWith('#') ? stored : '#556cd6';
    } catch (e) {
      return '#556cd6';
    }
  });

  // Runtime-collected CSS system colors (e.g. --md-sys-color-surface-container)
  const [cssSysColors, setCssSysColors] = useState<Record<string, string>>({});

  const getCssSystemColors = (): Record<string, string> => {
    if (typeof window === 'undefined') return {};
    const styles = getComputedStyle(document.documentElement);
    const res: Record<string, string> = {};
    for (let i = 0; i < styles.length; i++) {
      const prop = styles.item(i);
      if (!prop) continue;
      if (prop.startsWith('--md-sys-color-')) {
        const value = styles.getPropertyValue(prop).trim();
        const key = prop.replace('--md-sys-color-', '');
        res[key] = value;
      }
    }
    return res;
  };

  const theme = useMemo(
    () => (mode === 'light' ? createLightTheme(primarySeed, cssSysColors) : createDarkTheme(primarySeed, cssSysColors)),
    [mode, primarySeed, cssSysColors],
  );
  const toggleTheme = () => {
    setMode((prevMode) => {
      const next = prevMode === 'light' ? 'dark' : 'light';
      try {
        localStorage.setItem('bookvito-theme-mode', next);
      } catch (e) {
        /* ignore */
      }
      return next;
    });
  };

  // Apply CSS variables for tokens and primary seed so non-MUI parts can use them
  useEffect(() => {
    const root = document.documentElement;
    root.style.setProperty('--primary-seed', primarySeed);
    root.style.setProperty('--radii-medium', `${designTokens.radii.medium}px`);
    root.style.setProperty('--motion-short', `${designTokens.motion.short}ms`);
    root.style.setProperty('--shadow-dp1', designTokens.shadow.dp1);
    root.style.setProperty('--shadow-dp2', designTokens.shadow.dp2);
    // transition for color/background changes
    root.style.setProperty('transition', `background var(--motion-short) ease, color var(--motion-short) ease`);
  }, [primarySeed]);

  // Ensure the CSS theme class is applied so the CSS variables defined under
  // `.light-mc` / `.dark-mc` become available on the page. Without this the
  // `var(--md-sys-color-...)` fallbacks will trigger and you'll see default
  // colors (green) instead of theme tokens.
  useEffect(() => {
    const docEl = document.documentElement;
  const classesToRemove = ['light', 'dark', 'light-mc', 'dark-mc', 'light-hc', 'dark-hc', 'light-medium-contrast', 'dark-medium-contrast', 'light-high-contrast', 'dark-high-contrast'];
    classesToRemove.forEach((c) => docEl.classList.remove(c));
  const classToAdd = mode === 'light' ? 'light' : 'dark';
    docEl.classList.add(classToAdd);
    // After we add the class, read system CSS variables and store them so
    // the MUI theme can include them (see createLightTheme/createDarkTheme).
    try {
      const sys = getCssSystemColors();
      setCssSysColors(sys);
    } catch (e) {
      // ignore
    }
    return () => {
      docEl.classList.remove(classToAdd);
    };
  }, [mode]);

  // After applying the CSS theme class, try to read the CSS token for the
  // primary color and use it to update the MUI primary seed. We keep a
  // fallback (initial state) in case the variable isn't present.
  useEffect(() => {
    const docEl = document.documentElement;

    const getCssVar = (name: string) => {
      try {
        return getComputedStyle(docEl).getPropertyValue(name).trim();
      } catch (e) {
        return '';
      }
    };

    const cssPrimary = getCssVar('--md-sys-color-primary');
    if (cssPrimary) {
      // cssPrimary is typically in the form "rgb(r g b)" or "rgb(r, g, b)"
      const toHex = (val: string) => {
        const rgbMatch = val.match(/\d{1,3}(?:[, ]+\d{1,3}){2}/);
        if (!rgbMatch) return '';
        const nums = rgbMatch[0].split(/[ ,]+/).map((n) => Number(n));
        if (nums.length < 3 || nums.some((n) => Number.isNaN(n))) return '';
        return (
          '#' +
          nums
            .slice(0, 3)
            .map((n) => n.toString(16).padStart(2, '0'))
            .join('')
        );
      };

      const hex = toHex(cssPrimary);
      if (hex) setPrimarySeedState(hex);
    }
  }, [mode]);

  const setPrimarySeed = (hex: string) => setPrimarySeedState(hex);

  // persist primary seed when it changes
  useEffect(() => {
    try {
      localStorage.setItem('bookvito-primary-seed', primarySeed);
    } catch (e) {
      // ignore
    }
  }, [primarySeed]);

  return (
    <ThemeContext.Provider value={{ toggleTheme, mode, setPrimarySeed }}>
      <MuiThemeProvider theme={theme}>
        <CssBaseline />
        {children}
      </MuiThemeProvider>
    </ThemeContext.Provider>
  );
};