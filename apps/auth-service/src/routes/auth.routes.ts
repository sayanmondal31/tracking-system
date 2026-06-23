import { Router } from "express";
import { validate } from "../middleware/validate";
import {
  loginSchema,
  logoutSchema,
  refreshTokenSchema,
  registerSchema,
} from "../schemas/auth.schema";
import { AuthController } from "../controllers/auth.controller";
import { authenticate } from "../middleware/auth.middleware";
import { OtpSchema } from "../schemas/otp.schema";

const router = Router();

const authController = new AuthController();

router
  .route("/register")
  .post(validate(registerSchema), authController.register);

router.route("/login").post(validate(loginSchema), authController.login);
router.route("/me").get(authenticate, authController.me);
router
  .route("/refresh")
  .post(validate(refreshTokenSchema), authController.refresh);

router.route("/verify").post(validate(OtpSchema), authController.verify);
router
  .route("/logout")
  .post(authenticate, validate(logoutSchema), authController.logout);

export default router;
