// scraper/src/portals/occ/index.ts
// OccPortal — implementación completa del scraper de OCC Mundial
const { BasePortal } = require('../base.portal');
const { slugify, buildOccUrl } = require('./url');
const { isTitleRelevant, parseCardText, stripHtml } = require('./parse');
import type { JobRaw, SearchParams } from '../base.portal';

const OCC_BASE = 'https://www.occ.com.mx';
const DETAIL_BASE = 'https://oferta.occ.com.mx/offer';

class OccPortal extends BasePortal {
  name = 'occ';

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    const jobs: JobRaw[] = [];
    const seenIds = new Set<string>();
    const maxPages = params.maxPages ?? 3;
    const browser = await this.createBrowser();

    try {
      const page = await this.newPage(browser);
      const roleSlug = slugify(params.role);
      const baseUrl = buildOccUrl(params);

      for (let p = 1; p <= maxPages; p++) {
        const pageUrl = p === 1
          ? baseUrl
          : baseUrl.includes('?') ? `${baseUrl}&page=${p}` : `${baseUrl}?page=${p}`;

        console.log(`[occ] Página ${p}: ${pageUrl}`);
        await page.goto(pageUrl, { waitUntil: 'networkidle', timeout: 30000 });

        const finalUrl = page.url();
        if (p === 1 && !finalUrl.includes(`de-${roleSlug}`)) {
          console.log(`[occ] WARN: OCC redirigió — URL final no contiene "de-${roleSlug}"`);
          break;
        }

        await page.waitForSelector('[id^="jobcard-"]', { timeout: 15000 })
          .catch(() => console.log(`[occ] Timeout esperando job cards (página ${p})`));

        await this.delay(1000, 2000);

        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const cardData: Array<{ jobId: string; lines: string[] }> = await (page as any).evaluate(() => {
          // @ts-ignore
          const els = Array.from(document.querySelectorAll('[id^="jobcard-"]'));
          // @ts-ignore
          return els.map((el: any) => ({
            jobId: el.getAttribute('data-id') || el.id.replace('jobcard-', ''),
            lines: el.innerText.split('\n').map((s: string) => s.trim()).filter(Boolean),
          }));
        });

        console.log(`[occ] Página ${p}: ${cardData.length} job cards`);
        if (cardData.length === 0) break;

        let newOnPage = 0;
        for (const card of cardData) {
          if (!card.jobId || seenIds.has(card.jobId)) continue;

          const { title, salary, company, location } = parseCardText(card.lines);
          if (!title) continue;

          if (!isTitleRelevant(title, params.role)) {
            console.log(`[occ] Descartando "${title}" (irrelevante para "${params.role}")`);
            continue;
          }

          let description = title;
          let modality = '';
          try {
            const resp = await page.request.get(
              `${DETAIL_BASE}/${card.jobId}/d/j?ipo=41&iapo=1`,
              { timeout: 8000 }
            );
            if (resp.ok()) {
              const json = await resp.json();
              if (json?.o?.ld) description = stripHtml(json.o.ld).slice(0, 2000);
              if (json?.o?.rm) modality = json.o.rm;
            }
          } catch { /* continuar con description = title */ }

          seenIds.add(card.jobId);
          newOnPage++;
          jobs.push({
            title,
            company: company || 'Empresa confidencial',
            description, salary, modality, location,
            url: `${OCC_BASE}/empleo/oferta-${card.jobId}/`,
            portal: 'occ',
            externalId: card.jobId,
          });

          await this.delay(300, 700);
        }

        console.log(`[occ] Página ${p}: ${newOnPage} vacantes nuevas`);
        if (newOnPage === 0) break;
        if (p < maxPages) await this.delay(1500, 2500);
      }

      console.log(`[occ] Total extraídas: ${jobs.length} vacantes`);
    } catch (err) {
      console.error('[occ] Error durante scraping:', err);
    } finally {
      await browser.close();
    }

    return jobs;
  }
}

module.exports = { OccPortal };
