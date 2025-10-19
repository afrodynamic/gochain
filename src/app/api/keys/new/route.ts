import { KeyPairSchema } from '@/app/_utils/types';
import { NextRequest, NextResponse } from 'next/server';

export const POST = async (request: NextRequest) => {
  const base = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:9090';
  const url = `${base}/keys/new`;

  try {
    const seed = await request.text();
    const response = await fetch(url, {
      method: 'POST',
      body: seed,
      headers: {},
    });

    if (!response.ok) {
      try {
        const e = await response.json();

        return NextResponse.json(
          { data: null, error: e.error || 'keys/new failed' },
          { status: 400 }
        );
      } catch {
        return NextResponse.json(
          { data: null, error: 'keys/new failed' },
          { status: 400 }
        );
      }
    }

    const json = await response.json();
    const parsed = KeyPairSchema.safeParse(json);

    if (!parsed.success) {
      return NextResponse.json(
        { data: null, error: String(parsed.error) },
        { status: 200 }
      );
    }

    return NextResponse.json(parsed.data, { status: 200 });
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);

    return NextResponse.json(
      { data: null, error: 'keys/new proxy error: ' + msg },
      { status: 500 }
    );
  }
};
