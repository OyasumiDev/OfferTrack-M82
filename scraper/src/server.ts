import express from "express";
import dotenv from "dotenv";
import { scrapeRoutes } from "./routes/scrape.routes";
import { embedRoutes } from "./routes/embed.routes";
import { healthRoutes } from "./routes/health.routes";
import { errorHandler } from "./middleware/error-handler";

dotenv.config();

const app = express();
const PORT = process.env.SCRAPER_PORT || 3001;

app.use(express.json({ limit: "10mb" }));
app.use("/scrape", scrapeRoutes);
app.use("/embed", embedRoutes);
app.use("/health", healthRoutes);
app.use(errorHandler);

app.listen(PORT, () => {
  console.log(`OfferTrack Scraper running on http://localhost:${PORT}`);
});
