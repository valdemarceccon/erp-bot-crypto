// src/layouts/MainLayout.tsx
import {
  AppBar,
  Box,
  Toolbar,
  Typography,
  IconButton,
  MenuItem,
  Menu,
  Drawer,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  MenuItemProps,
  Button,
} from '@mui/material';
import { Menu as MenuIcon, AccountCircle, Key } from '@mui/icons-material';
import Link from 'next/link';
import React, { forwardRef, useState } from 'react';
import PersonIcon from '@mui/icons-material/Person';

type MainLayoutProps = {
  children: React.ReactNode;
  loggedIn: boolean;
};

interface MenuItemsProps {
  onToggleDrawer(open: boolean): void;
}

const MenuItems: React.FC<MenuItemsProps> = ({ onToggleDrawer }) => {
  return (
    <Box
      role="presentation"
      onClick={() => onToggleDrawer(false)}
      onKeyDown={() => onToggleDrawer(false)}
    >
      <List>
        <ListItem key="bot_key" disablePadding>
          <ListItemButton>
            <ListItemIcon>
              <Key />
            </ListItemIcon>
            <ListItemText primary="Chaves" />
          </ListItemButton>
        </ListItem>
        <ListItem key="users" disablePadding>
          <ListItemButton>
            <ListItemIcon>
              <PersonIcon />
            </ListItemIcon>
            <ListItemText primary="Usuários" />
          </ListItemButton>
        </ListItem>
      </List>
    </Box>
  );
}

const MainLayout: React.FC<MainLayoutProps> = ({ children, loggedIn }) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const [state, setState] = React.useState(false);

  const toggleDrawer =
    (open: boolean) =>
      (event: React.KeyboardEvent | React.MouseEvent) => {
        if (
          event.type === 'keydown' &&
          ((event as React.KeyboardEvent).key === 'Tab' ||
            (event as React.KeyboardEvent).key === 'Shift')
        ) {
          return;
        }

        setState(open);
      };

  return (
    <Box sx={{ flexGrow: 1, minHeight: '100vh', display: 'flex', flexDirection: 'column' }}>
      <AppBar position="static">
        <Toolbar>
          <IconButton edge="start" color="inherit" aria-label="menu" sx={{ mr: 2 }} onClick={toggleDrawer(!state)}>
            <MenuIcon />
          </IconButton>
          <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
            Dashboard
          </Typography>
          {loggedIn ? (
            <>
              <IconButton onClick={handleMenu} color="inherit">
                <AccountCircle />
              </IconButton>
              <Menu
                anchorEl={anchorEl}
                open={Boolean(anchorEl)}
                onClose={handleClose}
              >
                <MenuItem onClick={handleClose}>Logout</MenuItem>
              </Menu>
            </>
          ) : (
            <>
              <Button href="/register" >
                Register
              </Button>

              <Button href="/login" >
                Login
              </Button>
            </>
          )}


        </Toolbar>
      </AppBar>
      <Drawer
        anchor="left"
        open={state}
        onClose={toggleDrawer(false)}>
        <MenuItems onToggleDrawer={toggleDrawer} />
      </Drawer>

      <Box sx={{ flexGrow: 1, p: 3 }}>{children}</Box>
    </Box>
  );
};

export default MainLayout;
