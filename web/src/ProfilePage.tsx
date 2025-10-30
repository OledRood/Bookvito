import React, { useState } from 'react';
import {
  Box,
  Typography,
  Container,
  Paper,
  Button,
  Switch,
  FormControlLabel,
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
  Badge,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import { useNavigate } from 'react-router-dom';
import { useNotification } from '../contexts/NotificationContext';

const ProfilePage: React.FC = () => {
  const { showNotification } = useNotification();
  const { user, logout } = useAuth();
 
  const navigate = useNavigate();

  const [notificationsEnabled, setNotificationsEnabled] = useState(true);
  const [openDeleteDialog, setOpenDeleteDialog] = useState(false);
  const [deleteReason, setDeleteReason] = useState('');
  const [selectedReason, setSelectedReason] = useState('');
  const [unreadNotificationsCount, setUnreadNotificationsCount] = useState(3); // Индикатор непрочитанных

  const handleNotificationsChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setNotificationsEnabled(event.target.checked);
    showNotification(`Уведомления ${event.target.checked ? 'включены' : 'отключены'}`, 'info');
    setUnreadNotificationsCount(0); // Сбрасываем счетчик при взаимодействии
  };

  const handleDeleteAccount = () => {
    setOpenDeleteDialog(true);
  };

  const handleCloseDeleteDialog = () => {
    setOpenDeleteDialog(false);
    setDeleteReason('');
    setSelectedReason('');
  };

  const handleConfirmDelete = () => {
    console.log('Deleting account for user:', user?.email);
    console.log('Reason:', selectedReason === 'other' ? deleteReason : selectedReason);
    // Симуляция API-запроса на удаление аккаунта (с задержкой)
    setTimeout(async () => {
      try {
        showNotification('Аккаунт успешно удален.', 'success');
        logout(); // Выход пользователя после удаления
        navigate('/login'); // Перенаправление на страницу входа
      } catch (error) {
        showNotification('Ошибка при удалении аккаунта.', 'error');
      }
    }, 1000);
    handleCloseDeleteDialog();
  };

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

  // Стилизация шестиугольного аватара (базовый пример с использованием clip-path)
  const hexagonalClipPath = {
    clipPath: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)',
    shapeOutside: 'polygon(50% 0%, 100% 25%, 100% 75%, 50% 100%, 0% 75%, 0% 25%)',
    width: 120,
    height: 120,
    backgroundColor: 'grey.300', // Фоновый цвет-заглушка
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: 'hidden',
  };

  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4, mt: 4 }}>
        <Box sx={{ display: 'flex', flexDirection: 'column', alignItems: 'center', mb: 4 }}>
          <Box sx={hexagonalClipPath}>
            <img
              src={user.avatar || 'https://via.placeholder.com/150/cccccc/ffffff?text=User'}
              alt="User Avatar"
              style={{ width: '100%', height: '100%', objectFit: 'cover' }}
            />
          </Box>
          <Typography variant="h5" component="h1" sx={{ mt: 2 }}>
            {user.name}
          </Typography>
          <Typography variant="body1" color="text.secondary">
            {user.email}
          </Typography>
        </Box>

        <Box sx={{ mb: 3 }}>
          <Typography variant="h6" gutterBottom>Настройки</Typography>
          {/* Переключатель темы убран со страницы профиля (дублирует глобальный переключатель в хедере) */}
          <Badge badgeContent={unreadNotificationsCount} color="error" invisible={unreadNotificationsCount === 0}>
            <FormControlLabel
              control={<Switch checked={notificationsEnabled} onChange={handleNotificationsChange} />}
              label="Уведомления"
              sx={{ mr: 0 }} // Убираем лишний отступ, чтобы Badge был ближе
            />
          </Badge>

        </Box>

        <Box sx={{ mt: 4 }}>
          <Button
            variant="outlined"
            color="error"
            onClick={handleDeleteAccount}
            fullWidth
          >
            Удалить профиль
          </Button>
          <Button
            variant="contained"
            color="primary"
            onClick={logout}
            fullWidth
            sx={{ mt: 2 }}
          >
            Выйти
          </Button>
        </Box>
      </Paper>

      <Dialog open={openDeleteDialog} onClose={handleCloseDeleteDialog}>
        <DialogTitle>Удалить профиль?</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Вы уверены, что хотите удалить свой профиль? Это действие необратимо.
            Пожалуйста, укажите причину удаления:
          </DialogContentText>
          <FormControl component="fieldset" sx={{ mt: 2, width: '100%' }}>
            <FormLabel component="legend">Причина удаления</FormLabel>
            <RadioGroup
              aria-label="delete-reason"
              name="delete-reason-group"
              value={selectedReason}
              onChange={(event) => setSelectedReason(event.target.value)}
            >
              <FormControlLabel value="not_using" control={<Radio />} label="Больше не пользуюсь сервисом" />
              <FormControlLabel value="privacy_concerns" control={<Radio />} label="Проблемы с конфиденциальностью" />
              <FormControlLabel value="bad_experience" control={<Radio />} label="Неудовлетворительный опыт использования" />
              <FormControlLabel value="other" control={<Radio />} label="Другое" />
            </RadioGroup>
          </FormControl>
          {selectedReason === 'other' && (
            <TextField
              autoFocus
              margin="dense"
              id="reason"
              label="Укажите причину"
              type="text"
              fullWidth
              variant="standard"
              value={deleteReason}
              onChange={(e) => setDeleteReason(e.target.value)}
              sx={{ mt: 2 }}
            />
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDeleteDialog}>Отмена</Button>
          <Button onClick={handleConfirmDelete} color="error" disabled={!selectedReason}>
            Удалить
          </Button>
        </DialogActions>
      </Dialog>
    </Container>
  );
};

export default ProfilePage;