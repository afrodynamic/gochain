import { HealthStatusResponseSchema } from '@/app/_utils/types';
import { NextResponse } from 'next/server';

export const GET = async () => {
  const baseUrl =
    process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';
  const healthUrl = `${baseUrl}/health`;

  try {
    const healthResponse = await fetch(healthUrl);

    if (healthResponse.ok) {
      const healthData = await healthResponse.json();
      const parsed = HealthStatusResponseSchema.safeParse(healthData);

      if (parsed.success) {
        console.log('Health data:', parsed.data);
        return NextResponse.json(parsed.data, { status: 200 });
      }

      return NextResponse.json(
        {
          data: {
            status: 'down',
            timestamp: new Date().toISOString(),
          },
        },
        { status: 200 }
      );
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : String(error);
    return NextResponse.json(
      {
        data: null,
        error: 'Failed to fetch health status: ' + message,
      },
      { status: 500 }
    );
  }
};
