import bcrypt from "bcryptjs";
import { v4 as uuid } from "uuid";
import { UserRepository } from "../repositories/user.repository";
import { AppError } from "../errors/AppError";
import { JwtService } from "./jwt.service";
import type { User } from "../types/user.types";
import type { LoginResponse } from "../types/auth.types";
import { RefreshTokenRepository } from "../repositories/refresh-token.repository";
import { sha256 } from "../utils/hash";
import { OtpRepository } from "../repositories/otp.repository";
import { generateOTP } from "../utils/generateOtp";
import type { VerifyOtpResponse } from "../types/otp.types";
import type { IAuthSvc } from "../interface/IAuthServer.interface";

export class AuthService implements IAuthSvc {
  private userRepository = new UserRepository();
  private refreshTokenRepository = new RefreshTokenRepository();
  private otpRepository = new OtpRepository();

  private jwtService = new JwtService();

  private readonly refreshTokenTTl = 30 * 24 * 60 * 60 * 1000;

  private createRefreshTokenSession = async (
    userId: string,
    email: string,
  ): Promise<string> => {
    const refrestoken = this.jwtService.generateRefreshToken(userId, email);

    const tokenHash = sha256(refrestoken);

    // save refresh token
    await this.refreshTokenRepository.create(
      uuid(),
      userId,
      tokenHash,
      new Date(Date.now() + this.refreshTokenTTl),
    );

    return refrestoken;
  };

  register = async (
    email: string,
    password: string,
  ): Promise<{ newUser: User | null; otp: string }> => {
    const existingUSer = await this.userRepository.findByEmail(email);

    if (existingUSer) {
      throw new AppError("User already exists", 409);
    }

    const passwordHash = await bcrypt.hash(password, 10);

    const newUser = await this.userRepository.create(
      uuid(),
      email,
      passwordHash,
    );

    const otp = generateOTP();

    await this.otpRepository.setOtp(otp, email);

    return { newUser, otp };
  };

  login = async (email: string, password: string): Promise<LoginResponse> => {
    const user = await this.userRepository.findByEmail(email);

    if (!user) {
      throw new AppError("Invalid credentials", 401);
    }

    const passwordMatch = await bcrypt.compare(password, user.password_hash);

    if (!passwordMatch) {
      throw new AppError("Invalid credential", 401);
    }

    const accessToken = this.jwtService.generateAccessToken(user.id, email);

    const refreshToken = await this.createRefreshTokenSession(user.id, email);

    return {
      accessToken,
      refreshToken,
    };
  };

  refreshToken = async (token: string) => {
    const payload = this.jwtService.verifyToken(token);

    const tokenHash = sha256(token);

    const existingToken =
      await this.refreshTokenRepository.findByToken(tokenHash);

    if (!existingToken) {
      throw new AppError("Invalid refresh token", 401);
    }

    await this.refreshTokenRepository.deleteByID(existingToken.id);

    const accessToken = this.jwtService.generateAccessToken(
      payload.sub,
      payload.email,
    );

    const refreshToken = await this.createRefreshTokenSession(
      payload.sub,
      payload.email,
    );

    return {
      accessToken,
      refreshToken,
    };
  };

  verifyOtp = async (
    otp: string,
    email: string,
  ): Promise<VerifyOtpResponse> => {
    const otpHash = sha256(otp);

    const existingOtpHash = await this.otpRepository.getOtp(email);

    if (!existingOtpHash) {
      throw new AppError("Otp expired!", 401);
    }

    if (otpHash != existingOtpHash) {
      throw new AppError("Otp invalid!", 400);
    }

    const user = await this.userRepository.update(email);

    return {
      id: user.id,
      email: user.email,
      is_verified: user.is_verified,
    } as VerifyOtpResponse;
  };

  logout = async (token: string) => {
    const tokenHaash = sha256(token);

    const existingToken =
      await this.refreshTokenRepository.findByToken(tokenHaash);

    if (!existingToken) {
      return;
    }

    await this.refreshTokenRepository.deleteByID(existingToken.id);
  };
}
