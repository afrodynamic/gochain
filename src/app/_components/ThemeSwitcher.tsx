import {
  DarkMode as DarkModeIcon,
  LightMode as LightModeIcon,
  BrightnessMedium as SystemModeIcon,
} from '@mui/icons-material';
import { Fade, IconButton, Skeleton } from '@mui/material';
import { useColorScheme } from '@mui/material/styles';
import { FC } from 'react';

export const ThemeSwitcher: FC = () => {
  const { mode, setMode } = useColorScheme();

  const toggleMode = () => {
    if (mode === 'light') {
      setMode('dark');
    } else if (mode === 'dark') {
      setMode('system');
    } else {
      setMode('light');
    }
  };

  return (
    <IconButton onClick={toggleMode} color="primary">
      <Fade in={true} timeout={500}>
        {mode === 'dark' ? (
          <DarkModeIcon />
        ) : mode === 'light' ? (
          <LightModeIcon />
        ) : mode === 'system' ? (
          <SystemModeIcon />
        ) : (
          <Skeleton variant="circular" width={24} height={24} />
        )}
      </Fade>
    </IconButton>
  );
};
