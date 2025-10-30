import React, { useState } from 'react';
import {
  Box,
  Typography,
  Container,
  Paper,
  Button,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  TextField,
  RadioGroup,
  Radio,
  FormControl,
  FormLabel,
  FormControlLabel,
} from '@mui/material';
import { useAuth } from './AuthContext';
import userService from '../src/services/userService';
 
import { useNavigate } from 'react-router-dom';
import { useNotification } from '../contexts/NotificationContext';

const ProfilePage: React.FC = () => {
  const { showNotification } = useNotification();
  const { user, logout, saveTokens } = useAuth();
  // temporarily hide notifications item in profile UI (do not delete) — toggle this flag to show again
  const HIDE_PROFILE_NOTIFICATIONS = true;
  
  const navigate = useNavigate();

  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [deleteReason, setDeleteReason] = useState('');
  const [selectedReason, setSelectedReason] = useState('');
  const [unreadNotificationsCount, setUnreadNotificationsCount] = useState(3);
  const [avatarDialogOpen, setAvatarDialogOpen] = useState(false);
  const [availableAvatars, setAvailableAvatars] = useState<string[]>([]);
  const [selectedAvatarFile, setSelectedAvatarFile] = useState<string | null>(null);
  const [nameDialogOpen, setNameDialogOpen] = useState(false);
  const [nameInput, setNameInput] = useState<string>(user?.name || '');

  const handleNotificationsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNotificationsEnabled(event.target.checked);
    showNotification(`Уведомления ${event.target.checked ? 'включены' : 'отключены'}`, 'info');
    setUnreadNotificationsCount(0);
  };

  const handleDeleteAccount = () => setOpenDeleteDialog(true);
  const handleCloseDeleteDialog = () => {
    setOpenDeleteDialog(false);
    setDeleteReason('');
    setSelectedReason('');
  };

  const handleConfirmDelete = async () => {
    console.log('Deleting account for user:', user?.email);
    const reason = selectedReason === 'other' ? deleteReason : selectedReason;
    showNotification('Удаление аккаунта...', 'info');
    try {
      // use userService which attaches token and supports refresh via axios interceptor
      await userService.deleteProfile();
      showNotification('Аккаунт успешно удален.', 'success');
      handleCloseDeleteDialog();
      logout();
      navigate('/login');
    } catch (err: any) {
      console.error(err);
      showNotification('Не удалось удалить аккаунт', 'error');
    }
  };

  const toggleNotifications = async () => {
    const newVal = !notificationsEnabled;
    // optimistically update UI
    setNotificationsEnabled(newVal);
    setUnreadNotificationsCount(0);
    showNotification(`Уведомления ${newVal ? 'включены' : 'отключены'}`, 'info');
    try {
      await userService.updateProfile({ notifications: newVal });
    } catch (err) {
      console.error('Failed to update notification preference', err);
      showNotification('Не удалось сохранить настройку уведомлений', 'error');
      // revert optimistic update
      setNotificationsEnabled(!newVal);
    }
  };

  // load avatars manifest when dialog is opened
  React.useEffect(() => {
    if (!avatarDialogOpen) return;
    let mounted = true;
    (async () => {
      try {
        const res = await fetch('/avatars/manifest.json');
        if (!res.ok) return;
        const data = await res.json();
        if (mounted && Array.isArray(data)) setAvailableAvatars(data as string[]);
      } catch (e) {
        // ignore
      }
    })();
    return () => {
      mounted = false;
    };
  }, [avatarDialogOpen]);

    if (!user) {
    return (
      <Container maxWidth="md" sx={{ mt: 4 }}>
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h5">Пожалуйста, войдите, чтобы просмотреть профиль.</Typography>
          <Button variant="contained" color="primary" sx={{ mt: 2 }} onClick={() => navigate('/login')}>
            Войти
          </Button>
        </Paper>
      </Container>
    );
  }

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 6, display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 2 }}>
        <Box
          sx={{
            width: 140,
            height: 140,
            borderRadius: '999px',
            bgcolor: (theme: any) => theme.palette.background.paper,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            boxShadow: 2,
            overflow: 'hidden',
          }}
        >
          <img
            role="button"
            onClick={() => setAvatarDialogOpen(true)}
            src={`/avatars/${user.avatar || 'avatar1.png'}`}
            alt="User Avatar"
            style={{ width: '100%', height: '100%', objectFit: 'cover', cursor: 'pointer' }}
          />
        </Box>

        <Typography variant="h4" component="h1" sx={{ fontWeight: 600 }}>
          {user.name}
        </Typography>

        <Box sx={{ width: '100%', px: 1 }}>
          <Paper
            role="button"
            onClick={() => { setNameInput(user.name); setNameDialogOpen(true); }}
            sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2, p: 1.25, mb: 1.25, borderRadius: 2, cursor: 'pointer' }}
            elevation={0}
          >
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ width: 40, height: 40, borderRadius: 1.5, display: 'flex', alignItems: 'center', justifyContent: 'center', bgcolor: (theme: any) => theme.palette.action.hover }}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4z" fill="currentColor" />
                  <path d="M6 20a6 6 0 0112 0" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              </Box>
              <Typography variant="subtitle1">Изменить имя</Typography>
            </Box>
            <Box>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 6l6 6-6 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </Box>
          </Paper>
        </Box>

        <Box sx={{ width: '100%', px: 1 }}>
          

          {!HIDE_PROFILE_NOTIFICATIONS && (
            <Paper role="button" onClick={() => toggleNotifications()} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2, p: 1.25, mb: 1.25, borderRadius: 2, cursor: 'pointer' }} elevation={0}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ width: 40, height: 40, borderRadius: 1.5, display: 'flex', alignItems: 'center', justifyContent: 'center', bgcolor: (theme: any) => theme.palette.action.hover }}>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M12 22a2 2 0 002-2H10a2 2 0 002 2zm6-6V11c0-3.07-1.63-5.64-4.5-6.32V4a1.5 1.5 0 10-3 0v.68C7.63 5.36 6 7.92 6 11v5l-1.7 1.7A1 1 0 005 20h14a1 1 0 00.7-1.7L18 16z" fill="currentColor" />
                </svg>
              </Box>
              <Typography variant="subtitle1">Уведомления</Typography>
            </Box>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              {unreadNotificationsCount > 0 && (
                <Box sx={{ bgcolor: 'error.main', color: 'common.white', px: 1, py: 0.5, borderRadius: 2, fontSize: 12 }}>{unreadNotificationsCount}</Box>
              )}
              <Box>
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M9 6l6 6-6 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              </Box>
            </Box>
            </Paper>
          )}

          {/* Logout (penultimate) */}
          <Paper
            role="button"
            onClick={() => {
              // show a brief notification and logout
              showNotification('Выход из аккаунта...', 'info');
              setTimeout(() => {
                logout();
                navigate('/login');
              }, 300);
            }}
            sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2, p: 1.25, mb: 1.25, borderRadius: 2, cursor: 'pointer' }}
            elevation={0}
          >
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ width: 40, height: 40, borderRadius: 1.5, display: 'flex', alignItems: 'center', justifyContent: 'center', bgcolor: (theme: any) => theme.palette.action.hover }}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M16 17l5-5-5-5" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  <path d="M21 12H9" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  <path d="M13 5H6a2 2 0 00-2 2v10a2 2 0 002 2h7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              </Box>
              <Typography variant="subtitle1">Выйти</Typography>
            </Box>
            <Box>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 6l6 6-6 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </Box>
          </Paper>

          <Paper role="button" onClick={handleDeleteAccount} sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 2, p: 1.25, mb: 1.25, borderRadius: 2, cursor: 'pointer' }} elevation={0}>
            <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
              <Box sx={{ width: 40, height: 40, borderRadius: 1.5, display: 'flex', alignItems: 'center', justifyContent: 'center', bgcolor: (theme: any) => theme.palette.action.hover }}>
                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                  <path d="M3 6h18" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  <path d="M19 6l-1 14a2 2 0 01-2 2H8a2 2 0 01-2-2L5 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  <path d="M10 11v6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  <path d="M14 11v6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                  <path d="M9 6V4a1 1 0 011-1h4a1 1 0 011 1v2" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
                </svg>
              </Box>
              <Typography variant="subtitle1">Удалить профиль</Typography>
            </Box>
            <Box>
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                <path d="M9 6l6 6-6 6" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </Box>
          </Paper>
        </Box>
      </Box>

      <Dialog open={openDeleteDialog} onClose={handleCloseDeleteDialog}>
        <DialogTitle>Удалить профиль?</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Вы уверены, что хотите удалить свой профиль? Это действие необратимо.
            Пожалуйста, укажите причину удаления:
          </DialogContentText>
          <FormControl component="fieldset" sx={{ mt: 2, width: '100%' }}>
            <FormLabel component="legend">Причина удаления</FormLabel>
            <RadioGroup aria-label="delete-reason" name="delete-reason-group" value={selectedReason} onChange={(event) => setSelectedReason(event.target.value)}>
              <FormControlLabel value="not_using" control={<Radio />} label="Больше не пользуюсь сервисом" />
              <FormControlLabel value="privacy_concerns" control={<Radio />} label="Проблемы с конфиденциальностью" />
              <FormControlLabel value="bad_experience" control={<Radio />} label="Неудовлетворительный опыт использования" />
              <FormControlLabel value="other" control={<Radio />} label="Другое" />
            </RadioGroup>
          </FormControl>
          {selectedReason === 'other' && (
            <TextField autoFocus margin="dense" id="reason" label="Укажите причину" type="text" fullWidth variant="standard" value={deleteReason} onChange={(e) => setDeleteReason(e.target.value)} sx={{ mt: 2 }} />
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog}>Отмена</Button>
          <Button onClick={handleConfirmDelete} color="error" disabled={!selectedReason}>
            Удалить
          </Button>
        </DialogActions>
      </Dialog>
      {/* Avatar selection dialog */}
      <Dialog open={avatarDialogOpen} onClose={() => setAvatarDialogOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Выберите аватар</DialogTitle>
        <DialogContent>
            <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mt: 1 }}>
            {availableAvatars.length === 0 && (
              <Typography variant="body2">Нет доступных аватаров. Поместите изображения в /avatars и обновите manifest.json.</Typography>
            )}
            {availableAvatars.map((f) => (
              <Box
                key={f}
                onClick={() => setSelectedAvatarFile(f)}
                sx={{
                  width: 72,
                  height: 72,
                  borderRadius: '50%',
                  overflow: 'hidden',
                  cursor: 'pointer',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  border: selectedAvatarFile === f ? '3px solid' : '1px solid',
                  borderColor: selectedAvatarFile === f ? 'primary.main' : 'divider',
                }}
              >
                <img src={`/avatars/${f}`} alt={f} style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
              </Box>
            ))}
          </Box>
        </DialogContent>
        <DialogActions>
          <Button
            variant="outlined"
            sx={{ borderColor: 'grey.300', color: 'text.secondary' }}
            onClick={() => setAvatarDialogOpen(false)}
          >
            Отмена
          </Button>
          <Button
            variant="contained"
            color="primary"
            disabled={!selectedAvatarFile}
            onClick={async () => {
              if (!selectedAvatarFile) return;
              showNotification('Сохраняем аватар...', 'info');
              try {
                // use userService so token attach + refresh happens automatically
                const updatedUser = await userService.updateProfile({ avatar: selectedAvatarFile });
                // update auth context and localStorage
                try {
                  const access = localStorage.getItem('accessToken') || localStorage.getItem('token') || '';
                  const refresh = localStorage.getItem('refreshToken') || '';
                  saveTokens(access, refresh, updatedUser as any);
                } catch (e) {
                  localStorage.setItem('user', JSON.stringify(updatedUser));
                }
                showNotification('Аватар обновлён', 'success');
                setAvatarDialogOpen(false);
                setSelectedAvatarFile(null);
              } catch (err: any) {
                console.error(err);
                showNotification('Не удалось сохранить аватар', 'error');
              }
            }}
          >
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
      {/* Name edit dialog */}
      <Dialog open={nameDialogOpen} onClose={() => setNameDialogOpen(false)} maxWidth="xs" fullWidth>
        <DialogTitle>Изменить имя</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            id="name"
            label="Имя"
            type="text"
            fullWidth
            variant="standard"
            value={nameInput}
            onChange={(e) => setNameInput(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button variant="outlined" sx={{ borderColor: 'grey.300', color: 'text.secondary' }} onClick={() => setNameDialogOpen(false)}>
            Отмена
          </Button>
          <Button
            variant="contained"
            color="primary"
            disabled={!nameInput || nameInput.trim() === '' || nameInput === user.name}
            onClick={async () => {
              showNotification('Сохраняем имя...', 'info');
              try {
                const updatedUser = await userService.updateProfile({ name: nameInput });
                try {
                  const access = localStorage.getItem('accessToken') || localStorage.getItem('token') || '';
                  const refresh = localStorage.getItem('refreshToken') || '';
                  saveTokens(access, refresh, updatedUser as any);
                } catch (e) {
                  localStorage.setItem('user', JSON.stringify(updatedUser));
                }
                showNotification('Имя обновлено', 'success');
                setNameDialogOpen(false);
              } catch (err: any) {
                console.error(err);
                showNotification('Не удалось сохранить имя', 'error');
              }
            }}
          >
            Сохранить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default ProfilePage;