'use client';
import { useKeysNew } from '@/app/(features)/(dashboard)/useKeysNew';
import { CopyableField } from '@/app/_components/core/CopyableField';
import type { KeyPair } from '@/app/_utils/types';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import HourglassEmptyIcon from '@mui/icons-material/HourglassEmpty';
import VpnKeyIcon from '@mui/icons-material/VpnKey';
import {
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  LinearProgress,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material';
import { useMemo, useRef, useState } from 'react';

export function KeyGenerator() {
  const [seed, setSeed] = useState('');
  const generatedAtRef = useRef<string | null>(null);
  const { mutate, data, status, isPending, reset } = useKeysNew();

  const generated = Boolean(data);
  if (generated && !generatedAtRef.current)
    generatedAtRef.current = new Date().toISOString();

  const { label, color, icon } = useMemo(() => {
    if (generated)
      return {
        label: 'GENERATED',
        color: 'success' as const,
        icon: <CheckCircleIcon />,
      };
    return {
      label: 'READY',
      color: 'default' as const,
      icon: <HourglassEmptyIcon />,
    };
  }, [generated]);

  const lastGenerated = generatedAtRef.current
    ? new Date(generatedAtRef.current).toLocaleString()
    : '—';
  const keys = (data as KeyPair | undefined) ?? null;

  return (
    <Card elevation={6} className="max-w-lg w-full">
      {isPending && <LinearProgress />}

      <CardHeader
        avatar={<VpnKeyIcon fontSize="large" />}
        title="Keypair"
        subheader={
          status === 'error' ? 'Failed to generate' : 'Generate a new keypair'
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

          <TextField
            label="Optional Seed"
            size="small"
            fullWidth
            value={seed}
            onChange={(e) => setSeed(e.target.value)}
            placeholder="leave blank for random"
          />

          <Box display="flex" gap={1}>
            <Button
              variant="contained"
              onClick={() => mutate({ seed: seed || undefined })}
              disabled={isPending}
            >
              Generate
            </Button>

            <Button
              variant="outlined"
              onClick={() => {
                reset();
                setSeed('');
                generatedAtRef.current = null;
              }}
              disabled={isPending}
            >
              Reset
            </Button>
          </Box>

          <Stack
            direction="row"
            alignItems="center"
            justifyContent="space-between"
          >
            <Typography variant="body1" sx={{ fontWeight: 600 }}>
              Generated At
            </Typography>

            <Tooltip title={lastGenerated}>
              <Box>
                <Typography variant="body2" color="text.secondary">
                  {lastGenerated}
                </Typography>
              </Box>
            </Tooltip>
          </Stack>

          <Stack spacing={1}>
            <CopyableField label="Public Key" value={keys?.public} />
            <CopyableField label="Private Key" value={keys?.private} />
          </Stack>
        </Stack>
      </CardContent>
    </Card>
  );
}
