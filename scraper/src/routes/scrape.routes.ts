// scraper/src/routes/scrape.routes.ts
import express = require('express');
import type { Request, Response, NextFunction } from 'express';
import type { SearchParams } from '../portals/base.portal';

const scrapeRoutes = express.Router();
const { getPortal } = require('../portals/dispatcher');
const { normalizeJob } = require('../utils/normalizer');
const { embedJob } = require('../embeddings/embed.service');

// Filtro de relevancia centralizado — aplica a todos los portales antes del embed
function isTitleRelevant(title: string, role: string): boolean {
  if (!title) return false;
  const words = role.toLowerCase().replace(/[^a-záéíóúüñ\s]/gi, ' ').split(/\s+/).filter((w: string) => w.length >= 3);
  if (words.length === 0) return true;
  const t = title.toLowerCase();
  return words.some((w: string) => t.includes(w));
}

scrapeRoutes.post('/', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const {
      role, location, salaryMin, salaryMax, modality,
      portals: requestedPortals,
    } = req.body as {
      role: string;
      location?: string;
      salaryMin?: number;
      salaryMax?: number;
      modality?: string;
      portals?: string[];
    };

    if (!role || role.trim() === '') {
      res.status(400).json({ error: 'El campo "role" es requerido' });
      return;
    }

    const selectedPortals: string[] = (requestedPortals && requestedPortals.length > 0)
      ? requestedPortals
      : ['occ'];

    const params: SearchParams = { role: role.trim(), location, salaryMin, salaryMax, modality };
    const allJobs: any[] = [];

    for (const portalName of selectedPortals) {
      const portal = getPortal(portalName);
      if (!portal) {
        console.warn(`[scrape] Portal desconocido: "${portalName}" — valores válidos: occ, computrabajo, indeed`);
        continue;
      }

      console.log(`[scrape] Scrapeando ${portalName} — role: "${role}", location: "${location || 'cualquiera'}"`);
      const rawJobs = await portal.scrape(params);

      for (const raw of rawJobs) {
        if (!isTitleRelevant(raw.title, role)) {
          console.log(`[scrape] Descartado (irrelevante): "${raw.title}"`);
          continue;
        }
        const normalized = normalizeJob(raw);
        const embedding = await embedJob(raw);
        allJobs.push({ ...normalized, embedding });
      }
    }

    res.json({ jobs: allJobs, count: allJobs.length });
  } catch (err) {
    next(err);
  }
});

module.exports = { scrapeRoutes };
