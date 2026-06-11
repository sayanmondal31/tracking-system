import { JwtPayload } from "./jwt.types";

// Typescript understands what ```[req.user]``` is
// declaration merging
declare global {
  namespace Express {
    interface Request {
      user?: JwtPayload;
    }
  }
}

export {};
