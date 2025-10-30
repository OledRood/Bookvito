import React, { useState } from 'react';
import { Box, Typography, TextField, Button, Container, Paper } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { useForm, SubmitHandler } from 'react-hook-form';
import { yupResolver } from '@hookform/resolvers/yup';
import * as yup from 'yup';
import { useNotification } from './contexts/NotificationContext';

interface ForgotPasswordFormInputs {
  email: string;
}

const schema = yup.object().shape({
  email: yup.string().email('Введите корректный email').required('Email обязателен'),
});

const ForgotPasswordPage: React.FC = () => {
  const { showNotification } = useNotification();
  const [loading, setLoading] = useState<boolean>(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ForgotPasswordFormInputs>({
    resolver: yupResolver(schema),
  });

  const onSubmit: SubmitHandler<ForgotPasswordFormInputs> = async (data) => {
    setLoading(true);
    try {
      // Симуляция API-запроса на сброс пароля
      await new Promise((resolve) => setTimeout(resolve, 1500)); // Симуляция задержки сети
      console.log('Password reset request for:', data.email);
      showNotification('Если ваш email зарегистрирован, инструкции по сбросу пароля отправлены на вашу почту.', 'success');
    } catch (err) {
      showNotification('Произошла ошибка при отправке запроса. Попробуйте еще раз.', 'error');
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="xs">
      <Paper sx={{ p: 4, mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography variant="h5" component="h1" gutterBottom>
          Восстановление пароля
        </Typography>
        <Typography variant="body2" color="text.secondary" align="center" sx={{ mb: 2 }}>
          Введите ваш email, и мы отправим вам ссылку для сброса пароля.
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
          <Button
            type="submit"
            fullWidth
            variant="contained"
            color="primary"
            sx={{ mt: 3, mb: 2 }}
            disabled={loading}
          >
            {loading ? 'Отправка...' : 'Сбросить пароль'}
          </Button>
          <Box sx={{ display: 'flex', justifyContent: 'center', width: '100%' }}>
            <RouterLink to="/login" style={{ textDecoration: 'none' }}>
              <Button color="primary" size="small">Вернуться ко входу</Button>
            </RouterLink>
          </Box>
        </Box>
      </Paper>
    </Container>
  );
};

export default ForgotPasswordPage;