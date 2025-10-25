import type { KeyPair } from '@/app/_utils/types';
import { walletClient } from '@/lib/rpc/clients';
import { useMutation } from '@tanstack/react-query';

const createKeys = async (opts?: {
  seed?: string;
  mode?: 'random' | 'deterministic';
}): Promise<KeyPair> => {
  const seed =
    opts?.mode === 'deterministic' && opts?.seed
      ? new TextEncoder().encode(opts.seed)
      : undefined;
  const response = await walletClient.newKey({ seed });

  return {
    privateKey: response.priv,
    publicKey: response.pub,
    address: response.addr,
  };
};

export const useKeysNew = () =>
  useMutation<
    KeyPair,
    Error,
    { seed?: string; mode?: 'random' | 'deterministic' }
  >({
    mutationFn: ({ seed }) =>
      createKeys({ seed, mode: seed ? 'deterministic' : 'random' }),
  });
