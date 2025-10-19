'use client';
import { InitColorSchemeScript } from '@mui/material';
import { AppRouterCacheProvider } from '@mui/material-nextjs/v13-appRouter';
import CssBaseline from '@mui/material/CssBaseline';
import { ThemeProvider } from '@mui/material/styles';
import type {} from '@mui/material/themeCssVarsAugmentation';
import { FC, ReactNode, useMemo } from 'react';

import { getTheme } from './theme';

export const ThemeRegistry: FC<{ children: ReactNode }> = ({ children }) => {
  const theme = useMemo(() => getTheme(), []);

  return (
    <>
      <InitColorSchemeScript attribute="data" />

      <AppRouterCacheProvider options={{ enableCssLayer: true }}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          {children}
        </ThemeProvider>
      </AppRouterCacheProvider>
    </>
  );
};
