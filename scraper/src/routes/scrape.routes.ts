import { Router } from "express";

export const scrapeRoutes = Router();

scrapeRoutes.post("/", async (req, res) => {
  // TODO: implementar scraping
  res.json({ jobs: [] });
});
