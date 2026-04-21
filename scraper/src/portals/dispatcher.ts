// scraper/src/portals/dispatcher.ts
// Rutea el nombre de portal a su implementación — no contiene lógica de scraping
const { OccPortal } = require('./occ');
const { ComputrabajoPortal } = require('./computrabajo');
const { IndeedPortal } = require('./indeed');

const REGISTRY: Record<string, any> = {
  occ: new OccPortal(),
  computrabajo: new ComputrabajoPortal(),
  indeed: new IndeedPortal(),
};

// Devuelve la instancia del portal o null si no existe
function getPortal(name: string): any | null {
  return REGISTRY[name] ?? null;
}

// Lista los portales disponibles
function availablePortals(): string[] {
  return Object.keys(REGISTRY);
}

module.exports = { getPortal, availablePortals, REGISTRY };
