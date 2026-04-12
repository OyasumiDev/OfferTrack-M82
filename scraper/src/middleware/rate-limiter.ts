import { Request, Response, NextFunction } from "express";

const requests = new Map<string, number>();

export function rateLimiter(req: Request, res: Response, next: NextFunction) {
  const key = req.ip || "local";
  const now = Date.now();
  const last = requests.get(key) || 0;
  if (now - last < 500) {
    return res.status(429).json({ error: "Too many requests" });
  }
  requests.set(key, now);
  next();
}
