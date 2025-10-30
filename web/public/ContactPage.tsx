import React from 'react';
import { Box, Typography } from '@mui/material';

const ContactPage: React.FC = () => {
  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" component="h1" gutterBottom>Контакты</Typography>
      <Typography variant="body1">Свяжитесь с нами по адресу example@example.com</Typography>
    </Box>
  );
};

export default ContactPage;