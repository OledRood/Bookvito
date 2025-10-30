import React, { ReactNode } from 'react';
import { AuthProvider } from './AuthContext';

export const AuthProviderWithManager: React.FC<{ children: ReactNode }> = ({ children }) => (
  <AuthProvider>
    {children}
  </AuthProvider>
);