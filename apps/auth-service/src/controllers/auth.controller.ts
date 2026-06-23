import type { Request, Response } from "express";
import { AuthService } from "../services/auth.service";

import { JwtService } from "../services/jwt.service";
import type { User } from "../types/user.types";

export class AuthController {
  private authService = new AuthService();
  private jwtService = new JwtService();

  register = async (req: Request, res: Response) => {
    const { email, password } = req.body;

    // const {user,otp} = await this.authService.register(email, password) as {user:User,otp:string};

    const { newUser: user, otp } = (await this.authService.register(
      email,
      password,
    )) as { newUser: User; otp: string };

    return res.status(201).json({
      success: true,
      data: {
        id: user?.id,
        email: user?.email,
        otp: otp,
      },
    });
  };

  login = async (req: Request, res: Response) => {
    const { email, password } = req.body;

    const { accessToken, refreshToken } = await this.authService.login(
      email,
      password,
    );

    return res.status(200).json({
      success: true,
      data: {
        accessToken,
        refreshToken,
      },
    });
  };

  me = async (req: Request, res: Response) => {
    return res.status(200).json({
      success: true,
      data: req.user,
    });
  };

  refresh = async (req: Request, res: Response) => {
    const { refreshToken } = req.body;

    const tokens = await this.authService.refreshToken(refreshToken);

    return res.status(200).json({
      success: true,
      data: tokens,
    });
  };

  verify = async (req: Request, res: Response) => {
    const { email, otp } = req.body;

    const verifyUser = await this.authService.verifyOtp(otp, email);

    return res.status(200).json({
      success: true,
      data: verifyUser,
    });
  };

  logout = async (req: Request, res: Response) => {
    const { refreshToken } = req.body;

    // Extract access token from the Authorization header
    const authHeader: string = req.headers.authorization || "";
    const accessToken = authHeader.replace("Bearer ", "");

    this.authService.logout(refreshToken, accessToken);

    return res.status(200).json({
      success: true,
    });
  };
}
