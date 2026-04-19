// scraper/src/embeddings/models.ts
const EMBED_MODEL: string = process.env.EMBED_MODEL || 'BAAI/bge-small-en-v1.5';
const EMBED_DIMENSIONS: number = parseInt(process.env.EMBED_DIMENSIONS || '384');

module.exports = { EMBED_MODEL, EMBED_DIMENSIONS };
