import type { HealthStatusResponse } from '@/app/_utils/types';
import { useQuery } from '@tanstack/react-query';

const fetchHealthStatus = async (): Promise<HealthStatusResponse> => {
  try {
    const baseUrl =
      process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:9090';
    const apiUrl = `${baseUrl}/health`;
    const response = await fetch(apiUrl);

    if (!response.ok) {
      throw new Error('Failed to fetch health status');
    }

    const data = await response.json();

    return {
      status: data.status,
      timestamp: data.timestamp,
    };
  } catch (error) {
    console.error(error);
    throw new Error('Failed to fetch health status');
  }
};

export const useHealthStatusQuery = () => {
  return useQuery({
    queryKey: ['healthStatus'],
    queryFn: fetchHealthStatus,
  });
};
