import type { Request, Response, NextFunction } from "express";
import { JwtService } from "../services/jwt.service";
import { AppError } from "../errors/AppError";

const jwtService = new JwtService();

export const authenticate = (
  req: Request,
  _res: Response,
  next: NextFunction,
) => {
  const authHeaders = req.headers.authorization;

  if (!authHeaders) {
    throw new AppError("Unauthorized", 401);
  }

  const token = authHeaders.replace("Bearer ", "");

  const payload = jwtService.verifyToken(token);

  req.user = payload;

  next();
};
