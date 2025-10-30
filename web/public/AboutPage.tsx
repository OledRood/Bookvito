import React from 'react';
import { Box, Typography, Paper, Container } from '@mui/material';

const AboutPage: React.FC = () => {
  return (
    <Container maxWidth="md">
      <Paper sx={{ p: 4, mt: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          О нас
        </Typography>
        <Typography variant="body1">Это страница с информацией о проекте Bookvito.</Typography>
        <Typography variant="body1" sx={{ mt: 2 }}>Наша миссия - создать лучшую платформу для обмена и продажи книг.</Typography>
      </Paper>
    </Container>
  );
};

export default AboutPage;