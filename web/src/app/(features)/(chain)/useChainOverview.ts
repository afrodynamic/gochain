import {
  buildOverview,
  deriveTransactions,
  fetchRecentBlocks,
  fetchRecentTransactions,
  type ChainOverview,
  type ChainTransactionRow,
} from '@/app/_utils/chainQueries';
import { useQuery } from '@tanstack/react-query';

export const useChainOverview = () =>
  useQuery<ChainOverview>({
    queryKey: ['chain', 'overview'],
    queryFn: async () => {
      const blocks = await fetchRecentBlocks();
      return buildOverview(blocks);
    },
    refetchInterval: 10_000,
  });

export const useRecentTransactions = () =>
  useQuery<ChainTransactionRow[]>({
    queryKey: ['chain', 'transactions'],
    queryFn: async () => {
      const transactions = await fetchRecentTransactions();

      return deriveTransactions(transactions);
    },
    refetchInterval: 10_000,
  });
