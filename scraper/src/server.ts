// scraper/src/server.ts
import express = require('express');
import dotenv = require('dotenv');
import path = require('path');

dotenv.config({ path: path.resolve(__dirname, '../../.env') });

const { scrapeRoutes } = require('./routes/scrape.routes');
const { embedRoutes } = require('./routes/embed.routes');
const { healthRoutes } = require('./routes/health.routes');
const { errorHandler } = require('./middleware/error-handler');
const { rateLimiter } = require('./middleware/rate-limiter');

const app = express();
const PORT = process.env.SCRAPER_PORT || 3001;

app.use(express.json({ limit: '10mb' }));
app.use('/scrape', rateLimiter, scrapeRoutes);
app.use('/embed', embedRoutes);
app.use('/health', healthRoutes);
app.use(errorHandler);

app.listen(PORT, () => {
  console.log(`OfferTrack Scraper corriendo en puerto ${PORT}`);

  // Pre-calentar el modelo de embeddings en background para que el TEST 2 no espere
  const { getModel } = require('./embeddings/embed.service');
  getModel()
    .then(() => console.log('[embed] Modelo BAAI/bge-small-en-v1.5 listo'))
    .catch((e: Error) => console.warn('[embed] Aviso cargando modelo:', e.message));
});
