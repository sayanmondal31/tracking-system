import express from "express";

// Routes
import authRoutes from "./routes/auth.routes";
import { errorMiddleware } from "./middleware/error.middleware";

const app = express();

app.use(express.json());

app.get("/health", (_, res) => {
  return res.status(200).json({
    service: "auth-service",
    status: "healthy",
  });
});

app.use("/auth", authRoutes);

app.use(errorMiddleware);

export default app;
