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
const navItems = ['Home'];

const hrefFor = (item: string) =>
  item === 'Home' ? '/' : `/#${item.toLowerCase()}`;

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
                key={item}
                component={NextLink}
                href={hrefFor(item)}
                color="primary"
              >
                {item}
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
                <ListItem key={item} disablePadding>
                  <ListItemButton
                    component={NextLink}
                    href={hrefFor(item)}
                    sx={{ textAlign: 'center' }}
                  >
                    <ListItemText primary={item} />
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
