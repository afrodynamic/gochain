'use client';
import { useRecentTransactions } from '@/app/(features)/(chain)/useChainOverview';
import {
  MaterialReactTable,
  type MRT_ColumnDef,
} from '@/app/_components/table/MaterialReactTable';
import { type ChainTransactionRow } from '@/app/_utils/chainQueries';
import {
  Alert,
  Box,
  Chip,
  CircularProgress,
  Container,
  Grid,
  Stack,
  Typography,
} from '@mui/material';
import { useMemo } from 'react';

const columns: MRT_ColumnDef<ChainTransactionRow>[] = [
  {
    accessorKey: 'txHash',
    header: 'Tx Hash',
    Cell: ({ row }) => row.txHash,
  },
  {
    accessorKey: 'blockHeight',
    header: 'Block',
    size: 100,
  },
  {
    accessorKey: 'from',
    header: 'From',
  },
  {
    accessorKey: 'to',
    header: 'To',
  },
  {
    accessorKey: 'amount',
    header: 'Amount',
    size: 120,
  },
  {
    accessorKey: 'fee',
    header: 'Fee',
    size: 120,
  },
  {
    accessorKey: 'timestamp',
    header: 'Timestamp',
    Cell: ({ row }) => row.timestamp.toLocaleString(),
  },
  {
    accessorKey: 'status',
    header: 'Status',
    size: 120,
    Cell: ({ row }) => <Chip label={row.status} color="success" size="small" />,
  },
];

export function TransactionsTable() {
  const { data, isLoading, error } = useRecentTransactions();

  const rows = useMemo(() => data ?? [], [data]);

  return (
    <Box>
      <Container maxWidth="lg" sx={{ py: { xs: 10, md: 12 } }}>
        <Stack spacing={4}>
          <Grid container alignItems="center" spacing={2}>
            <Grid size={{ xs: 12, sm: 8 }}>
              <Typography variant="h3" fontWeight={800} gutterBottom>
                Recent Transactions
              </Typography>

              <Typography variant="body1">
                View mined transfers as they are sealed into blocks with
                timestamps and fee details.
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
              emptyContent="No transactions yet"
            />
          )}
        </Stack>
      </Container>
    </Box>
  );
}
