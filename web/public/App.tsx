import React, { Suspense, lazy } from 'react';
import { Routes, Route, Navigate } from 'react-router-dom';
import { Box, CircularProgress } from '@mui/material';
import Layout from './Layout';
import PrivateRoute from './PrivateRoute';
import RoleRoute from './RoleRoute';
// TO_BE_REMOVED_BEFORE_RELEASE
const HomePage = lazy(() => import('./HomePage'));
const AboutPage = lazy(() => import('./AboutPage'));
const ProfilePage = lazy(() => import('./ProfilePage'));
const NotFoundPage = lazy(() => import('./NotFoundPage'));
const LoginPage = lazy(() => import('./LoginPage'));
const RegisterPage = lazy(() => import('./RegisterPage'));
const ForgotPasswordPage = lazy(() => import('./ForgotPasswordPage'));
const CreateBookPage = lazy(() => import('./CreateBookPage'));
const BookDetailPage = lazy(() => import('./BookDetailPage'));
const BooksPage = lazy(() => import('./BooksPage'));
const SearchPage = lazy(() => import('./SearchPage'));
const MyBooksPage = lazy(() => import('./MyBooksPage'));
const ReservedBooksPage = lazy(() => import('./ReservedBooksPage'));
const MyShelfPage = lazy(() => import('./MyShelfPage'));
const ReturnBookPage = lazy(() => import('./ReturnBookPage'));
const ReadBooksPage = lazy(() => import('./ReadBooksPage'));
const BookStatsPage = lazy(() => import('./BookStatsPage'));
const ShelfEditPage = lazy(() => import('./ShelfEditPage'));
const AdminPage = lazy(() => import('./AdminPage'));
const ModerPage = lazy(() => import('./ModerPage'));
const LegacyBookRedirect = lazy(() => import('./LegacyBookRedirect'));
const ColorPalettePage = lazy(() => import('./ColorPalettePage'));

const RouteFallback: React.FC = () => (
  <Box sx={{ minHeight: '40vh', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
    <CircularProgress />
  </Box>
);

const withSuspense = (node: React.ReactNode) => (
  <Suspense fallback={<RouteFallback />}>
    {node}
  </Suspense>
);

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={withSuspense(<HomePage />)} />
        <Route path="home" element={<Navigate to="/" replace />} />
        <Route path="create" element={<Navigate to="/books/new" replace />} />
        <Route path="books/new" element={<PrivateRoute>{withSuspense(<CreateBookPage />)}</PrivateRoute>} />
        <Route path="books" element={<PrivateRoute>{withSuspense(<BooksPage />)}</PrivateRoute>}>
          <Route path="my" element={withSuspense(<MyBooksPage />)} />
          <Route path="my/stats/:bookId" element={withSuspense(<BookStatsPage />)} />
          <Route path="reserved" element={withSuspense(<ReservedBooksPage />)} />
          <Route path="shelf" element={withSuspense(<MyShelfPage />)} />
          <Route path="shelf/:bookId/edit" element={withSuspense(<ShelfEditPage />)} />
          <Route path="read" element={withSuspense(<ReadBooksPage />)} />
          <Route path="return/:bookId" element={withSuspense(<ReturnBookPage />)} />
        </Route>
        <Route path="about" element={withSuspense(<AboutPage />)} />
        <Route path="book/:bookId" element={withSuspense(<LegacyBookRedirect />)} />
        <Route path="books/:bookIdSlug" element={withSuspense(<BookDetailPage />)} />
        <Route path="profile" element={<PrivateRoute>{withSuspense(<ProfilePage />)}</PrivateRoute>} />
        <Route path="search" element={withSuspense(<SearchPage />)} />
        <Route path="login" element={withSuspense(<LoginPage />)} />
        <Route path="register" element={withSuspense(<RegisterPage />)} />
        <Route path="forgot-password" element={withSuspense(<ForgotPasswordPage />)} />

        {/* Защищённые маршруты по роли */}
        <Route path="moder" element={<Navigate to="/moderation" replace />} />
        <Route path="moderation" element={<RoleRoute requiredRole="moder">{withSuspense(<ModerPage />)}</RoleRoute>} />
        <Route path="admin" element={<RoleRoute requiredRole="admin">{withSuspense(<AdminPage />)}</RoleRoute>} />

        <Route path="*" element={withSuspense(<NotFoundPage />)} />
        {/* TO_BE_REMOVED_BEFORE_RELEASE */}
        <Route path="color-palette" element={withSuspense(<ColorPalettePage />)} />
      </Route>
    </Routes>
  );
};

export default App;
