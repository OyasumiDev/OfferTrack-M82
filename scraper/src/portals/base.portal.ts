// scraper/src/portals/base.portal.ts
const { chromium } = require('playwright');

export interface JobRaw {
  title: string;
  company: string;
  description: string;
  salary?: string;
  modality?: string;
  location?: string;
  url: string;
  portal: string;
  externalId?: string; // ID numérico del portal (ej: data-id de OCC) para deduplicación en Qdrant
}

export interface SearchParams {
  role: string;
  location?: string;
  salaryMin?: number;
  salaryMax?: number;
  modality?: string;
  maxPages?: number; // Páginas a scraper por portal (default 3)
}

// Inyectado en cada nueva página para eliminar huellas de automatización
const STEALTH_SCRIPT = `
  // Eliminar navigator.webdriver (principal señal de bot)
  Object.defineProperty(navigator, 'webdriver', { get: () => undefined });

  // Simular chrome runtime (ausente en headless)
  window.chrome = {
    runtime: {},
    loadTimes: () => {},
    csi: () => {},
    app: {},
  };

  // Corregir navigator.plugins (vacío en headless)
  Object.defineProperty(navigator, 'plugins', {
    get: () => [{ name: 'Chrome PDF Plugin' }, { name: 'Chrome PDF Viewer' }],
  });

  // Corregir navigator.languages
  Object.defineProperty(navigator, 'languages', {
    get: () => ['es-MX', 'es', 'en-US', 'en'],
  });

  // Corregir Notification.permission (bloqueado en headless)
  if (typeof Notification !== 'undefined') {
    Object.defineProperty(Notification, 'permission', { get: () => 'default' });
  }
`;

abstract class BasePortal {
  abstract name: string;
  abstract scrape(params: SearchParams): Promise<JobRaw[]>;

  protected async createBrowser() {
    return chromium.launch({
      // false = Chrome visible — necesario para pasar Cloudflare en OCC
      // En un servidor sin pantalla usar xvfb-run o cambiar a 'new' y evaluar
      headless: false,
      args: [
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--disable-blink-features=AutomationControlled',
        '--disable-features=IsolateOrigins,site-per-process',
        '--disable-dev-shm-usage',
      ],
    });
  }

  protected async newPage(browser: any) {
    const { getRandomUserAgent } = require('../utils/user-agent');
    const context = await browser.newContext({
      userAgent: getRandomUserAgent(),
      locale: 'es-MX',
      timezoneId: 'America/Mexico_City',
      viewport: { width: 1366, height: 768 },
      extraHTTPHeaders: {
        'Accept-Language': 'es-MX,es;q=0.9,en-US;q=0.8,en;q=0.7',
        'Accept-Encoding': 'gzip, deflate, br',
        'Sec-Fetch-Dest': 'document',
        'Sec-Fetch-Mode': 'navigate',
        'Sec-Fetch-Site': 'none',
        'Upgrade-Insecure-Requests': '1',
      },
    });
    const page = await context.newPage();
    await page.addInitScript(STEALTH_SCRIPT);
    return page;
  }

  protected delay(min?: number, max?: number): Promise<void> {
    const { delay } = require('../utils/delay');
    return delay(min, max);
  }
}

module.exports = { BasePortal };
