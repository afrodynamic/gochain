'use client';
import { useHealthStatusQuery } from '@/app/(features)/(dashboard)/useHealthQuery';
import type { HealthStatusResponse } from '@/app/_utils/types';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import ErrorIcon from '@mui/icons-material/Error';
import HealthAndSafetyIcon from '@mui/icons-material/HealthAndSafety';
import HelpOutlineIcon from '@mui/icons-material/HelpOutline';
import {
  Box,
  Card,
  CardContent,
  CardHeader,
  Chip,
  LinearProgress,
  Stack,
  Tooltip,
  Typography,
} from '@mui/material';
import { useMemo } from 'react';

export function Health() {
  const { data, status } = useHealthStatusQuery();

  const { label, color, icon } = useMemo(() => {
    const s =
      (data as HealthStatusResponse | undefined)?.status ?? 'unavailable';
    if (s === 'ok') {
      return {
        label: 'OK',
        color: 'success' as const,
        icon: <CheckCircleIcon />,
      };
    } else if (s === 'down') {
      return { label: 'ERROR', color: 'error' as const, icon: <ErrorIcon /> };
    } else {
      return {
        label: 'UNAVAILABLE',
        color: 'default' as const,
        icon: <HelpOutlineIcon />,
      };
    }
  }, [data]);

  const timestamp = (data as HealthStatusResponse | undefined)?.timestamp;
  const exact = timestamp ? new Date(timestamp) : null;
  const lastChecked = exact ? exact.toLocaleString() : 'Unavailable';

  return (
    <Card
      elevation={6}
      className="max-w-lg w-full border"
      sx={{ borderColor: 'primary.main' }}
    >
      {status === 'pending' && <LinearProgress />}
      <CardHeader
        avatar={<HealthAndSafetyIcon fontSize="large" />}
        title="Blockchain Health"
        subheader={
          status === 'error' ? 'Failed to fetch status' : 'Live status overview'
        }
        sx={{ '& .MuiCardHeader-title': { fontWeight: 700, fontSize: 22 } }}
      />
      <CardContent>
        <Stack spacing={2}>
          <Stack
            direction="row"
            alignItems="center"
            justifyContent="space-between"
          >
            <Typography variant="body1" sx={{ fontWeight: 600 }}>
              Status
            </Typography>

            <Chip
              icon={icon}
              label={label}
              color={color}
              variant={color === 'default' ? 'outlined' : 'filled'}
              className="font-bold"
            />
          </Stack>

          <Stack
            direction="row"
            alignItems="center"
            justifyContent="space-between"
          >
            <Typography variant="body1" sx={{ fontWeight: 600 }}>
              Last Checked
            </Typography>

            <Tooltip title={lastChecked}>
              <Box>
                <Typography variant="body2" color="text.secondary">
                  {lastChecked}
                </Typography>
              </Box>
            </Tooltip>
          </Stack>
        </Stack>
      </CardContent>
    </Card>
  );
}
