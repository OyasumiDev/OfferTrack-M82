// scraper/src/utils/normalizer.ts
import type { JobRaw } from '../portals/base.portal';

function cleanText(text: string): string {
  return text.trim().replace(/\s+/g, ' ');
}

function inferModality(text: string): string {
  const lower = text.toLowerCase();
  if (lower.includes('remoto') || lower.includes('home office') || lower.includes('trabajo desde casa') || lower.includes('teletrabajo')) {
    return 'remote';
  }
  if (lower.includes('híbrido') || lower.includes('hibrido') || lower.includes('hybrid')) {
    return 'hybrid';
  }
  return 'onsite';
}

function parseSalary(salaryText: string): { min: number; max: number; currency: string } {
  const result = { min: 0, max: 0, currency: 'MXN' };
  if (!salaryText) return result;

  const clean = salaryText.replace(/,/g, '').replace(/\./g, '');
  const numbers = clean.match(/\d+/g);
  if (!numbers || numbers.length === 0) return result;

  if (numbers.length >= 2) {
    result.min = parseInt(numbers[0]);
    result.max = parseInt(numbers[1]);
  } else {
    result.min = parseInt(numbers[0]);
    result.max = parseInt(numbers[0]);
  }

  if (salaryText.toLowerCase().includes('usd') || salaryText.includes('$') && salaryText.toLowerCase().includes('dolar')) {
    result.currency = 'USD';
  }

  return result;
}

// ID determinístico: usa externalId numérico del portal si existe, si no SHA-256 de la URL canónica.
// URL canónica = sin fragmento (#...) ni parámetros de posición (lc=Score-N) → mismo job = mismo ID.
function canonicalUrl(url: string): string {
  return url.trim().split('#')[0].split('?lc=')[0];
}

function stableId(raw: JobRaw): string {
  if (raw.externalId && /^\d+$/.test(raw.externalId)) {
    const n = raw.externalId.padStart(32, '0');
    return `${n.slice(0,8)}-${n.slice(8,12)}-4${n.slice(13,16)}-${n.slice(16,20)}-${n.slice(20,32)}`;
  }
  const hash = require('crypto').createHash('sha256').update(canonicalUrl(raw.url)).digest('hex');
  return `${hash.slice(0,8)}-${hash.slice(8,12)}-4${hash.slice(13,16)}-${hash.slice(16,20)}-${hash.slice(20,32)}`;
}

function normalizeJob(raw: JobRaw): JobRaw & {
  id: string;
  salaryMin: number;
  salaryMax: number;
  currency: string;
  modalityNormalized: string;
  scrapedAt: string;
} {
  const salaryData = parseSalary(raw.salary || '');
  const modalityNormalized = raw.modality
    ? inferModality(raw.modality)
    : inferModality(raw.description);

  return {
    ...raw,
    id: stableId(raw),
    title: cleanText(raw.title),
    company: cleanText(raw.company),
    description: cleanText(raw.description),
    location: raw.location ? cleanText(raw.location) : '',
    url: canonicalUrl(raw.url),
    salaryMin: salaryData.min,
    salaryMax: salaryData.max,
    currency: salaryData.currency,
    modalityNormalized,
    scrapedAt: new Date().toISOString(),
  };
}

module.exports = { normalizeJob, inferModality, parseSalary, cleanText };
