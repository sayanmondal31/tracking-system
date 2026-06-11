import { createClient } from "redis";
import { env } from "./env";
import { AppError } from "../errors/AppError";

export const redis = createClient({
  url: env.REDIS_URL,
});

redis.on("error", (err) => {
  console.error("Redis Error:", err);
});

export async function connectRedis() {
  await redis.connect();
}
