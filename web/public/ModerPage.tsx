import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  List,
  ListItem,
  ListItemText,
  CircularProgress,
  Button,
  Chip,
  Divider,
  Stack,
} from '@mui/material';
import ReportGmailerrorredIcon from '@mui/icons-material/ReportGmailerrorred';
import CheckCircleOutlineIcon from '@mui/icons-material/CheckCircleOutline';
import ArchiveIcon from '@mui/icons-material/Archive';
import moderService from '../src/services/moderService';
import { useNotification } from '../contexts/NotificationContext';

const ModerPage: React.FC = () => {
  const { showNotification } = useNotification();
  const [reports, setReports] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [processing, setProcessing] = useState<string | null>(null);

  const load = async () => {
    setLoading(true);
    try {
      const data = await moderService.getReports();
      setReports(data || []);
    } catch (e: any) {
      showNotification(e?.response?.data?.error || 'Ошибка при загрузке жалоб', 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handle = async (action: () => Promise<any>, successMsg: string, reportId: string) => {
    setProcessing(reportId);
    try {
      await action();
      showNotification(successMsg, 'success');
      // Убираем жалобу из списка после обработки
      setReports(prev => prev.filter(r => r.id !== reportId));
    } catch (e: any) {
      showNotification(e?.response?.data?.error || 'Ошибка', 'error');
    } finally {
      setProcessing(null);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', p: 2 }}>
      <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 2 }}>
        <ReportGmailerrorredIcon color="warning" />
        <Typography variant="h5" fontWeight={600}>
          Жалобы на объявления
        </Typography>
        {!loading && (
          <Chip
            label={reports.length}
            color={reports.length > 0 ? 'warning' : 'default'}
            size="small"
            sx={{ ml: 1 }}
          />
        )}
      </Box>

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 4 }}>
          <CircularProgress />
        </Box>
      ) : reports.length === 0 ? (
        <Paper
          elevation={0}
          sx={{
            p: 4,
            borderRadius: 3,
            textAlign: 'center',
            bgcolor: 'var(--md-sys-color-surface-container-low)',
          }}
        >
          <CheckCircleOutlineIcon sx={{ fontSize: 48, color: 'success.main', mb: 1 }} />
          <Typography variant="body1" color="text.secondary">
            Нет активных жалоб
          </Typography>
        </Paper>
      ) : (
        <Paper sx={{ borderRadius: 3, overflow: 'hidden' }}>
          <List disablePadding>
            {reports.map((report: any, idx: number) => (
              <React.Fragment key={report.id}>
                {idx > 0 && <Divider />}
                <ListItem
                  alignItems="flex-start"
                  sx={{ flexDirection: 'column', py: 2, px: 2 }}
                >
                  {/* Заголовок книги */}
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 1, mb: 0.5, width: '100%' }}>
                    <Typography variant="subtitle2" fontWeight={600}>
                      {report.book?.title || 'Книга не найдена'}
                    </Typography>
                    {report.book?.author && (
                      <Typography variant="caption" color="text.secondary">
                        — {report.book.author}
                      </Typography>
                    )}
                    <Chip
                      label={`ID: ${report.book?.id?.slice(0, 8)}...`}
                      size="small"
                      sx={{ ml: 'auto', fontSize: '0.65rem' }}
                    />
                  </Box>

                  {/* Автор жалобы */}
                  <Typography variant="caption" color="text.secondary" sx={{ mb: 0.5 }}>
                    От: {report.user?.name || 'Неизвестный пользователь'} ({report.user?.email})
                  </Typography>

                  {/* Текст жалобы */}
                  <Paper
                    elevation={0}
                    sx={{
                      p: 1.5,
                      borderRadius: 2,
                      bgcolor: 'var(--md-sys-color-surface-container)',
                      width: '100%',
                      mb: 1.5,
                    }}
                  >
                    <Typography variant="body2">{report.reason}</Typography>
                  </Paper>

                  {/* Дата */}
                  <Typography variant="caption" color="text.secondary" sx={{ mb: 1.5 }}>
                    {new Date(report.created_at).toLocaleString('ru-RU')}
                  </Typography>

                  {/* Действия */}
                  <Stack direction="row" spacing={1}>
                    <Button
                      size="small"
                      variant="contained"
                      color="error"
                      startIcon={<ArchiveIcon />}
                      disabled={processing === report.id}
                      onClick={() => handle(
                        () => moderService.archiveBook(report.book?.id),
                        'Книга архивирована',
                        report.id
                      )}
                    >
                      Архивировать
                    </Button>
                    <Button
                      size="small"
                      variant="outlined"
                      color="success"
                      disabled={processing === report.id}
                      onClick={() => handle(
                        () => moderService.resolveReport(report.id),
                        'Жалоба закрыта',
                        report.id
                      )}
                    >
                      Закрыть
                    </Button>
                    <Button
                      size="small"
                      variant="text"
                      color="inherit"
                      disabled={processing === report.id}
                      onClick={() => handle(
                        () => moderService.dismissReport(report.id),
                        'Жалоба отклонена',
                        report.id
                      )}
                    >
                      Отклонить
                    </Button>
                  </Stack>
                </ListItem>
              </React.Fragment>
            ))}
          </List>
        </Paper>
      )}
    </Box>
  );
};

export default ModerPage;
