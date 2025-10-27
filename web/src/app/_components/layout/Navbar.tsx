'use client';
import { ThemeSwitcher } from '@/app/_components/ThemeSwitcher';
import MenuIcon from '@mui/icons-material/Menu';
import {
  AppBar,
  Box,
  Button,
  CssBaseline,
  Divider,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemText,
  Link as MUILink,
  Toolbar,
  Typography,
} from '@mui/material';
import NextLink from 'next/link';
import { useState } from 'react';

const drawerWidth = 240;
const navItems = [
  { label: 'Home', href: '/' },
  { label: 'Blocks', href: '/blocks' },
  { label: 'Wallet', href: '/wallet' },
  { label: 'Transactions', href: '/transactions' },
];

export function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false);
  const handleDrawerToggle = () => setMobileOpen((v) => !v);

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />

      <AppBar
        component="nav"
        className="backdrop-blur-sm"
        color="transparent"
        enableColorOnDark
      >
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerToggle}
            sx={{ mr: 2, display: { sm: 'none' } }}
          >
            <MenuIcon />
          </IconButton>

          <MUILink
            component={NextLink}
            href="/"
            sx={{ flexGrow: 1, display: { xs: 'none', sm: 'block' } }}
            underline="none"
          >
            <Typography variant="h6" fontWeight="bold">
              Gochain
            </Typography>
          </MUILink>

          <Box
            sx={{ display: { xs: 'none', sm: 'block' } }}
            className="space-x-2"
          >
            {navItems.map((item) => (
              <Button
                key={item.label}
                component={NextLink}
                href={item.href}
                color="primary"
              >
                {item.label}
              </Button>
            ))}
          </Box>

          <div className="flex flex-grow" />

          <ThemeSwitcher />
        </Toolbar>

        <Divider />
      </AppBar>

      <nav>
        <Drawer
          variant="temporary"
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{ keepMounted: true }}
          sx={{
            display: { xs: 'block', sm: 'none' },
            '& .MuiDrawer-paper': {
              boxSizing: 'border-box',
              width: drawerWidth,
            },
          }}
        >
          <Box onClick={handleDrawerToggle} sx={{ textAlign: 'center' }}>
            <Typography variant="h6" sx={{ my: 2 }}>
              Gochain
            </Typography>

            <Divider />

            <List>
              {navItems.map((item) => (
                <ListItem key={item.label} disablePadding>
                  <ListItemButton
                    component={NextLink}
                    href={item.href}
                    sx={{ textAlign: 'center' }}
                  >
                    <ListItemText primary={item.label} />
                  </ListItemButton>
                </ListItem>
              ))}
            </List>
          </Box>
        </Drawer>
      </nav>

      <Toolbar />
    </Box>
  );
}
