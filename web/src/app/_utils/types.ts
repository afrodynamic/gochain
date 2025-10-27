import { z } from 'zod';

export const HealthStatusResponseSchema = z.object({
  status: z.enum(['ok', 'down']),
  timestamp: z.string().transform((val) => new Date(val)),
});

export type HealthStatusResponse = z.infer<typeof HealthStatusResponseSchema>;

export const KeyPairSchema = z.object({
  address: z.string(),
  publicKey: z.string(),
  privateKey: z.string(),
});
export type KeyPair = z.infer<typeof KeyPairSchema>;

export type ChainTransaction = {
  hash: string;
  from: string;
  to: string;
  amount: string;
  fee: string;
  nonce: string;
  blockHeight: number;
  timestamp: Date;
  status: string;
};

export type BlockDetails = {
  height: number;
  hash: string;
  prevHash: string;
  timestamp: Date;
  transactions: ChainTransaction[];
};

export type AccountBalance = {
  address: string;
  balance: string;
};

export type SubmittedTx = {
  txHash: string;
  status?: string;
};
