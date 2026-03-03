import api from './api';

export const getReports = async () => {
  const resp = await api.get('moder/reports');
  return resp.data || [];
};

export const resolveReport = async (reportId) => {
  const resp = await api.put(`moder/reports/${reportId}/resolve`);
  return resp.data;
};

export const dismissReport = async (reportId) => {
  const resp = await api.put(`moder/reports/${reportId}/dismiss`);
  return resp.data;
};

export const archiveBook = async (bookId) => {
  const resp = await api.put(`moder/books/${bookId}/archive`);
  return resp.data;
};

export const reportBook = async (bookId, reason) => {
  const resp = await api.post(`books/${bookId}/report`, { reason });
  return resp.data;
};

export default { getReports, resolveReport, dismissReport, archiveBook, reportBook };
