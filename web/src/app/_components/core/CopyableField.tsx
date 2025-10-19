import { Check, ContentCopy } from '@mui/icons-material';
import { Box, Stack, Tooltip, Typography } from '@mui/material';
import IconButton from '@mui/material/IconButton';
import { useState } from 'react';

export function CopyableField({
  label,
  value,
}: {
  label: string;
  value?: string;
}) {
  const [copied, setCopied] = useState(false);
  const handleCopy = async () => {
    if (!value) return;
    await navigator.clipboard.writeText(value);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  return (
    <Stack spacing={0.5}>
      <Stack direction="row" justifyContent="space-between" alignItems="center">
        <Typography variant="subtitle2">{label}</Typography>

        <Tooltip title={copied ? 'Copied!' : 'Copy'}>
          <IconButton size="small" onClick={handleCopy}>
            {copied ? <Check color="success" /> : <ContentCopy />}
          </IconButton>
        </Tooltip>
      </Stack>

      <Box
        sx={{
          p: 1,
          bgcolor: 'action.hover',
          borderRadius: 1,
          fontFamily: 'monospace',
          overflowX: 'auto',
        }}
      >
        {value ?? '—'}
      </Box>
    </Stack>
  );
}
