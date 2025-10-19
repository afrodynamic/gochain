import type { KeyPair } from '@/app/_utils/types';
import { useMutation } from '@tanstack/react-query';

const createKeys = async (seed?: string): Promise<KeyPair> => {
  const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:9090';
  const response = await fetch(`${base}/keys/new`, {
    method: 'POST',
    body: seed ?? '',
  });

  if (!response.ok) {
    throw new Error('Failed to create keys');
  }

  const json = await response.json();

  return json;
};

export const useKeysNew = () =>
  useMutation<KeyPair, Error, { seed?: string }>({
    mutationFn: ({ seed }) => createKeys(seed),
  });
