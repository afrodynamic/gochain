import { walletClient } from '@/lib/rpc/clients';
import { ConnectError } from '@connectrpc/connect';
import { useMutation, useQueryClient } from '@tanstack/react-query';

const toError = (error: unknown, fallback: string) => {
  if (error instanceof ConnectError) {
    return new Error(error.rawMessage || fallback);
  }

  if (error instanceof Error) {
    return error;
  }

  return new Error(fallback);
};

type WalletKeys = {
  address: string;
  privateKey: string;
  publicKey: string;
};

const encodeSeed = (seed?: string): Uint8Array => {
  if (!seed) {
    return new Uint8Array();
  }

  return new TextEncoder().encode(seed);
};

const createWallet = async ({
  seed,
}: {
  seed?: string;
}): Promise<WalletKeys> => {
  try {
    const response = await walletClient.newKey({ seed: encodeSeed(seed) });

    return {
      address: response.addr,
      privateKey: response.priv,
      publicKey: response.pub,
    };
  } catch (error) {
    throw toError(error, 'Unable to create wallet');
  }
};

type BalanceResult = {
  address: string;
  balance: string;
};

const lookupBalance = async ({
  address,
}: {
  address: string;
}): Promise<BalanceResult> => {
  if (!address.trim()) {
    throw new Error('Address is required');
  }

  try {
    const response = await walletClient.balance({ addr: address });

    return {
      address,
      balance: response.balance.toString(),
    };
  } catch (error) {
    throw toError(error, 'Unable to fetch balance');
  }
};

type WalletTransferInput = {
  from: string;
  to: string;
  amount: string;
  privateKey: string;
  fee?: string;
};

type WalletTransferResult = {
  txId: string;
  status: string;
};

const sendWalletTx = async ({
  from,
  to,
  amount,
  privateKey,
  fee,
}: WalletTransferInput): Promise<WalletTransferResult> => {
  if (!from.trim() || !to.trim()) {
    throw new Error('Both sender and recipient are required');
  }

  if (!privateKey.trim()) {
    throw new Error('Private key is required to sign the transaction');
  }

  const normalizedAmount = amount.trim();

  if (!normalizedAmount) {
    throw new Error('Amount is required');
  }

  let amountBigInt: bigint;

  try {
    amountBigInt = BigInt(normalizedAmount);
  } catch {
    throw new Error('Amount must be an integer value');
  }

  if (amountBigInt <= 0) {
    throw new Error('Amount must be positive');
  }

  let feeBigInt: bigint = BigInt(1);

  if (fee && fee.trim()) {
    try {
      feeBigInt = BigInt(fee.trim());
    } catch {
      throw new Error('Fee must be a valid integer');
    }
  }

  try {
    const built = await walletClient.buildTx({
      from,
      to,
      amount: amountBigInt,
      feeHint: { maxFeePerGas: feeBigInt, maxPriorityFee: feeBigInt },
    });

    const signed = await walletClient.signTx({
      priv: privateKey,
      tx: built.tx,
    });

    const broadcast = await walletClient.broadcast({ signed: signed.signed });

    const status = await walletClient.txStatus({ txId: broadcast.txId });

    return {
      txId: broadcast.txId,
      status: status.status,
    };
  } catch (error) {
    throw toError(error, 'Unable to send transaction');
  }
};

export const useCreateWallet = () =>
  useMutation<WalletKeys, Error, { seed?: string }>({
    mutationFn: createWallet,
  });

export const useWalletBalance = () =>
  useMutation<BalanceResult, Error, { address: string }>({
    mutationFn: lookupBalance,
  });

export const useSendWalletTransaction = () => {
  const queryClient = useQueryClient();

  return useMutation<WalletTransferResult, Error, WalletTransferInput>({
    mutationFn: sendWalletTx,
    onSuccess: async () => {
      await Promise.all([
        queryClient.invalidateQueries({ queryKey: ['chain', 'overview'] }),
        queryClient.invalidateQueries({ queryKey: ['chain', 'transactions'] }),
      ]);
    },
  });
};
