const cheerio = require('cheerio') as typeof import('cheerio');
const { BasePortal } = require('../base.portal');
import type { JobRaw, SearchParams } from '../base.portal';
import type { Element } from 'domhandler';
import type { CheerioAPI } from 'cheerio';
import { readEmbeddedResults } from './embedded';
import { BuildIndeedURL } from './url';
import { INDEED_SELECTORS } from './selectors';

const INDEED_BASE = 'https://mx.indeed.com';

class IndeedPortal extends BasePortal {
  name = 'indeed';

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    const maxPages = params.maxPages ?? 3;
    const seen = new Set<string>();
    const jobs: JobRaw[] = [];
    const browser = await this.createBrowser();

    try {
      const page = await this.newPage(browser);

      for (let p = 1; p <= maxPages; p++) {
        const url = BuildIndeedURL({
          puesto: params.role,
          ciudad: params.location,
          page: p,
        });

        console.log(`[indeed] Página ${p}: ${url}`);
        await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });
        await this.delay(1500, 2500);

        // Intento 1: JSON embebido en window.mosaic
        const embedded = await readEmbeddedResults(page);
        if (embedded && embedded.length > 0) {
          console.log(`[indeed] Página ${p}: ${embedded.length} resultados vía JSON embebido`);
          for (const r of embedded) {
            if (r.expired || !r.jobkey || seen.has(r.jobkey)) continue;
            seen.add(r.jobkey);
            jobs.push({
              title: r.displayTitle || '',
              company: r.company || 'Empresa confidencial',
              description: r.snippet || r.displayTitle || '',
              salary: r.salarySnippet?.text || '',
              modality: r.remoteLocation ? 'remoto' : '',
              location: r.formattedLocation || '',
              url: `${INDEED_BASE}/viewjob?jk=${r.jobkey}`,
              portal: 'indeed',
              externalId: r.jobkey,
            });
          }
          if (embedded.length < 15) break;
          if (p < maxPages) await this.delay(3000, 5000);
          continue;
        }

        // Intento 2: fallback DOM con Cheerio
        console.log(`[indeed] Página ${p}: JSON embebido vacío — usando extracción DOM`);
        const html = await page.content();
        const $ = cheerio.load(html);

        const cards = $(INDEED_SELECTORS.cardContainer);
        console.log(`[indeed] Página ${p}: ${cards.length} tarjetas DOM encontradas`);

        if (cards.length === 0) {
          const raw = html.substring(0, 500);
          console.log(`[indeed] WARN: página sin resultados. Preview: ${raw}`);
          break;
        }

        let newOnPage = 0;
        cards.each((_i: number, el: Element) => {
          const job = this.extractFromCard($, el);
          if (!job || !job.externalId || seen.has(job.externalId)) return;
          seen.add(job.externalId!);
          jobs.push(job);
          newOnPage++;
        });

        console.log(`[indeed] Página ${p}: ${newOnPage} nuevas vacantes (DOM)`);
        if (cards.length < 15 || newOnPage === 0) break;
        if (p < maxPages) await this.delay(3000, 5000);
      }
    } catch (err) {
      console.error('[indeed] Error durante scraping:', err);
    } finally {
      await browser.close();
    }

    console.log(`[indeed] Total: ${jobs.length} vacantes`);
    return jobs;
  }

  private extractFromCard($: CheerioAPI, el: Element): JobRaw | null {
    try {
      const jk = (el as any).attribs?.['data-jk'] ||
        $('a[data-jk]', el).first().attr('data-jk') || '';

      const title =
        $(INDEED_SELECTORS.title, el).first().attr('title') ||
        $(INDEED_SELECTORS.titleFallback, el).first().text().trim() ||
        $('h2.jobTitle', el).first().text().trim();

      if (!title || !jk) return null;

      const company = $(INDEED_SELECTORS.company, el).first().text().trim() || 'Empresa confidencial';
      const location = $(INDEED_SELECTORS.location, el).first().text().trim() || '';
      const salary = $(INDEED_SELECTORS.salary, el).first().text().trim() || '';

      return {
        title,
        company,
        description: title,
        salary,
        modality: '',
        location,
        url: `${INDEED_BASE}/viewjob?jk=${jk}`,
        portal: 'indeed',
        externalId: jk,
      };
    } catch {
      return null;
    }
  }
}

module.exports = { IndeedPortal };
