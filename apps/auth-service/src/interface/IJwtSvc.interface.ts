import type { JwtPayload } from "../types/jwt.types";

export interface IJWTSvc {
  generateAccessToken(userId: string, email: string): string;
  generateRefreshToken(userId: string, email: string): string;
  verifyToken(token: string): JwtPayload;
}
