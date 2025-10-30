import React, { useState } from 'react';
import { Box, Typography, TextField, Button, Container, Paper } from '@mui/material';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useForm, SubmitHandler } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useAuth } from './AuthContext';
import { useNotification } from './NotificationContext';

interface RegisterFormInputs {
  name: string;
  email: string;
  password: string;
  confirmPassword: string;
}

const schema = yup.object().shape({
  name: yup.string().required('Имя обязательно'),
  email: yup.string().email('Введите корректный email').required('Email обязателен'),
  password: yup.string().min(6, 'Пароль должен быть не менее 6 символов').required('Пароль обязателен'),
  confirmPassword: yup
    .string()
    .oneOf([yup.ref('password')], 'Пароли должны совпадать')
    .required('Подтверждение пароля обязательно'),
});

const RegisterPage: React.FC = () => {
  const { register: authRegister } = useAuth();
  const { showNotification } = useNotification();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<RegisterFormInputs>({
    resolver: yupResolver(schema),
  });

  const onSubmit: SubmitHandler<RegisterFormInputs> = async (data) => {
    setLoading(true);
    setError(null);
    try {
      const result = await authRegister(data.name, data.email, data.password);
      if (result && result.success) {
        showNotification('Регистрация прошла успешно! Вы вошли в систему.', 'success');
        navigate('/'); // Перенаправление на главную или дашборд после успешной регистрации
      } else {
        const msg = (result && result.error) || 'Ошибка регистрации. Возможно, пользователь с таким email уже существует.';
        showNotification(msg, 'error');
      }
    } catch (err) {
      showNotification('Произошла ошибка при регистрации. Попробуйте еще раз.', 'error');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="xs">
      <Paper sx={{ p: 4, mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography variant="h5" component="h1" gutterBottom>
          Регистрация
        </Typography>
        <Box component="form" onSubmit={handleSubmit(onSubmit)} noValidate sx={{ mt: 1, width: '100%' }}>
          <TextField
            margin="normal"
            required
            fullWidth
            id="name"
            label="Имя"
            autoComplete="name"
            autoFocus
            {...register('name')}
            error={!!errors.name}
            helperText={String(errors.name?.message || '')}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            id="email"
            label="Email адрес"
            autoComplete="email"
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
            autoComplete="new-password"
            {...register('password')}
            error={!!errors.password}
            helperText={String(errors.password?.message || '')}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            id="confirmPassword"
            label="Подтвердите пароль"
            type="password"
            autoComplete="new-password"
            {...register('confirmPassword')}
            error={!!errors.confirmPassword}
            helperText={String(errors.confirmPassword?.message || '')}
          />
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            sx={{ mt: 3, mb: 2 }}
            disabled={loading}
          >
            {loading ? 'Регистрация...' : 'Зарегистрироваться'}
          </Button>
          <Box sx={{ display: 'flex', justifyContent: 'center', width: '100%' }}>
            <RouterLink to="/login" style={{ textDecoration: 'none' }}>
              <Button color="primary" size="small">Уже есть аккаунт? Войти</Button>
            </RouterLink>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default RegisterPage;