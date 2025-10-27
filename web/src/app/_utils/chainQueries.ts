import { bytesToHex } from '@/app/_utils/encoding';
import type { BlockDetails, ChainTransaction } from '@/app/_utils/types';
import { chainClient } from '@/lib/rpc/clients';

export type ChainOverview = {
  blocks: BlockDetails[];
  latestBlock: BlockDetails | null;
  totalBlocks: number;
  totalTransactions: number;
};

export type ChainTransactionRow = {
  txHash: string;
  blockHeight: number;
  from: string;
  to: string;
  amount: string;
  fee: string;
  status: string;
  timestamp: Date;
};

const rpcBaseUrl = (
  process.env.NEXT_PUBLIC_RPC_URL || 'http://localhost:8080'
).replace(/\/$/, '');
const demoBaseUrl = `${rpcBaseUrl}/demo`;

type BlocksApiResponse = {
  blocks: Array<{
    hash: string;
    height: number;
    prevHash: string;
    timestamp: string;
    transactions: RawTransaction[];
  }>;
};

type RawTransaction = {
  hash: string;
  from: string;
  to: string;
  amount: number;
  fee: number;
  nonce: number;
  blockHeight: number;
  timestamp: string;
  status: string;
};

type TransactionsApiResponse = {
  transactions: RawTransaction[];
};

const parseTransaction = (raw: RawTransaction): ChainTransaction => ({
  hash: raw.hash,
  from: raw.from,
  to: raw.to,
  amount: raw.amount.toString(),
  fee: raw.fee.toString(),
  nonce: raw.nonce.toString(),
  blockHeight: raw.blockHeight,
  timestamp: new Date(raw.timestamp),
  status: raw.status,
});

const parseBlock = (
  raw: BlocksApiResponse['blocks'][number]
): BlockDetails => ({
  hash: raw.hash,
  height: raw.height,
  prevHash: raw.prevHash,
  timestamp: new Date(raw.timestamp),
  transactions: (raw.transactions ?? []).map(parseTransaction),
});

const toBlockDetails = (
  response: Awaited<ReturnType<typeof chainClient.getBlock>>
): BlockDetails => ({
  height: Number(response.height),
  hash: bytesToHex(response.hash),
  prevHash: bytesToHex(response.prevHash),
  timestamp: new Date(),
  transactions: [],
});

export const fetchBlockAt = async (height: number): Promise<BlockDetails> => {
  if (!Number.isFinite(height) || height < 0) {
    throw new Error('Height must be non-negative');
  }

  const response = await chainClient.getBlock({ height: BigInt(height) });

  return toBlockDetails(response);
};

const fetchJson = async <T>(path: string): Promise<T> => {
  const response = await fetch(path, {
    headers: { 'Content-Type': 'application/json' },
    cache: 'no-store',
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error(`Request failed with status ${response.status}`);
  }

  return (await response.json()) as T;
};

export const fetchRecentBlocks = async (
  limit = 64
): Promise<BlockDetails[]> => {
  const data = await fetchJson<BlocksApiResponse>(
    `${demoBaseUrl}/blocks?limit=${limit}`
  );

  return data.blocks.map(parseBlock);
};

export const fetchRecentTransactions = async (
  limit = 128
): Promise<ChainTransaction[]> => {
  const data = await fetchJson<TransactionsApiResponse>(
    `${demoBaseUrl}/transactions?limit=${limit}`
  );

  return data.transactions.map(parseTransaction);
};

export const buildOverview = (blocks: BlockDetails[]): ChainOverview => {
  const latestBlock = blocks.at(-1) ?? null;
  const totalBlocks = blocks.length;
  const totalTransactions = blocks.reduce(
    (count, block) => count + block.transactions.length,
    0
  );

  return {
    blocks,
    latestBlock,
    totalBlocks,
    totalTransactions,
  };
};

export const deriveTransactions = (
  transactions: ChainTransaction[]
): ChainTransactionRow[] =>
  [...transactions].reverse().map((tx) => ({
    txHash: tx.hash,
    blockHeight: tx.blockHeight,
    from: tx.from,
    to: tx.to,
    amount: tx.amount,
    fee: tx.fee,
    status: tx.status,
    timestamp: tx.timestamp,
  }));
