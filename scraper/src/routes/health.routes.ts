// scraper/src/routes/health.routes.ts
import express = require('express');
import type { Request, Response } from 'express';

const healthRoutes = express.Router();

healthRoutes.get('/', (_req: Request, res: Response) => {
  res.json({
    status: 'ok',
    service: 'offertrack-scraper',
    timestamp: new Date().toISOString(),
  });
});

module.exports = { healthRoutes };
