import React, { useEffect, useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  Select,
  MenuItem,
  FormControl,
  CircularProgress,
  Divider,
  Chip,
} from '@mui/material';
import PeopleIcon from '@mui/icons-material/People';
import MenuBookIcon from '@mui/icons-material/MenuBook';
import SwapHorizIcon from '@mui/icons-material/SwapHoriz';
import ReportIcon from '@mui/icons-material/Report';
import adminService from '../src/services/adminService';
import { useNotification } from '../contexts/NotificationContext';

const roleColor: Record<string, 'default' | 'primary' | 'error' | 'warning'> = {
  user: 'default',
  moder: 'primary',
  admin: 'error',
};

const roleLabel: Record<string, string> = {
  user: 'Пользователь',
  moder: 'Модератор',
  admin: 'Администратор',
};

const AdminPage: React.FC = () => {
  const { showNotification } = useNotification();
  const [stats, setStats] = useState<any>(null);
  const [users, setUsers] = useState<any[]>([]);
  const [loadingStats, setLoadingStats] = useState(true);
  const [loadingUsers, setLoadingUsers] = useState(true);
  const [updatingRole, setUpdatingRole] = useState<string | null>(null);

  useEffect(() => {
    let mounted = true;
    adminService.getAdminStats().then(data => {
      if (mounted) { setStats(data); setLoadingStats(false); }
    }).catch(e => {
      if (mounted) { showNotification(e?.response?.data?.error || 'Ошибка при загрузке статистики', 'error'); setLoadingStats(false); }
    });
    adminService.getAdminUsers().then(data => {
      if (mounted) { setUsers(data); setLoadingUsers(false); }
    }).catch(e => {
      if (mounted) { showNotification(e?.response?.data?.error || 'Ошибка при загрузке пользователей', 'error'); setLoadingUsers(false); }
    });
    return () => { mounted = false; };
  }, []);

  const handleRoleChange = async (userId: string, newRole: string) => {
    setUpdatingRole(userId);
    try {
      await adminService.updateUserRole(userId, newRole);
      setUsers(prev => prev.map(u => u.id === userId ? { ...u, role: newRole } : u));
      showNotification('Роль обновлена', 'success');
    } catch (e: any) {
      showNotification(e?.response?.data?.error || 'Ошибка при обновлении роли', 'error');
    } finally {
      setUpdatingRole(null);
    }
  };

  return (
    <Box sx={{ maxWidth: 800, mx: 'auto', p: 2 }}>
      <Typography variant="h5" gutterBottom fontWeight={600}>
        Панель администратора
      </Typography>

      {/* Статистика */}
      <Typography variant="subtitle1" color="text.secondary" gutterBottom>
        Статистика платформы
      </Typography>
      {loadingStats ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Grid container spacing={2} sx={{ mb: 4 }}>
          {[
            { label: 'Пользователи', value: stats?.total_users ?? 0, icon: <PeopleIcon /> },
            { label: 'Книги', value: stats?.total_books ?? 0, icon: <MenuBookIcon /> },
            { label: 'Активные обмены', value: stats?.active_exchanges ?? 0, icon: <SwapHorizIcon /> },
            { label: 'Жалобы', value: stats?.pending_reports ?? 0, icon: <ReportIcon /> },
          ].map(item => (
            <Grid item xs={6} sm={3} key={item.label}>
              <Paper
                elevation={0}
                sx={{
                  p: 2,
                  borderRadius: 3,
                  textAlign: 'center',
                  bgcolor: 'var(--md-sys-color-surface-container-low)',
                }}
              >
                <Box sx={{ color: 'var(--md-sys-color-primary)', mb: 0.5 }}>{item.icon}</Box>
                <Typography variant="h5" fontWeight={700}>{item.value}</Typography>
                <Typography variant="caption" color="text.secondary">{item.label}</Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      )}

      <Divider sx={{ my: 2 }} />

      {/* Управление пользователями */}
      <Typography variant="subtitle1" color="text.secondary" gutterBottom>
        Управление пользователями
      </Typography>
      {loadingUsers ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 3 }}>
          <CircularProgress />
        </Box>
      ) : (
        <Paper sx={{ borderRadius: 3, overflow: 'hidden' }}>
          <List disablePadding>
            {users.map((u: any, idx: number) => (
              <ListItem
                key={u.id}
                divider={idx < users.length - 1}
                sx={{ py: 1.5 }}
                secondaryAction={
                  <FormControl size="small" sx={{ minWidth: 140 }}>
                    <Select
                      value={u.role || 'user'}
                      onChange={e => handleRoleChange(u.id, e.target.value)}
                      disabled={updatingRole === u.id}
                      sx={{ borderRadius: 2 }}
                    >
                      <MenuItem value="user">Пользователь</MenuItem>
                      <MenuItem value="moder">Модератор</MenuItem>
                      <MenuItem value="admin">Администратор</MenuItem>
                    </Select>
                  </FormControl>
                }
              >
                <ListItemAvatar>
                  <Avatar src={u.avatar ? `/images/${u.avatar}` : undefined} alt={u.name}>
                    {u.name?.[0]?.toUpperCase()}
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary={
                    <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
                      <Typography variant="body2" fontWeight={500}>{u.name}</Typography>
                      <Chip
                        label={roleLabel[u.role] || u.role}
                        color={roleColor[u.role] || 'default'}
                        size="small"
                        sx={{ height: 20, fontSize: '0.65rem' }}
                      />
                    </Box>
                  }
                  secondary={u.email}
                  primaryTypographyProps={{ color: 'text.primary' }}
                  secondaryTypographyProps={{ color: 'text.secondary', fontSize: '0.75rem' }}
                />
              </ListItem>
            ))}
          </List>
        </Paper>
      )}
    </Box>
  );
};

export default AdminPage;
