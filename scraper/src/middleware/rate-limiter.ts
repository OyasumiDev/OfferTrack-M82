// scraper/src/middleware/rate-limiter.ts
import type { Request, Response, NextFunction } from 'express';

const requests = new Map<string, number>();

function rateLimiter(req: Request, res: Response, next: NextFunction): void {
  const key = req.ip || 'local';
  const now = Date.now();
  const last = requests.get(key) || 0;

  if (now - last < 1000) {
    res.status(429).json({ error: 'Too many requests — espera 1 segundo entre llamadas' });
    return;
  }

  requests.set(key, now);
  next();
}

module.exports = { rateLimiter };
