import React from 'react';
import { Navigate, useParams } from 'react-router-dom';

const LegacyBookRedirect: React.FC = () => {
  const { bookId } = useParams<{ bookId: string }>();

  if (!bookId) {
    return <Navigate to="/" replace />;
  }

  return <Navigate to={`/books/${bookId}`} replace />;
};

export default LegacyBookRedirect;
