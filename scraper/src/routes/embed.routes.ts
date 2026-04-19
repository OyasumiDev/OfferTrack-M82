// scraper/src/routes/embed.routes.ts
import express = require('express');
import type { Request, Response, NextFunction } from 'express';

const embedRoutes = express.Router();
const { embedTexts, embedQuery } = require('../embeddings/embed.service');
const { EMBED_MODEL, EMBED_DIMENSIONS } = require('../embeddings/models');

embedRoutes.post('/texts', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const { texts } = req.body as { texts: string[] };
    if (!Array.isArray(texts) || texts.length === 0) {
      res.status(400).json({ error: 'Se requiere un array "texts" no vacío' });
      return;
    }
    const embeddings: number[][] = await embedTexts(texts);
    res.json({ embeddings, dimensions: EMBED_DIMENSIONS, model: EMBED_MODEL });
  } catch (err) {
    next(err);
  }
});

embedRoutes.post('/query', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const { text } = req.body as { text: string };
    if (!text) {
      res.status(400).json({ error: 'Se requiere el campo "text"' });
      return;
    }
    const embedding: number[] = await embedQuery(text);
    res.json({ embedding, dimensions: EMBED_DIMENSIONS, model: EMBED_MODEL });
  } catch (err) {
    next(err);
  }
});

module.exports = { embedRoutes };
