// scraper/src/portals/computrabajo/index.ts
// ComputrabajoPortal — implementación aislada, sin imports de otros portales
const cheerio = require('cheerio') as typeof import('cheerio');
const { BasePortal } = require('../base.portal');
import type { Element } from 'domhandler';
import type { CheerioAPI } from 'cheerio';
import type { JobRaw, SearchParams } from '../base.portal';

// slugify propio — algoritmo idéntico al de OCC pero en su propio archivo
function slugify(s: string): string {
  return s
    .toLowerCase()
    .normalize('NFD').replace(/[\u0300-\u036f]/g, '')
    .replace(/\./g, ' ')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');
}

class ComputrabajoPortal extends BasePortal {
  name = 'computrabajo';
  private readonly baseUrl = 'https://mx.computrabajo.com';

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    const jobs: JobRaw[] = [];
    const seenUrls = new Set<string>();
    const maxPages = params.maxPages ?? 3;
    const browser = await this.createBrowser();

    try {
      const page = await this.newPage(browser);

      // URL canónica per Prompt Maestro §3.1: /trabajo-de-{puesto}-en-{ciudad}
      const roleSlug = slugify(params.role);
      const locSlug = params.location ? slugify(params.location.split(',')[0].trim()) : '';
      const baseUrl = locSlug
        ? `${this.baseUrl}/trabajo-de-${roleSlug}-en-${locSlug}`
        : `${this.baseUrl}/trabajo-de-${roleSlug}`;

      for (let p = 1; p <= maxPages; p++) {
        const pageUrl = p === 1 ? baseUrl : `${baseUrl}?p=${p}`;
        console.log(`[computrabajo] Navegando página ${p}: ${pageUrl}`);

        await page.goto(pageUrl, { waitUntil: 'domcontentloaded', timeout: 30000 });
        await this.delay(800, 1500);

        await page.waitForSelector('article.box_offer, [class*="offer_"]', { timeout: 10000 })
          .catch(() => {});

        const html = await page.content();
        const $ = cheerio.load(html);

        let newOnPage = 0;
        $('article.box_offer, [class*="offer_"]').each((_i: number, el: Element) => {
          const job = this.extractJobFromElement($, el);
          if (job && !seenUrls.has(job.url)) {
            seenUrls.add(job.url);
            jobs.push(job);
            newOnPage++;
          }
        });

        console.log(`[computrabajo] Página ${p}: ${newOnPage} vacantes nuevas`);
        if (newOnPage === 0) break;
        if (p < maxPages) await this.delay(800, 1500);
      }

      console.log(`[computrabajo] Total extraídas: ${jobs.length} vacantes`);
    } catch (err) {
      console.error('[computrabajo] Error durante scraping:', err);
    } finally {
      await browser.close();
    }

    return jobs;
  }

  private extractJobFromElement($: CheerioAPI, el: Element): JobRaw | null {
    try {
      const title =
        $('h2 a, [class*="title"] a', el).first().text().trim() ||
        $('h2', el).first().text().trim();

      const company =
        $('[class*="company"] a, [class*="empresa"]', el).first().text().trim() ||
        $('p.fs16', el).first().text().trim();

      const salary =
        $('[class*="salary"], [class*="salario"]', el).first().text().trim() || '';

      const modality =
        $('[class*="modality"], [class*="modalidad"], [class*="remote"]', el).first().text().trim() || '';

      const relativeUrl =
        $('h2 a', el).first().attr('href') ||
        $('a[href*="/empleo/"]', el).first().attr('href') || '';

      const rawLoc =
        $('[class*="location"], [class*="ciudad"]', el).first().text().trim() ||
        $('p[class*="fs13"]', el).first().text().trim();

      const isDate = /^(hace|más de|hoy|ayer|\d+\s+de\s+\w)/i.test(rawLoc);
      const hrefCity = (relativeUrl.match(/.*-en-([a-záéíóúüñ-]+)-[a-f0-9]{32}/i) || [])[1] || '';
      const location = isDate || !rawLoc
        ? (hrefCity ? hrefCity.replace(/-/g, ' ') : '')
        : rawLoc;

      const fullUrl = relativeUrl.startsWith('http')
        ? relativeUrl
        : `${this.baseUrl}${relativeUrl}`;

      if (!title || !relativeUrl) return null;

      return {
        title,
        company: company || 'Empresa confidencial',
        description: $('[class*="description"], [class*="descripcion"]', el).text().trim() || title,
        salary, modality, location,
        url: fullUrl,
        portal: 'computrabajo',
      };
    } catch {
      return null;
    }
  }
}

module.exports = { ComputrabajoPortal };
