import type { LoginResponse } from "../types/auth.types";
import type { VerifyOtpResponse } from "../types/otp.types";
import type { User } from "../types/user.types";

export interface IAuthSvc {
  register(
    email: string,
    password: string,
  ): Promise<{ newUser: User | null; otp: string }>;
  login(email: string, password: string): Promise<LoginResponse>;
  refreshToken(token: string): Promise<{
    accessToken: string;
    refreshToken: string;
  }>;
  verifyOtp(otp: string, email: string): Promise<VerifyOtpResponse>;
  logout(token: string, accessToken: String): Promise<void>;
}
