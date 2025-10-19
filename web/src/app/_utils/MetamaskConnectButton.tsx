'use client';
import { Box, Button } from '@mui/material';
import { useAccount, useConnect, useDisconnect } from 'wagmi';

export const MetamaskConnectButton = () => {
  const { address } = useAccount();
  const { connectors, connect } = useConnect();
  const { disconnect } = useDisconnect();

  return (
    <Box>
      {address ? (
        <Button
          variant="contained"
          color="error"
          size="small"
          onClick={() => disconnect()}
        >
          Disconnect
        </Button>
      ) : (
        <Box display="flex" gap={1}>
          {connectors.map((connector) => (
            <Button
              key={connector.uid}
              variant="contained"
              color="success"
              onClick={() => connect({ connector })}
              className="font-bold"
            >
              Connect {connector.name}
            </Button>
          ))}
        </Box>
      )}
    </Box>
  );
};
