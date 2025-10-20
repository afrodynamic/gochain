import type { KeyPair } from '@/app/_utils/types';
import { useMutation } from '@tanstack/react-query';

const createKeys = async (opts?: {
  seed?: string;
  passphrase?: string;
  mode?: 'random' | 'deterministic';
}): Promise<KeyPair> => {
  const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:9090';
  const body = JSON.stringify(opts ?? { mode: 'random' });
  const response = await fetch(`${base}/keys/new`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body,
  });

  if (!response.ok) {
    const text = await response.text();
    throw new Error(`Failed to create keys: ${text}`);
  }

  return response.json();
};

export const useKeysNew = () =>
  useMutation<
    KeyPair,
    Error,
    { seed?: string; passphrase?: string; mode?: 'random' | 'deterministic' }
  >({
    mutationFn: ({ seed, passphrase }) =>
      createKeys({ seed, passphrase, mode: seed ? 'deterministic' : 'random' }),
  });
