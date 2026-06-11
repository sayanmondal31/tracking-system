import { z } from "zod";

const envSchema = z.object({
  PORT: z.string().default("3001"),

  DATABASE_URL: z.string(),

  JWT_SECRET: z.string(),

  REDIS_URL: z.string(),
});

export const env = envSchema.parse(process.env);
