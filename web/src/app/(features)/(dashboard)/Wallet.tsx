'use client';
import { CopyableField } from '@/app/_components/core/CopyableField';
import { MetamaskConnectButton } from '@/app/_utils/MetamaskConnectButton';
import { OpenInNew } from '@mui/icons-material';
import {
  Button,
  Card,
  CardContent,
  CardHeader,
  Chip,
  LinearProgress,
  MenuItem,
  Stack,
  TextField,
  Tooltip,
  Typography,
} from '@mui/material';
import { useMemo, useState } from 'react';
import { parseEther } from 'viem';
import {
  useAccount,
  useBalance,
  useChainId,
  useChains,
  useConnect,
  useSendTransaction,
  useSignMessage,
  useSwitchChain,
} from 'wagmi';

export function Wallet() {
  const { isPending: isConnecting } = useConnect();
  const { address, isConnected } = useAccount();
  const chains = useChains();
  const chainId = useChainId();

  const { switchChain, isPending: isSwitching } = useSwitchChain();

  const {
    data: bal,
    status: balStatus,
    refetch: refetchBal,
  } = useBalance({ address, chainId, query: { enabled: isConnected } });

  const { signMessageAsync, status: sigStatus } = useSignMessage();
  const {
    sendTransactionAsync,
    status: sendStatus,
    isPending: isSending,
  } = useSendTransaction();

  const [amount, setAmount] = useState('0.001');
  const [lastSig, setLastSig] = useState('');
  const [lastTx, setLastTx] = useState('');

  const activeChain = useMemo(
    () => chains.find((c) => c.id === chainId) ?? null,
    [chainId, chains]
  );

  return (
    <Card elevation={6} className="max-w-lg w-full">
      {(isConnecting || isSending || isSwitching) && <LinearProgress />}

      <CardHeader
        title="Wallet"
        subheader="MetaMask"
        sx={{ '& .MuiCardHeader-title': { fontWeight: 700, fontSize: 22 } }}
        action={<MetamaskConnectButton />}
      />

      <CardContent>
        <Stack spacing={2}>
          {!isConnected ? (
            <Typography variant="body1">
              Connect a wallet to view details.
            </Typography>
          ) : (
            <>
              <CopyableField label="Wallet Address" value={address || ''} />

              <Stack direction="row" spacing={1} alignItems="center">
                <TextField
                  select
                  label="Network"
                  size="small"
                  value={activeChain?.id ?? ''}
                  onChange={(e) =>
                    switchChain({ chainId: Number(e.target.value) })
                  }
                  sx={{ minWidth: 220 }}
                >
                  {chains.map((c) => (
                    <MenuItem key={c.id} value={c.id}>
                      {c.name}
                    </MenuItem>
                  ))}
                </TextField>

                {activeChain && (
                  <Tooltip title="Obtain faucet funds" arrow placement="top">
                    <Button
                      href="https://cloud.google.com/application/web3/faucet"
                      target="_blank"
                      rel="noreferrer"
                      endIcon={<OpenInNew />}
                    >
                      Faucet
                    </Button>
                  </Tooltip>
                )}
              </Stack>

              <Stack direction="row" spacing={1} alignItems="center">
                <Typography variant="body1" sx={{ fontWeight: 600 }}>
                  Balance
                </Typography>
                <Chip
                  label={
                    balStatus === 'pending'
                      ? 'Loading…'
                      : bal
                      ? `${bal.formatted} ${bal.symbol}`
                      : '—'
                  }
                />
                <Button size="small" onClick={() => refetchBal()}>
                  Refresh
                </Button>
              </Stack>

              <Button
                variant="outlined"
                onClick={async () => {
                  const sig = await signMessageAsync({
                    message: `Hello from ${address}`,
                  });
                  setLastSig(sig);
                }}
                disabled={sigStatus === 'pending'}
              >
                Sign Message
              </Button>

              {lastSig && (
                <Typography variant="body2" color="text.secondary">
                  Sig: {lastSig.slice(0, 20)}…
                </Typography>
              )}

              <Stack direction="row" spacing={1}>
                <TextField
                  label="Amount (native)"
                  size="small"
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                />
                <Button
                  variant="contained"
                  onClick={async () => {
                    const tx = await sendTransactionAsync({
                      to: address!,
                      value: parseEther(amount || '0'),
                    });
                    setLastTx(tx);
                  }}
                  disabled={sendStatus === 'pending'}
                  size="small"
                >
                  Send To Self
                </Button>
              </Stack>

              {lastTx && (
                <Typography variant="body2" color="text.secondary">
                  Tx: {lastTx}
                </Typography>
              )}
            </>
          )}
        </Stack>
      </CardContent>
    </Card>
  );
}
