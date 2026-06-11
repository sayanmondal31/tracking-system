import { Pool } from "pg";
import type { QueryResult } from "pg";
import fs from "fs/promises";
import path from "path";

export async function runMigrations(pool: Pool) {
  // Create migration tracking table
  await pool.query(`
    CREATE TABLE IF NOT EXISTS migrations (
      id SERIAL,
      name VARCHAR(255) PRIMARY KEY,
      executed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    );
  `);

  const migrationsDir = path.join(__dirname, "migrations");

  const files = await fs.readdir(migrationsDir);

  const upFiles = files.filter((file) => file.endsWith(".up.sql")).sort();

  for (const file of upFiles) {
    const migrateName = file.replace(".up.sql", "");

    const existingMigration = await pool.query(
      `
        SELECT id FROM migrations WHERE name = $1;
      `,
      [migrateName],
    );

    if ((existingMigration.rowCount ?? 0) > 0) {
      continue;
    }

    console.log(`Running migration: ${migrateName}`);

    const sql = await fs.readFile(path.join(migrationsDir, file), "utf-8");

    await pool.query("BEGIN");

    try {
      await pool.query(sql);

      await pool.query(`INSERT INTO migrations(name) VALUES($1)`, [
        migrateName,
      ]);

      await pool.query("COMMIT");
    } catch (error) {
      await pool.query("ROLLBACK");
      throw error;
    }
  }
}
