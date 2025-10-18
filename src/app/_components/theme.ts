import { createTheme, responsiveFontSizes } from '@mui/material/styles';
import { Fira_Code, Inter } from 'next/font/google';

export const inter = Inter({
  variable: '--font-inter',
  subsets: ['latin'],
  display: 'swap',
});

export const firaCode = Fira_Code({
  variable: '--font-fira',
  subsets: ['latin'],
  display: 'swap',
});

export const getTheme = () => {
  return responsiveFontSizes(
    createTheme({
      typography: {
        fontFamily: `${inter.style.fontFamily}, ${firaCode.style.fontFamily}, system-ui, sans-serif`,
        h1: {
          fontFamily: `${firaCode.style.fontFamily}, monospace`,
          fontWeight: 700,
        },
        h2: {
          fontFamily: `${firaCode.style.fontFamily}, monospace`,
          fontWeight: 600,
        },
        h3: {
          fontFamily: `${firaCode.style.fontFamily}, monospace`,
          fontWeight: 500,
        },
        h4: {
          fontFamily: `${firaCode.style.fontFamily}, monospace`,
          fontWeight: 500,
        },
        h5: {
          fontFamily: `${firaCode.style.fontFamily}, monospace`,
          fontWeight: 400,
        },
        h6: {
          fontFamily: `${firaCode.style.fontFamily}, monospace`,
          fontWeight: 300,
        },
      },
      colorSchemes: {
        dark: true,
        light: true,
      },
      cssVariables: {
        colorSchemeSelector: 'data',
      },
    })
  );
};

export const muiColors: Array<
  'default' | 'primary' | 'secondary' | 'success' | 'info' | 'error' | 'warning'
> = ['default', 'primary', 'secondary', 'success', 'info', 'error', 'warning'];
