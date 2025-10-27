'use client';
import { useChainOverview } from '@/app/(features)/(chain)/useChainOverview';
import {
  MaterialReactTable,
  type MRT_ColumnDef,
} from '@/app/_components/table/MaterialReactTable';
import type { BlockDetails } from '@/app/_utils/types';
import {
  Alert,
  Box,
  CircularProgress,
  Container,
  Grid,
  Stack,
  Typography,
} from '@mui/material';
import { useMemo } from 'react';

type BlockRow = BlockDetails & { index: number; txCount: number };

const columns: MRT_ColumnDef<BlockRow>[] = [
  {
    accessorKey: 'height',
    header: 'Height',
    size: 120,
  },
  {
    accessorKey: 'timestamp',
    header: 'Timestamp',
    Cell: ({ row }) => row.timestamp.toLocaleString(),
  },
  {
    accessorKey: 'txCount',
    header: 'Transactions',
    size: 140,
  },
  {
    accessorKey: 'hash',
    header: 'Block Hash',
    Cell: ({ row }) => row.hash,
  },
  {
    accessorKey: 'prevHash',
    header: 'Previous Hash',
    Cell: ({ row }) => row.prevHash || '—',
  },
];

export function BlockExplorer() {
  const { data, isLoading, error } = useChainOverview();

  const rows = useMemo(() => {
    if (!data) {
      return [] as BlockRow[];
    }

    return [...data.blocks].reverse().map((block, index) => ({
      ...block,
      index,
      txCount: block.transactions.length,
    }));
  }, [data]);

  return (
    <Box>
      <Container maxWidth="lg" sx={{ py: { xs: 10, md: 12 } }}>
        <Stack spacing={4}>
          <Grid container alignItems="center" spacing={2}>
            <Grid size={{ xs: 12, sm: 8 }}>
              <Typography variant="h3" fontWeight={800} gutterBottom>
                Block Explorer
              </Typography>

              <Typography variant="body1">
                Browse the most recent blocks, including timestamps and
                transaction activity.
              </Typography>
            </Grid>
          </Grid>

          {error && <Alert severity="error">{error.message}</Alert>}

          {isLoading ? (
            <Box py={6} textAlign="center">
              <CircularProgress />
            </Box>
          ) : (
            <MaterialReactTable
              columns={columns}
              data={rows}
              emptyContent="No blocks discovered yet"
            />
          )}
        </Stack>
      </Container>
    </Box>
  );
}
