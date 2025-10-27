'use client';
import {
  useCreateWallet,
  useSendWalletTransaction,
  useWalletBalance,
} from '@/app/(features)/(chain)/useWallet';
import { CopyableField } from '@/app/_components/core/CopyableField';
import {
  Alert,
  Box,
  Button,
  Card,
  CardContent,
  CardHeader,
  Container,
  Divider,
  Grid,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { useEffect, useState } from 'react';

type WalletForm = {
  address: string;
  privateKey: string;
  publicKey: string;
};

export function WalletDashboard() {
  const [seed, setSeed] = useState('');
  const [wallet, setWallet] = useState<WalletForm | null>(null);
  const [balanceAddress, setBalanceAddress] = useState('');
  const [transferForm, setTransferForm] = useState({
    from: '',
    to: '',
    amount: '1',
    privateKey: '',
    fee: '1',
  });

  const createWallet = useCreateWallet();
  const walletBalance = useWalletBalance();
  const sendTx = useSendWalletTransaction();

  useEffect(() => {
    if (wallet) {
      setBalanceAddress(wallet.address);
      setTransferForm((prev) => ({
        ...prev,
        from: wallet.address,
        privateKey: wallet.privateKey,
      }));
    }
  }, [wallet]);

  return (
    <Box>
      <Container maxWidth="md" sx={{ py: { xs: 10, md: 12 } }}>
        <Stack spacing={4}>
          <Typography variant="h3" fontWeight={800}>
            Wallet
          </Typography>

          <Typography variant="body1" color="text.secondary">
            Generate keys, inspect balances, and broadcast transactions on the
            gochain demo network.
          </Typography>

          <Card
            elevation={6}
            className="border"
            sx={{ borderColor: 'primary.main' }}
          >
            <CardHeader
              title="Create a Wallet"
              subheader="Optionally provide a seed to deterministically derive a key"
              sx={{ '& .MuiCardHeader-title': { fontWeight: 700 } }}
            />

            <CardContent>
              <Stack spacing={2}>
                <TextField
                  label="Seed (optional)"
                  value={seed}
                  onChange={(event) => setSeed(event.target.value)}
                  placeholder="Type any text to derive a deterministic wallet"
                />

                {createWallet.error && (
                  <Alert severity="error">{createWallet.error.message}</Alert>
                )}

                <Box
                  alignItems={{ xs: 'stretch', sm: 'center' }}
                  display="flex"
                  flexDirection={{ xs: 'column', sm: 'row' }}
                >
                  <Button
                    variant="contained"
                    size="large"
                    onClick={async () => {
                      try {
                        const result = await createWallet.mutateAsync({
                          seed: seed.trim() || undefined,
                        });
                        setWallet(result);
                      } catch {
                        /* handled by mutation state */
                      }
                    }}
                    disabled={createWallet.status === 'pending'}
                  >
                    {createWallet.status === 'pending'
                      ? 'Creating…'
                      : 'Generate Wallet'}
                  </Button>
                </Box>

                {wallet && (
                  <Stack spacing={1.5}>
                    <CopyableField label="Address" value={wallet.address} />

                    <CopyableField
                      label="Public Key"
                      value={wallet.publicKey}
                    />

                    <CopyableField
                      label="Private Key"
                      value={wallet.privateKey}
                    />
                  </Stack>
                )}
              </Stack>
            </CardContent>
          </Card>

          <Card
            elevation={6}
            className="border"
            sx={{ borderColor: 'primary.main' }}
          >
            <CardHeader
              title="Check Balance"
              subheader="Query wallet adapter balances for any address"
              sx={{ '& .MuiCardHeader-title': { fontWeight: 700 } }}
            />

            <CardContent>
              <Stack spacing={2}>
                <TextField
                  label="Address"
                  value={balanceAddress}
                  onChange={(event) => setBalanceAddress(event.target.value)}
                  fullWidth
                />

                <Box
                  alignItems={{ xs: 'stretch', sm: 'center' }}
                  display="flex"
                  flexDirection={{ xs: 'column', sm: 'row' }}
                >
                  <Button
                    size="large"
                    variant="contained"
                    onClick={async () => {
                      try {
                        await walletBalance.mutateAsync({
                          address: balanceAddress,
                        });
                      } catch {
                        /* handled by mutation state */
                      }
                    }}
                    disabled={walletBalance.status === 'pending'}
                  >
                    {walletBalance.status === 'pending'
                      ? 'Checking…'
                      : 'Check Balance'}
                  </Button>
                </Box>

                {walletBalance.error && (
                  <Alert severity="error">{walletBalance.error.message}</Alert>
                )}

                {walletBalance.data && (
                  <Alert severity="success">
                    Balance for {walletBalance.data.address}:{' '}
                    {walletBalance.data.balance}
                  </Alert>
                )}
              </Stack>
            </CardContent>
          </Card>

          <Card
            elevation={6}
            className="border"
            sx={{ borderColor: 'primary.main' }}
          >
            <CardHeader
              title="Send Transaction"
              subheader="Build, sign, and broadcast a gochain transfer"
              sx={{ '& .MuiCardHeader-title': { fontWeight: 700 } }}
            />

            <CardContent>
              <Stack spacing={2}>
                <Grid container spacing={2} alignItems="stretch">
                  <Grid size={{ xs: 12, sm: 6 }}>
                    <TextField
                      label="From"
                      value={transferForm.from}
                      onChange={(event) =>
                        setTransferForm((prev) => ({
                          ...prev,
                          from: event.target.value,
                        }))
                      }
                      fullWidth
                    />
                  </Grid>

                  <Grid size={{ xs: 12, sm: 6 }}>
                    <TextField
                      label="To"
                      value={transferForm.to}
                      onChange={(event) =>
                        setTransferForm((prev) => ({
                          ...prev,
                          to: event.target.value,
                        }))
                      }
                      fullWidth
                    />
                  </Grid>
                </Grid>

                <Grid container spacing={2} alignItems="stretch">
                  <Grid size={{ xs: 12, sm: 4 }}>
                    <TextField
                      label="Amount"
                      value={transferForm.amount}
                      onChange={(event) =>
                        setTransferForm((prev) => ({
                          ...prev,
                          amount: event.target.value,
                        }))
                      }
                      fullWidth
                    />
                  </Grid>

                  <Grid size={{ xs: 12, sm: 4 }}>
                    <TextField
                      label="Fee"
                      value={transferForm.fee}
                      onChange={(event) =>
                        setTransferForm((prev) => ({
                          ...prev,
                          fee: event.target.value,
                        }))
                      }
                      fullWidth
                    />
                  </Grid>
                </Grid>

                {sendTx.error && (
                  <Alert severity="error">{sendTx.error.message}</Alert>
                )}

                <Box
                  alignItems={{ xs: 'stretch', sm: 'center' }}
                  display="flex"
                  flexDirection={{ xs: 'column', sm: 'row' }}
                >
                  <Button
                    variant="contained"
                    size="large"
                    onClick={async () => {
                      try {
                        await sendTx.mutateAsync({
                          ...transferForm,
                        });
                        setTransferForm((prev) => ({ ...prev, amount: '1' }));
                      } catch {
                        /* handled by mutation state */
                      }
                    }}
                    disabled={sendTx.status === 'pending'}
                  >
                    {sendTx.status === 'pending' ? 'Submitting…' : 'Submit'}
                  </Button>
                </Box>

                {sendTx.data && (
                  <Alert
                    severity="success"
                    sx={{ width: '100%', wordBreak: 'break-word' }}
                  >
                    <Stack spacing={0.5}>
                      <Typography component="span">
                        Transaction submitted successfully.
                      </Typography>

                      <Typography component="span" variant="body2">
                        Hash: <code>{sendTx.data.txId}</code>
                      </Typography>

                      <Typography component="span" variant="body2">
                        Status: {sendTx.data.status}
                      </Typography>
                    </Stack>
                  </Alert>
                )}
              </Stack>
            </CardContent>

            <Divider />

            <CardContent>
              <Typography variant="caption" color="text.secondary">
                Tip: balances in this demo chain are simulated. Create a few
                wallets to seed them with funds and experiment with transfers
                between addresses.
              </Typography>
            </CardContent>
          </Card>
        </Stack>
      </Container>
    </Box>
  );
}
