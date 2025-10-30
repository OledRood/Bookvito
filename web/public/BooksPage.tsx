import React from 'react';
import { Box, Typography, Paper, Container, List, ListItem, ListItemText, ListItemButton, Divider, Button } from '@mui/material';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import { Link as RouterLink, Outlet, useLocation, Navigate } from 'react-router-dom';

const bookRoutes = [
  { label: 'Мои книги', path: '/books/my' },
  { label: 'Забронированные', path: '/books/reserved' },
  { label: 'Моя полка', path: '/books/shelf' },
  { label: 'Прочитанное', path: '/books/read' },
];

const BooksPage: React.FC = () => {
  const location = useLocation();

  // Находим текущую вкладку по пути, чтобы подсветить ее
  const currentTab = bookRoutes.find(route => location.pathname.startsWith(route.path))?.path || false;

  // Если мы на /books, показываем мобильный список секций (tiles) вместо редиректа
  if (location.pathname === '/books') {
    return (
      <Container maxWidth="sm" sx={{ pt: 3 }}>
        <Typography variant="h4" component="h1" sx={{ my: 1 }}>
          Книги
        </Typography>

        <Paper sx={{ mt: 2, borderRadius: 3, overflow: 'hidden', boxShadow: 3 }}>
          <List disablePadding>
            {bookRoutes.map((route, idx) => (
              <Box key={route.path}>
                <ListItem disablePadding>
                      <ListItemButton
                        component={RouterLink}
                        to={route.path}
                        sx={{
                          px: 3,
                          ...(idx === 0 ? { pt: 2 } : {}),
                          ...(idx === bookRoutes.length - 1 ? { pb: 2 } : {}),
                          display: 'flex',
                          justifyContent: 'space-between',
                          alignItems: 'center',
                        }}
                      >
                    <ListItemText
                      primary={route.label}
                      primaryTypographyProps={{ color: 'text.primary', fontWeight: 500 }}
                    />
                    <ChevronRightIcon sx={{ color: 'text.secondary' }} />
                  </ListItemButton>
                </ListItem>
                {idx !== bookRoutes.length - 1 && <Divider component="li" sx={{ borderColor: 'divider' }} />}
              </Box>
            ))}
          </List>
        </Paper>
      </Container>
    );
  }

  // Child routes: show back button and the child content only (no global title or outer card)
  return (
    <Container maxWidth="lg">
      {location.pathname === '/books' ? (
        // landing view keeps the big title
        <>
          <Typography variant="h4" component="h1" sx={{ my: 3 }}>
            Книги
          </Typography>
          <Box sx={{ mb: 2 }} />
          <Paper sx={{ p: 3 }}>
            <Outlet />
          </Paper>
        </>
      ) : (
        // child view: only back button + child content (no outer Paper/title)
        <>
          <Box sx={{ mb: 2 }}>
            <Button component={RouterLink} to="/books" startIcon={<ArrowBackIcon />} variant="text">
              Назад
            </Button>
          </Box>
          <Outlet />
        </>
      )}
    </Container>
  );
};

export default BooksPage;