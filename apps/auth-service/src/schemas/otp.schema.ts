import z from "zod";

export const OtpSchema = z.object({
  email: z.email(),
  otp: z.string().length(6),
});

export type OtpInput = z.infer<typeof OtpSchema>;
