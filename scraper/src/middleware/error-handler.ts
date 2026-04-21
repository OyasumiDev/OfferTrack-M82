// scraper/src/middleware/error-handler.ts
import type { Request, Response, NextFunction } from 'express';

function errorHandler(
  err: Error & { status?: number },
  _req: Request,
  res: Response,
  _next: NextFunction
): void {
  console.error('[error-handler]', err.stack || err.message);
  res.status(err.status || 500).json({
    error: err.message || 'Internal server error',
    status: err.status || 500,
  });
}

module.exports = { errorHandler };
