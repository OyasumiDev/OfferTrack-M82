import { Router } from "express";
import { embedTexts, embedQuery } from "../embeddings/embed.service";

export const embedRoutes = Router();

embedRoutes.post("/texts", async (req, res) => {
  const { texts } = req.body;
  const embeddings = await embedTexts(texts);
  res.json({ embeddings });
});

embedRoutes.post("/query", async (req, res) => {
  const { text } = req.body;
  const embedding = await embedQuery(text);
  res.json({ embedding });
});
