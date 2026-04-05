import React from 'react';
import {
  Typography,
  Container,
  Box,
  IconButton,
  Menu,
  MenuItem,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  BottomNavigation,
  BottomNavigationAction,
  Paper,
  useTheme,
  useMediaQuery,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
  Button,
} from '@mui/material';
import { Outlet, Link as RouterLink, useLocation, useNavigate } from 'react-router-dom';
import { BOOKS_NEW_PATH, MODERATION_PATH } from '../src/routing/paths';
import ThemeSwitcher from './ThemeSwitcher';
import { useAuth } from './AuthContext';
import SeoManager from './SeoManager';
import AccountCircle from '@mui/icons-material/AccountCircle';
import HomeIcon from '@mui/icons-material/Home';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import BookIcon from '@mui/icons-material/MenuBook';
import PersonIcon from '@mui/icons-material/Person';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';
import GavelIcon from '@mui/icons-material/Gavel';
// TO_BE_REMOVED_BEFORE_RELEASE
import PaletteIcon from '@mui/icons-material/Palette';

const drawerWidth = 88;

const Layout: React.FC = () => {
  const { isAuthenticated, user, logout, isAdmin, isModer } = useAuth();

  // Nav items that depend on role — evaluated at render time
  const navItems = [
    { text: 'Главная', path: '/', icon: <HomeIcon /> },
    { text: 'Создать', path: BOOKS_NEW_PATH, icon: <AddCircleOutlineIcon />, requiresAuth: true },
    { text: 'Книги', path: '/books', icon: <BookIcon />, requiresAuth: true },
    { text: 'Профиль', path: '/profile', icon: <PersonIcon />, auth: true },
    ...(isModer ? [{ text: 'Модерация', path: MODERATION_PATH, icon: <GavelIcon />, requiresAuth: true }] : []),
    ...(isAdmin ? [{ text: 'Админ', path: '/admin', icon: <AdminPanelSettingsIcon />, requiresAuth: true }] : []),
    // TO_BE_REMOVED_BEFORE_RELEASE
    { text: 'Палетка', path: '/color-palette', icon: <PaletteIcon />, hidden: true },
  ];

  const location = useLocation();
  const hideAppBar = location.pathname.startsWith('/book/');
  const theme = useTheme();
  const isDesktop = useMediaQuery(theme.breakpoints.up('md'));
  const navigate = useNavigate();

  // dialog state for blocked actions
  const [authDialogOpen, setAuthDialogOpen] = React.useState(false);
  const [pendingAction, setPendingAction] = React.useState<string | null>(null);

  const openAuthDialog = (actionPath?: string) => {
    setPendingAction(actionPath || null);
    setAuthDialogOpen(true);
  };

  const closeAuthDialog = () => {
    setAuthDialogOpen(false);
    setPendingAction(null);
  };

  const handleAuthProceed = () => {
    closeAuthDialog();
    console.debug('[Layout] handleAuthProceed -> navigate to /login, pendingAction=', pendingAction);
    navigate('/login');
  };

  const authDialog = (
    <Dialog open={authDialogOpen} onClose={closeAuthDialog}>
      <DialogTitle>Требуется вход</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Для доступа к этой функции необходимо войти в аккаунт. Перейти на страницу входа?
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={closeAuthDialog}>Отмена</Button>
        <Button component={RouterLink} to="/login" onClick={handleAuthProceed} variant="contained">Войти</Button>
      </DialogActions>
    </Dialog>
  );

  // keep items which are not explicitly hidden; also respect auth flags
  const filteredNavItems = navItems.filter(item => item.hidden !== true && (!item.auth || isAuthenticated));

  if (isDesktop) {
    return (
      <>
        <SeoManager />
        <Box sx={{ display: 'flex' }}>
          <Drawer
            variant="permanent"
            sx={{
              width: drawerWidth,
              flexShrink: 0,
              [`& .MuiDrawer-paper`]: {
                width: drawerWidth,
                boxSizing: 'border-box',
                display: 'flex',
                flexDirection: 'column',
                height: '100%'
              },
            }}
          >
            <Box sx={{ flex: 1, overflow: 'auto' }}>
              <List>
                {filteredNavItems.map((item) => {
                  const disabled = Boolean(item.requiresAuth && !isAuthenticated);
                  return (
                    // add vertical spacing between buttons in desktop drawer
                    <ListItem key={item.text} disablePadding sx={{ mb: 2 }}>
                      <ListItemButton
                        component={disabled ? 'div' : RouterLink}
                        to={disabled ? undefined : item.path}
                        onClick={(e: React.MouseEvent) => {
                          if (disabled) {
                            e.preventDefault();
                            openAuthDialog(item.path);
                          }
                        }}
                        disableRipple
                        disableTouchRipple
                        selected={location.pathname === item.path}
                        sx={{
                          display: 'flex',
                          flexDirection: 'column',
                          justifyContent: 'center',
                          alignItems: 'center',
                          width: '96px',
                          height: '64px',
                          gap: '0px',
                          padding: '6px 0px',
                          color: disabled ? 'var(--md-sys-color-on-surface-variant)' : 'var(--md-sys-color-on-surface)',
                          cursor: disabled ? 'not-allowed' : 'pointer',
                          pointerEvents: 'auto',
                          transition: 'background-color .18s ease, color .18s ease, transform .12s ease',
                          '&:hover': {
                            backgroundColor: 'transparent',
                          },
                          '&.Mui-selected': {
                            color: 'var(--md-sys-color-primary)',
                            backgroundColor: 'transparent',
                          },
                          '&.Mui-selected:hover': {
                            backgroundColor: 'transparent',
                          },
                          '&:hover .nav-highlight:not([data-selected])': {
                            backgroundColor: 'var(--md-sys-color-surface-dim)',
                          },
                        }}
                      >
                        <ListItemIcon
                          sx={{
                            minWidth: 0,
                            display: 'flex',
                            justifyContent: 'center',
                            alignItems: 'center',
                            position: 'relative',
                            height: 40,
                            width: 56,
                            overflow: 'visible',
                            p: 0,
                          }}
                        >
                          <Box
                            className="nav-highlight"
                            data-selected={location.pathname === item.path ? '1' : undefined}
                            sx={{
                              width: 56,
                              height: 32,
                              borderRadius: 18,
                              backgroundColor: location.pathname === item.path ? 'var(--md-sys-color-primary-container)' : 'transparent',
                              transition: 'background-color 0.12s ease',
                              display: 'flex',
                              justifyContent: 'center',
                              alignItems: 'center',
                              position: 'relative',
                              zIndex: 0,
                            }}
                          >
                            {item.icon}
                          </Box>
                        </ListItemIcon>

                        <ListItemText
                          primary={item.text}
                          primaryTypographyProps={{
                            align: 'center',
                            sx: {
                              color: 'var(--M3/sys/light/on-surface-variant, var(--Schemes-OnSurfaceVariant, rgba(73, 69, 79, 1)))',
                              fontFamily: 'var(--Static-LabelMedium-Font, Roboto)',
                              fontSize: '12px',
                              fontWeight: 500,
                              lineHeight: '16px',
                              letterSpacing: '0.5px'
                            }
                          }}
                          sx={{ my: 0 }}
                        />
                      </ListItemButton>
                    </ListItem>
                  );
                })}
              </List>
            </Box>
            <Box sx={{ px: 2, py: 1.5, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
              <ThemeSwitcher />
            </Box>
          </Drawer>
          <Box component="div" sx={{ flexGrow: 1, px: 3, pt: 1.5, pb: 3 }}>
            <Outlet />
          </Box>
          {authDialog}
        </Box>
      </>
    );
  }

  // --- Мобильная версия с нижней навигацией ---
  return (
    <>
      <SeoManager />
      <Box component="div" sx={{ pb: 7 }}>
        <Outlet />
      </Box>
      <Paper
        sx={{ position: 'fixed', bottom: 0, left: 0, right: 0, bgcolor: 'var(--md-sys-color-surface)', borderTop: 1, borderColor: 'var(--md-sys-color-outline)', borderRadius: '16px 16px 0 0', overflow: 'hidden' }}
        elevation={3}
      >
        <BottomNavigation showLabels value={location.pathname}>
          {filteredNavItems.map((item) => {
            const disabled = Boolean(item.requiresAuth && !isAuthenticated);
            return (
              <BottomNavigationAction
                key={item.text}
                label={item.text}
                value={item.path}
                icon={item.icon}
                component={disabled ? 'button' : RouterLink}
                to={disabled ? undefined : item.path}
                onClick={(e: React.MouseEvent<HTMLButtonElement | HTMLAnchorElement>) => {
                  if (disabled) {
                    e.preventDefault();
                    openAuthDialog(item.path);
                  }
                }}
                aria-disabled={disabled}
                tabIndex={disabled ? -1 : 0}
                sx={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  minWidth: 72,
                  px: 1,
                  color: disabled ? 'var(--md-sys-color-on-surface-variant)' : 'inherit',
                  pointerEvents: 'auto',
                  transition: 'color .18s ease, transform .12s ease',
                  '& .MuiBottomNavigationAction-label': {
                    fontSize: '0.75rem',
                    mt: 0.5,
                    lineHeight: 1,
                  },
                  '&.Mui-selected': {
                    color: 'var(--md-sys-color-primary)',
                    transform: 'translateY(-2px)',
                    '& .MuiBottomNavigationAction-label': { fontWeight: 600 },
                  }
                }}
              />
            );
          })}
        </BottomNavigation>
      </Paper>
      {authDialog}
    </>
  );
};

export default Layout;
