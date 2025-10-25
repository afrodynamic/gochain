import type { HealthStatusResponse } from '@/app/_utils/types';
import { useQuery } from '@tanstack/react-query';

const fetchHealthStatus = async (): Promise<HealthStatusResponse> => {
  const baseUrl =
    process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
  const response = await fetch(`${baseUrl}/health`, { cache: 'no-store' });

  if (!response.ok) {
    throw new Error('health check failed');
  }

  const ts = new Date();

  return {
    status: 'ok',
    timestamp: ts,
  };
};

export const useHealthStatusQuery = () =>
  useQuery({
    queryKey: ['healthStatus'],
    queryFn: fetchHealthStatus,
    staleTime: 5_000,
  });
