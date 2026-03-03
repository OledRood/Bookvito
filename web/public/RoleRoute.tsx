import React from 'react';
import { Navigate } from 'react-router-dom';
import { Box, Typography } from '@mui/material';
import { useAuth } from './AuthContext';
import { UserRole } from './user';

interface RoleRouteProps {
  children: React.ReactElement;
  /** Минимальная роль: 'moder' — для moder + admin; 'admin' — только для admin */
  requiredRole: UserRole;
}

/**
 * RoleRoute — показывает страницу только пользователям с нужной ролью.
 * Если не аутентифицирован — редирект на /login.
 * Если аутентифицирован, но роль не подходит — 403 заглушка.
 */
const RoleRoute: React.FC<RoleRouteProps> = ({ children, requiredRole }) => {
  const { isAuthenticated, user } = useAuth();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  const role = user?.role ?? 'user';

  // admin может всё; moder видит moder-маршруты
  const allowed =
    role === 'admin' ||
    (requiredRole === 'moder' && role === 'moder');

  if (!allowed) {
    return (
      <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', justifyContent: 'center', minHeight: '60vh', gap: 2 }}>
        <Typography variant="h3" color="error">403</Typography>
        <Typography variant="h6" color="text.secondary">
          Недостаточно прав для просмотра этой страницы
        </Typography>
      </Box>
    );
  }

  return children;
};

export default RoleRoute;
