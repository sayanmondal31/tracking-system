import app from "./app";
import { db } from "./config/database";
import { connectRedis } from "./config/redis";
import { runMigrations } from "./db/migration_runner";

const PORT = Number(process.env.PORT || 3001);

async function bootStarp() {
  try {
    await runMigrations(db);
    await connectRedis();

    app.listen(PORT, () => {
      console.log(`Auth Service running on ${PORT}`);
    });
  } catch (error) {
    console.error("Bootstart failed", error);

    process.exit(1);
  }
}

bootStarp();
