import { z } from 'zod';

export const HealthStatusResponseSchema = z.object({
  status: z.enum(['ok', 'down']),
  timestamp: z.string().transform((val) => new Date(val)),
});

export type HealthStatusResponse = z.infer<typeof HealthStatusResponseSchema>;

export const KeyPairSchema = z.object({
  public: z.string(),
  private: z.string(),
});
export type KeyPair = z.infer<typeof KeyPairSchema>;
