// scraper/src/portals/occ/url.ts
// Construcción de URL canónica OCC — sin dependencias externas, sin Playwright
import type { SearchParams } from '../base.portal';

// Slugify per Prompt Maestro §2.3:
// 1. lowercase  2. quitar tildes  3. "." → espacio  4. [^a-z0-9]+→"-"  5. colapsar  6. trim
function slugify(s: string): string {
  return s
    .toLowerCase()
    .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
    .replace(/\./g, ' ')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');
}

// Ciudad → estado (slugs ya normalizados)
const CITY_TO_STATE: Record<string, string> = {
  'monterrey': 'nuevo-leon', 'san-nicolas-de-los-garza': 'nuevo-leon',
  'guadalupe': 'nuevo-leon', 'apodaca': 'nuevo-leon', 'san-pedro-garza-garcia': 'nuevo-leon',
  'guadalajara': 'jalisco', 'zapopan': 'jalisco', 'tlaquepaque': 'jalisco', 'tonala': 'jalisco',
  'mexico': 'ciudad-de-mexico', 'cdmx': 'ciudad-de-mexico', 'ciudad-de-mexico': 'ciudad-de-mexico',
  'puebla': 'puebla',
  'tijuana': 'baja-california', 'mexicali': 'baja-california', 'ensenada': 'baja-california',
  'leon': 'guanajuato', 'irapuato': 'guanajuato', 'celaya': 'guanajuato',
  'queretaro': 'queretaro', 'san-juan-del-rio': 'queretaro',
  'merida': 'yucatan',
  'cancun': 'quintana-roo', 'playa-del-carmen': 'quintana-roo', 'chetumal': 'quintana-roo',
  'hermosillo': 'sonora', 'ciudad-obregon': 'sonora',
  'chihuahua': 'chihuahua', 'ciudad-juarez': 'chihuahua',
  'culiacan': 'sinaloa', 'mazatlan': 'sinaloa',
  'torreon': 'coahuila', 'saltillo': 'coahuila', 'monclova': 'coahuila',
  'san-luis-potosi': 'san-luis-potosi',
  'aguascalientes': 'aguascalientes',
  'morelia': 'michoacan', 'uruapan': 'michoacan',
  'toluca': 'estado-de-mexico', 'ecatepec': 'estado-de-mexico', 'naucalpan': 'estado-de-mexico',
  'veracruz': 'veracruz', 'xalapa': 'veracruz',
  'acapulco': 'guerrero',
  'oaxaca': 'oaxaca',
  'villahermosa': 'tabasco',
  'tuxtla-gutierrez': 'chiapas',
  'durango': 'durango',
  'zacatecas': 'zacatecas',
  'tepic': 'nayarit',
  'colima': 'colima',
  'campeche': 'campeche',
  'la-paz': 'baja-california-sur', 'los-cabos': 'baja-california-sur',
  'pachuca': 'hidalgo',
  'cuernavaca': 'morelos',
  'tlaxcala': 'tlaxcala',
};

const OCC_BASE = 'https://www.occ.com.mx';

// Construye URL canónica de OCC per Prompt Maestro §2.2
// location acepta: "Ciudad, Estado" | "Ciudad" (autodetecta estado) | solo "Estado"
function buildOccUrl(params: SearchParams): string {
  const roleSlug = slugify(params.role);
  let url = `${OCC_BASE}/empleos/de-${roleSlug}/`;

  if (params.location) {
    const parts = params.location.split(',').map((s: string) => s.trim());
    if (parts.length >= 2) {
      // "Monterrey, Nuevo León" → estado + ciudad explícitos
      url += `en-${slugify(parts[1])}/en-la-ciudad-de-${slugify(parts[0])}/`;
    } else {
      const locSlug = slugify(parts[0]);
      const state = CITY_TO_STATE[locSlug];
      url += state
        ? `en-${state}/en-la-ciudad-de-${locSlug}/`
        : `en-${locSlug}/`;
    }
  }

  const mod = (params.modality || '').toLowerCase();
  if (mod.includes('remoto') || mod.includes('remote') || mod.includes('home office')) {
    url += 'tipo-home-office-remoto';
  }

  const qs: string[] = [];
  if (params.salaryMin && params.salaryMax) qs.push(`salary=${params.salaryMin},${params.salaryMax}`);
  else if (params.salaryMin) qs.push(`salary=${params.salaryMin},`);
  if (qs.length) url += `?${qs.join('&')}`;

  return url;
}

module.exports = { slugify, buildOccUrl, CITY_TO_STATE };
