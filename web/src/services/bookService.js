import api from './api';

export const getMyBooks = async () => {
  const resp = await api.get('books/my');
  return resp.data || [];
};

export const getReservedBooks = async () => {
  const resp = await api.get('books/reserved');
  return resp.data || [];
};

export const extendReservation = async (bookId) => {
  const resp = await api.put('books/reservation/extend', { book_id: bookId });
  return resp.data;
};

export const cancelReservation = async (bookId) => {
  // axios delete with body requires data in config
  const resp = await api.delete('books/reservation/cancel', { data: { book_id: bookId } });
  return resp.data;
};

export const getShelfBooks = async () => {
  const resp = await api.get('books/shelf');
  return resp.data || [];
};

export const getReadBooks = async () => {
  const resp = await api.get('books/read');
  return resp.data || [];
};

export const getMyBooksStats = async () => {
  const resp = await api.get('books/my/stats');
  return resp.data || {};
};

export const getBookStats = async (bookId) => {
  const resp = await api.get(`books/my/stats/${bookId}`);
  return resp.data || {};
};

export const getBook = async (bookId) => {
  const resp = await api.get(`books/${bookId}`);
  return resp.data || null;
};

export const borrowBook = async (bookId) => {
  const resp = await api.put('books/borrow', { book_id: bookId });
  return resp.data;
};

export default {
  getMyBooks,
  getReservedBooks,
  getShelfBooks,
  getReadBooks,
  extendReservation,
  cancelReservation,
  getMyBooksStats,
  getBookStats,
  getBook,
  borrowBook,
};
