import { db } from "../config/database";
import type { User } from "../types/user.types";

export class UserRepository {
  findByEmail = async (email: string): Promise<User | null> => {
    const result = await db.query<User>(
      `
      SELECT * 
      FROM users
      WHERE email=$1;
          `,
      [email],
    );

    return result.rows[0] ?? null;
  };

  create = async (
    id: string,
    email: string,
    password_hash: string,
  ): Promise<User | null> => {
    const result = await db.query<User>(
      `
        INSERT INTO users (
          id,
          email,
          password_hash
        ) 
        VALUES($1, $2, $3)
        RETURNING *;
        `,
      [id, email, password_hash],
    );

    return result.rows[0];
  };

  update = async (email: string): Promise<User> => {
    const result = await db.query<User>(
      `
      UPDATE users
      SET is_email_verified = TRUE
      WHERE email=$1
      RETURNING *
      `,
      [email],
    );

    return result.rows[0];
  };
}
