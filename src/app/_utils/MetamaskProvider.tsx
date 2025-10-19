'use client';
import { ReactNode } from 'react';
import { createConfig, http, WagmiProvider } from 'wagmi';
import {
  arbitrumSepolia,
  lineaSepolia,
  polygonAmoy,
  sepolia,
} from 'wagmi/chains';
import { metaMask } from 'wagmi/connectors';

const wagmiConfig = createConfig({
  ssr: true,
  chains: [sepolia, arbitrumSepolia, polygonAmoy, lineaSepolia],
  connectors: [
    metaMask({
      infuraAPIKey: process.env.NEXT_PUBLIC_INFURA_API_KEY,
    }),
  ],
  transports: {
    [sepolia.id]: http(sepolia.rpcUrls.default.http[0]),
    [arbitrumSepolia.id]: http(arbitrumSepolia.rpcUrls.default.http[0]),
    [polygonAmoy.id]: http(polygonAmoy.rpcUrls.default.http[0]),
    [lineaSepolia.id]: http(lineaSepolia.rpcUrls.default.http[0]),
  },
});

export const MetamaskProvider = ({ children }: { children: ReactNode }) => (
  <WagmiProvider config={wagmiConfig}>{children}</WagmiProvider>
);
