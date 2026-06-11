import { db } from "../config/database";
import type { RefreshToken } from "../types/refresh-token.types";

export class RefreshTokenRepository {
  create = async (
    id: string,
    userId: string,
    token: string,
    expiresAt: Date,
  ) => {
    const result = await db.query(
      `
          INSERT INTO refresh_tokens (
               id,
               user_id,
               token_hash,
               expires_at
          ) VALUES($1,$2,$3,$4)
           RETURNING *;
     `,
      [id, userId, token, expiresAt],
    );

    return result.rows[0];
  };

  findByToken = async (token: string): Promise<RefreshToken | null> => {
    const result = await db.query<RefreshToken>(
      `
          SELECT * FROM refresh_tokens 
          WHERE token_hash = $1;
          `,
      [token],
    );

    return result.rows[0] ?? null;
  };

  deleteByID = async (id: string): Promise<void> => {
    await db.query(
      `
          DELETE FROM refresh_tokens 
          WHERE id = $1;
          `,
      [id],
    );
  };
}
