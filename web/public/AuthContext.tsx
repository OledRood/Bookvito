import React, { createContext, useState, useContext, ReactNode, useEffect, useCallback, useMemo } from 'react';
import { User } from './user';
import userService from '../src/services/userService';

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isAdmin: boolean;
  isModer: boolean;
  hasRole: (role: string) => boolean;
  login: (email: string, password: string) => Promise<boolean>;
  register: (name: string, email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  saveTokens: (accessToken: string, refreshToken: string, userData: User) => void;
  clearTokens: () => void;
  setLoading: (loading: boolean) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);
  const [loading, setLoading] = useState<boolean>(true); // Для обработки начальной загрузки и проверки токена

  // При монтировании пытаемся восстановить пользователя из localStorage
  // и подписываемся на событие принудительного выхода из interceptor
  useEffect(() => {
    try {
      const stored = localStorage.getItem('user');
      const token = localStorage.getItem('accessToken') || localStorage.getItem('token');
      if (stored && token) {
        const parsed = JSON.parse(stored);
        setUser(parsed as User);
        setIsAuthenticated(true);
      }
    } catch (e) {
      console.warn('Failed to restore auth from localStorage', e);
    } finally {
      setLoading(false);
    }

    const handleUnauthorized = () => {
      localStorage.removeItem('accessToken');
      localStorage.removeItem('token');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('user');
      setUser(null);
      setIsAuthenticated(false);
    };
    window.addEventListener('auth:unauthorized', handleUnauthorized);
    return () => window.removeEventListener('auth:unauthorized', handleUnauthorized);
  }, []);

  const saveTokens = (accessToken: string, refreshToken: string, userData: User) => {
    // save both keys for compatibility across the app
    localStorage.setItem('accessToken', accessToken);
    localStorage.setItem('token', accessToken);
    if (refreshToken) localStorage.setItem('refreshToken', refreshToken);
    localStorage.setItem('user', JSON.stringify(userData));
    setUser(userData);
    setIsAuthenticated(true);
  };

  const clearTokens = () => {
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('user');
    setUser(null);
    setIsAuthenticated(false);
  };

  const login = useCallback(async (email: string, password: string): Promise<boolean> => {
    try {
      const data: any = await userService.login(email, password);
      // backend returns { token, user } per spec, or { accessToken }
      const t = data.token || data.accessToken;
      if (t) {
        // Save tokens first so subsequent requests (like fetching profile) include Authorization header
        try {
          localStorage.setItem('accessToken', t);
          localStorage.setItem('token', t);
          if (data.refreshToken) localStorage.setItem('refreshToken', data.refreshToken || '');
        } catch (e) {}

        // Try to get full user profile from API. If it fails, fall back to any user object returned by login (may be empty).
        let fetchedUser = data.user || {};
        try {
          const profile = await userService.getProfile();
          if (profile) fetchedUser = profile;
        } catch (e) {
          // ignore - keep fetchedUser as-is
        }

        saveTokens(t, data.refreshToken || '', fetchedUser);
        return true;
      }
      return false;
    } catch (error: any) {
      console.error('Login failed:', error);
      return false;
    }
  }, []);

  const register = useCallback(async (name: string, email: string, password: string): Promise<{ success: boolean; error?: string }> => {
    try {
      const res: any = await userService.register(name, email, password);
      if (res && res.success) {
        const data = res.data || {};
        const t = data.token || data.accessToken;
        if (t) {
          // Save tokens and then attempt to fetch full profile (backend may not return user in register response)
          try {
            localStorage.setItem('accessToken', t);
            localStorage.setItem('token', t);
            if (data.refreshToken) localStorage.setItem('refreshToken', data.refreshToken || '');
          } catch (e) {}

          let fetchedUser = data.user || {};
          try {
            const profile = await userService.getProfile();
            if (profile) fetchedUser = profile;
          } catch (e) {}

          saveTokens(t, data.refreshToken || '', fetchedUser);
          return { success: true };
        }
        // success but no token? treat as failure
        return { success: false, error: 'No token returned from server' };
      }
      const errMsg = (res && res.error) || 'Registration failed';
      console.error('Registration failed:', errMsg);
      return { success: false, error: errMsg };
    } catch (error: any) {
      const message = (error && error.message) || 'Registration failed';
      console.error('Registration failed (exception):', error);
      return { success: false, error: message };
    }
  }, []);

  const logout = useCallback(() => {
    userService.logout(); // revoke refresh token on server (fire-and-forget)
    clearTokens();
  }, []);

  const isAdmin = user?.role === 'admin';
  const isModer = user?.role === 'moder' || user?.role === 'admin';
  const hasRole = (role: string) => user?.role === role || (role === 'moder' && user?.role === 'admin');

  const value = useMemo(
    () => ({
      user,
      isAuthenticated,
      isAdmin,
      isModer,
      hasRole,
      login,
      register,
      logout,
      saveTokens,
      clearTokens,
      setLoading,
    }),
    [user, isAuthenticated, isAdmin, isModer, login, register, logout],
  );

  if (loading) {
    // Опционально: отобразить спиннер загрузки во время проверки статуса авторизации
    return <div>Загрузка авторизации...</div>;
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};