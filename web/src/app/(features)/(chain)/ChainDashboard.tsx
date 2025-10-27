'use client';
import { MetricCard } from '@/app/(features)/(chain)/MetricCard';
import { useChainOverview } from '@/app/(features)/(chain)/useChainOverview';
import { Health } from '@/app/(features)/(dashboard)/Health';
import { MetamaskWallet } from '@/app/(features)/(dashboard)/MetamaskWallet';
import RefreshIcon from '@mui/icons-material/Refresh';
import Masonry from '@mui/lab/Masonry';
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  CircularProgress,
  Container,
  LinearProgress,
  List,
  ListItem,
  ListItemText,
  Stack,
  Typography,
} from '@mui/material';

export function ChainDashboard() {
  const { data, isLoading, isFetching, error, refetch } = useChainOverview();

  const latestBlock = data?.latestBlock ?? null;
  const latestTransactions = latestBlock?.transactions ?? [];

  return (
    <Box>
      <Container maxWidth="lg" sx={{ py: { xs: 10, md: 12 } }}>
        <Stack spacing={4}>
          <Stack
            direction={{ xs: 'column', sm: 'row' }}
            spacing={2}
            alignItems={{ xs: 'flex-start', sm: 'center' }}
            className="w-full max-w-7xl"
            maxWidth="lg"
            padding={2}
          >
            <Box>
              <Typography variant="h3" fontWeight={800} gutterBottom>
                Gochain Dashboard
              </Typography>

              <Typography variant="body1">
                Inspect the latest gochain state at a glance.
              </Typography>
            </Box>

            <Box flexGrow={1} />

            <Button
              variant="outlined"
              startIcon={
                isFetching ? <CircularProgress size={16} /> : <RefreshIcon />
              }
              onClick={() => refetch()}
              disabled={isFetching}
            >
              Refresh
            </Button>
          </Stack>

          {error && <Alert severity="error">{error.message}</Alert>}

          <Masonry
            columns={{ xs: 1, md: 2, lg: 3 }}
            spacing={2}
            className="w-full p-2 max-w-7xl mx-auto"
          >
            <Box className="max-lg:flex max-lg:justify-center">
              <Health />
            </Box>

            <Box className="max-lg:flex max-lg:justify-center">
              <MetamaskWallet />
            </Box>

            <Box className="max-lg:flex max-lg:justify-center">
              <Card
                elevation={6}
                className="max-w-lg w-full border"
                sx={{ borderColor: 'primary.main' }}
              >
                {isFetching && <LinearProgress />}

                <CardHeader
                  title="Latest Block Details"
                  subheader="Raw values returned by the chain RPC"
                  sx={{ '& .MuiCardHeader-title': { fontWeight: 700 } }}
                />

                <CardContent>
                  {isLoading ? (
                    <Box py={4} textAlign="center">
                      <CircularProgress />
                    </Box>
                  ) : latestBlock ? (
                    <List>
                      <ListItem divider>
                        <ListItemText
                          primary="Height"
                          secondary={latestBlock.height}
                        />
                      </ListItem>

                      <ListItem divider>
                        <ListItemText
                          primary="Hash"
                          secondary={latestBlock.hash}
                        />
                      </ListItem>

                      <ListItem divider>
                        <ListItemText
                          primary="Timestamp"
                          secondary={latestBlock.timestamp.toLocaleString()}
                        />
                      </ListItem>

                      <ListItem>
                        <ListItemText
                          primary="Previous Hash"
                          secondary={latestBlock.prevHash || '—'}
                        />
                      </ListItem>

                      <ListItem>
                        <ListItemText
                          primary="Transactions"
                          secondary={
                            latestTransactions.length
                              ? latestTransactions
                                  .map((tx) => `${tx.hash.slice(0, 12)}…`)
                                  .slice(0, 3)
                                  .join(', ')
                              : 'None'
                          }
                        />
                      </ListItem>
                    </List>
                  ) : (
                    <Typography variant="body2" color="text.secondary">
                      No block data available yet. Submit a transaction to seal
                      a new block.
                    </Typography>
                  )}
                </CardContent>
              </Card>
            </Box>

            <Box className="max-lg:flex max-lg:justify-center">
              <MetricCard
                label="Latest Block"
                value={
                  isLoading
                    ? 'Loading…'
                    : latestBlock
                    ? `#${latestBlock.height}`
                    : '—'
                }
                helper={
                  latestBlock
                    ? `${
                        latestTransactions.length
                      } tx · ${latestBlock.timestamp.toLocaleString()}`
                    : undefined
                }
              />
            </Box>

            <Box className="max-lg:flex max-lg:justify-center">
              <MetricCard
                label="Total Blocks"
                value={isLoading ? 'Loading…' : data?.totalBlocks ?? '—'}
                helper="Includes the genesis block"
              />
            </Box>

            <Box className="max-lg:flex max-lg:justify-center">
              <MetricCard
                label="Total Transactions"
                value={isLoading ? 'Loading…' : data?.totalTransactions ?? '—'}
                helper="Counted across mined blocks"
              />
            </Box>
          </Masonry>
        </Stack>
      </Container>
    </Box>
  );
}
