// scraper/src/portals/occ/parse.ts
// Extracción y parseo de job cards de OCC — sin Playwright, sin state

const NOISE = new Set([
  'hoy', 'ayer', 'recomendada', 'vista recientemente.', 'vista recientemente',
  'nueva', 'nuevo',
  'ya estás postulado.', 'ya estás postulado', 'postulado', 'ya aplicaste',
  'ver oferta', 'guardar', 'aplica ya', 'compartir',
]);

const BENEFIT_PREFIXES = ['prestaciones', 'vales', 'fondo', 'seguro', 'plan de', 'excelente'];

function isBenefit(text: string): boolean {
  const lower = text.toLowerCase();
  return BENEFIT_PREFIXES.some(p => lower.startsWith(p));
}

function isNoise(text: string): boolean {
  return NOISE.has(text.toLowerCase()) || /^hace \d+/.test(text.toLowerCase());
}

function isSalary(text: string): boolean {
  return /\$\s*\d/.test(text);
}

// Descarta vacantes completamente irrelevantes al rol buscado (keyword mínimo 3 chars)
function isTitleRelevant(title: string, role: string): boolean {
  if (!title) return false;
  const words = role.toLowerCase().replace(/[^a-záéíóúüñ\s]/gi, ' ').split(/\s+/).filter(w => w.length >= 3);
  if (words.length === 0) return true;
  const t = title.toLowerCase();
  return words.some(w => t.includes(w));
}

// Extrae campos semánticos del array de texto plano de un job card de OCC
function parseCardText(lines: string[]): {
  title: string; salary: string; company: string; location: string; benefits: string[];
} {
  const clean = lines.map(l => l.trim()).filter(Boolean);
  let title = '';
  let salary = '';
  const benefits: string[] = [];
  const rest: string[] = [];

  let titleFound = false;
  for (const line of clean) {
    if (isNoise(line)) continue;
    if (!titleFound) { title = line; titleFound = true; continue; }
    if (isSalary(line)) { salary = line; continue; }
    if (isBenefit(line)) { benefits.push(line); continue; }
    rest.push(line);
  }

  // Último elemento de rest suele ser "Ciudad, Estado"
  const location = rest.length > 0 && rest[rest.length - 1].includes(',')
    ? rest[rest.length - 1]
    : rest[rest.length - 1] || '';
  const company = rest.length > 1
    ? rest[rest.length - 2]
    : (rest.length === 1 && !location ? rest[0] : '');

  return { title, salary, company, location, benefits };
}

// Limpia el HTML de la descripción que devuelve la API oferta.occ.com.mx
function stripHtml(html: string): string {
  return html
    .replace(/<[^>]+>/g, ' ')
    .replace(/&amp;/g, '&').replace(/&lt;/g, '<').replace(/&gt;/g, '>')
    .replace(/&nbsp;/g, ' ').replace(/&[a-z]+;/g, '')
    .replace(/\s+/g, ' ')
    .trim();
}

module.exports = { isTitleRelevant, parseCardText, stripHtml };
