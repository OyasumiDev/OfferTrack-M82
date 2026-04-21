// scraper/src/embeddings/embed.service.ts
const { EmbeddingModel, FlagEmbedding } = require('fastembed');
import type { JobRaw } from '../portals/base.portal';

let model: any = null;

async function getModel(): Promise<any> {
  if (!model) {
    // BGESmallENV15 = BAAI/bge-small-en-v1.5 → 384 dims
    model = await FlagEmbedding.init({ model: EmbeddingModel.BGESmallENV15 });
  }
  return model;
}

// fastembed devuelve Float32Array por batch — convertir a number[] para JSON correcto
function toArray(v: unknown): number[] {
  if (Array.isArray(v)) return v as number[];
  if (v instanceof Float32Array || ArrayBuffer.isView(v)) return Array.from(v as Float32Array);
  return Object.values(v as Record<string, number>);
}

async function embedTexts(texts: string[]): Promise<number[][]> {
  const m = await getModel();
  const results: number[][] = [];
  for await (const batch of m.embed(texts, 32)) {
    for (const vec of batch) {
      results.push(toArray(vec));
    }
  }
  return results;
}

async function embedQuery(text: string): Promise<number[]> {
  const m = await getModel();
  const raw = await m.queryEmbed(text);
  return toArray(raw);
}

async function embedJob(job: JobRaw): Promise<number[]> {
  const text = `${job.title} ${job.company} ${job.description}`.slice(0, 512);
  const embeddings = await embedTexts([text]);
  return embeddings[0];
}

module.exports = { getModel, embedTexts, embedQuery, embedJob };
