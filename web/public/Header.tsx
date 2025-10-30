import React from 'react';
import { AppBar, Toolbar, Typography, Button, Box } from '@mui/material';
import { Link } from 'react-router-dom';
import ThemeSwitcher from './ThemeSwitcher';

interface HeaderProps {
  isDarkMode: boolean;
  toggleTheme: () => void;
}

const Header: React.FC<HeaderProps> = ({ isDarkMode, toggleTheme }) => {
  return (
    <AppBar position="static" color="primary">
      <Toolbar>
        <Typography
          variant="h6"
          component={Link}
          to="/"
          sx={{ flexGrow: 1, textDecoration: 'none', color: 'inherit' }}
        >
          Material 3 App
        </Typography>
        <Box sx={{ display: { xs: 'none', sm: 'block' } }}>
          <Button color="inherit" component={Link} to="/">Главная</Button>
          <Button color="inherit" component={Link} to="/about">О нас</Button>
          <Button color="inherit" component={Link} to="/contact">Контакты</Button>
        </Box>
        <ThemeSwitcher isDarkMode={isDarkMode} toggleTheme={toggleTheme} />
      </Toolbar>
    </AppBar>
  );
};

export default Header;