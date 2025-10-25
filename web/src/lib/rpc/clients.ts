import { Chain } from '@/gen/chain/v1/chain_pb';
import { Wallet } from '@/gen/wallet/v1/wallet_pb';
import { createClient } from '@connectrpc/connect';
import { createGrpcWebTransport } from '@connectrpc/connect-web';

const transport = createGrpcWebTransport({
  baseUrl: process.env.NEXT_PUBLIC_RPC_URL || 'http://localhost:8080',
});

export const chainClient = createClient(Chain, transport);
export const walletClient = createClient(Wallet, transport);
