import React, { useState } from 'react';
import { Box, Typography, TextField, Button, Container, Paper } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useForm, SubmitHandler } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useAuth } from './AuthContext';
import { useNotification } from './NotificationContext';

interface LoginFormInputs {
  email: string;
  password: string;
}

const schema = yup.object().shape({
  email: yup.string().email('Введите корректный email').required('Email обязателен'),
  password: yup.string().min(6, 'Пароль должен быть не менее 6 символов').required('Пароль обязателен'),
});

const LoginPage: React.FC = () => {
  const { login } = useAuth();
  const { showNotification } = useNotification();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormInputs>({
    resolver: yupResolver(schema),
  });

  const onSubmit: SubmitHandler<LoginFormInputs> = async (data) => {
    setLoading(true);
    setError(null);
    try {
      const success = await login(data.email, data.password);
      if (success) {
        showNotification('Вы успешно вошли!', 'success');
        navigate('/'); // Перенаправление на главную или дашборд после успешного входа
      } else {
        showNotification('Неверный email или пароль.', 'error');
      }
    } catch (err) {
      showNotification('Произошла ошибка при входе. Попробуйте еще раз.', 'error');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="xs">
      <Paper sx={{ p: 4, mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography variant="h5" component="h1" gutterBottom>
          Вход
        </Typography>
        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate sx={{ mt: 1, width: '100%' }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email адрес"
            autoComplete="email"
            autoFocus
            {...register('email')}
            error={!!errors.email}
            helperText={String(errors.email?.message || '')}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            id="password"
            label="Пароль"
            type="password"
            autoComplete="current-password"
            {...register('password')}
            error={!!errors.password}
            helperText={String(errors.password?.message || '')}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            sx={{ mt: 3, mb: 2 }}
            disabled={loading}
          >
            {loading ? 'Вход...' : 'Войти'}
          </Button>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', width: '100%' }}>
            <RouterLink to="/forgot-password" style={{ textDecoration: 'none' }}>
              <Button color="primary" size="small">Забыли пароль?</Button>
            </RouterLink>
            <RouterLink to="/register" style={{ textDecoration: 'none' }}>
              <Button color="primary" size="small">Зарегистрироваться</Button>
            </RouterLink>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default LoginPage;