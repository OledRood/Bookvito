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

export const updateBook = async (bookId, payload) => {
  const resp = await api.put(`books/${bookId}`, payload);
  return resp.data;
};

export const returnBook = async (payload) => {
  const body = {
    book_id: payload.bookId,
    title: payload.title,
    author: payload.author,
  };
  if (payload.condition) body.condition = payload.condition;
  if (payload.description !== undefined) body.description = payload.description;
  if (payload.currentLocationId !== undefined) body.current_location_id = payload.currentLocationId || null;
  if (payload.imageUrl !== undefined) body.image_url = payload.imageUrl;

  const resp = await api.put('books/return', body);
  return resp.data;
};

export const uploadBookImage = async (file) => {
  const form = new FormData();
  form.append('image', file);
  const resp = await api.post('books/upload', form, {
    headers: { 'Content-Type': 'multipart/form-data' },
  });
  return resp.data || {};
};

export const setBookImage = async (bookId, imageUrl) => {
  const resp = await api.put(`books/image/${bookId}`, { image_url: imageUrl });
  return resp.data;
};

export const autoFillBook = async (query) => {
  const resp = await api.get('books/auto-fill', { params: { q: query } });
  return resp.data || null;
};

export const deleteBookImage = async (bookId) => {
  const resp = await api.delete(`books/image/${bookId}`);
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
  updateBook,
  returnBook,
  uploadBookImage,
  setBookImage,
  autoFillBook,
  deleteBookImage,
};
