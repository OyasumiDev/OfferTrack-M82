import { EmbeddingModel, FlagEmbedding } from "fastembed";

let model: FlagEmbedding | null = null;

async function getModel(): Promise<FlagEmbedding> {
  if (!model) {
    model = await FlagEmbedding.init({ model: EmbeddingModel.BGEBaseENV15 });
  }
  return model;
}

export async function embedTexts(texts: string[]): Promise<number[][]> {
  const m = await getModel();
  const results: number[][] = [];
  for await (const batch of m.embed(texts, 32)) {
    results.push(...batch);
  }
  return results;
}

export async function embedQuery(text: string): Promise<number[]> {
  const m = await getModel();
  return m.queryEmbed(text);
}
