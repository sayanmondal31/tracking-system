import { redis } from "../config/redis";
import { sha256 } from "../utils/hash";

export class OtpRepository {
  private readonly otpKey = "emailotp";
  private readonly otpExpiry = 300;

  setOtp = async (otp: string, email: string) => {
    const hashOtp = sha256(otp);

    await redis.set(`${this.otpKey}:${email}`, hashOtp, { EX: this.otpExpiry });
  };

  getOtp = async (email: string) => {
    const otp = await redis.get(`${this.otpKey}:${email}`);

    return otp;
  };
}
