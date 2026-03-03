import React from 'react';
import { Routes, Route } from 'react-router-dom';
import Layout from './Layout';
import HomePage from './HomePage';
import AboutPage from './AboutPage';
import ProfilePage from './ProfilePage';
import NotFoundPage from './NotFoundPage';
import LoginPage from './LoginPage';
import RegisterPage from './RegisterPage';
import ForgotPasswordPage from './ForgotPasswordPage';
import CreateBookPage from './CreateBookPage';
import BookDetailPage from './BookDetailPage';
import BooksPage from './BooksPage';
import SearchPage from './SearchPage';
import MyBooksPage from './MyBooksPage';
import ReservedBooksPage from './ReservedBooksPage';
import MyShelfPage from './MyShelfPage';
import ReturnBookPage from './ReturnBookPage';
import ReadBooksPage from './ReadBooksPage';
import BookStatsPage from './BookStatsPage';
import ShelfEditPage from './ShelfEditPage';
import AdminPage from './AdminPage';
import ModerPage from './ModerPage';
import PrivateRoute from './PrivateRoute';
import RoleRoute from './RoleRoute';
// TO_BE_REMOVED_BEFORE_RELEASE
import ColorPalettePage from './ColorPalettePage';

const App: React.FC = () => {
  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="create" element={<PrivateRoute><CreateBookPage /></PrivateRoute>} />
        <Route path="books" element={<PrivateRoute><BooksPage /></PrivateRoute>}>
          <Route path="my" element={<MyBooksPage />} />
          <Route path="my/stats/:bookId" element={<BookStatsPage />} />
          <Route path="reserved" element={<ReservedBooksPage />} />
          <Route path="shelf" element={<MyShelfPage />} />
          <Route path="shelf/:bookId/edit" element={<ShelfEditPage />} />
          <Route path="read" element={<ReadBooksPage />} />
          <Route path="return/:bookId" element={<ReturnBookPage />} />
        </Route>
        <Route path="about" element={<AboutPage />} />
        <Route path="book/:bookId" element={<BookDetailPage />} />
        <Route path="profile" element={<PrivateRoute><ProfilePage /></PrivateRoute>} />
        <Route path="search" element={<SearchPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />
        <Route path="forgot-password" element={<ForgotPasswordPage />} />

        {/* Защищённые маршруты по роли */}
        <Route path="moder" element={<RoleRoute requiredRole="moder"><ModerPage /></RoleRoute>} />
        <Route path="admin" element={<RoleRoute requiredRole="admin"><AdminPage /></RoleRoute>} />

        <Route path="*" element={<NotFoundPage />} />
        {/* TO_BE_REMOVED_BEFORE_RELEASE */}
        <Route path="color-palette" element={<ColorPalettePage />} />
      </Route>
    </Routes>
  );
};

export default App;
