export function normalizeText(text: string): string {
  return text.trim().replace(/\s+/g, " ").toLowerCase();
}

export function normalizeJob(raw: Record<string, unknown>): Record<string, unknown> {
  return {
    ...raw,
    title: normalizeText(String(raw.title || "")),
    company: normalizeText(String(raw.company || "")),
  };
}
