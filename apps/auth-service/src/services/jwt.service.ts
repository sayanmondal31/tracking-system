import jwt from "jsonwebtoken";
import { env } from "../config/env";
import type { JwtPayload } from "../types/jwt.types";
import { AppError } from "../errors/AppError";
import type { IJWTSvc } from "../interface/IJwtSvc.interface";

export class JwtService implements IJWTSvc {
  private readonly accessTokenExpiry = "1m";
  private readonly refreshTokenExpiry = "30d";

  generateAccessToken = (userId: string, email: string): string => {
    return jwt.sign(
      {
        sub: userId,
        email,
      },
      env.JWT_SECRET,
      { expiresIn: this.accessTokenExpiry },
    );
  };

  generateRefreshToken = (userId: string, email: string): string => {
    return jwt.sign(
      {
        sub: userId,
        email,
      },
      env.JWT_SECRET,
      { expiresIn: this.refreshTokenExpiry },
    );
  };

  verifyToken = (token: string): JwtPayload => {
    try {
      return jwt.verify(token, env.JWT_SECRET) as JwtPayload; // service should return typed data
    } catch (error) {
      throw new AppError("Invalid or expired token", 401);
    }
  };
}
