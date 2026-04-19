// scraper/src/portals/indeed/index.ts
const cheerio = require('cheerio') as typeof import('cheerio');
const { BasePortal } = require('../base.portal');
import type { Element } from 'domhandler';
import type { CheerioAPI } from 'cheerio';
import type { JobRaw, SearchParams } from '../base.portal';

class IndeedPortal extends BasePortal {
  name = 'indeed';
  private baseUrl = 'https://mx.indeed.com';

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    const jobs: JobRaw[] = [];
    const browser = await this.createBrowser();

    try {
      const page = await this.newPage(browser);

      const role = encodeURIComponent(params.role);
      const location = params.location ? encodeURIComponent(params.location) : 'México';
      const url = `${this.baseUrl}/jobs?q=${role}&l=${location}`;

      console.log(`[indeed] Navegando a: ${url}`);
      await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });
      await this.delay();

      await page.waitForSelector('[data-testid="slider_item"], .job_seen_beacon, .tapItem', {
        timeout: 15000,
      }).catch(() => console.log('[indeed] Selector no encontrado, usando HTML disponible'));

      const html = await page.content();
      const $ = cheerio.load(html);

      $('[data-testid="slider_item"], .job_seen_beacon, .tapItem').each((_i: number, el: Element) => {
        if (jobs.length >= 20) return false;
        const job = this.extractJobFromElement($, el);
        if (job) jobs.push(job);
      });

      console.log(`[indeed] Extraídas ${jobs.length} vacantes`);
    } catch (err) {
      console.error('[indeed] Error durante scraping:', err);
    } finally {
      await browser.close();
    }

    return jobs;
  }

  private extractJobFromElement($: CheerioAPI, el: Element): JobRaw | null {
    try {
      const title =
        $('[data-testid="jobTitle"] span, .jobTitle span', el).first().text().trim() ||
        $('h2.jobTitle', el).first().text().trim();

      const company =
        $('[data-testid="company-name"], .companyName', el).first().text().trim();

      const location =
        $('[data-testid="text-location"], .companyLocation', el).first().text().trim();

      const salary =
        $('[data-testid="attribute_snippet_testid"], .salary-snippet', el).first().text().trim() || '';

      const modality =
        $('[data-testid="attribute_snippet_testid"]:last-child, [class*="remote"]', el).first().text().trim() || '';

      const jk =
        (el as any).attribs?.['data-jk'] ||
        $('a[data-jk]', el).first().attr('data-jk') || '';

      const relativeUrl =
        $('a[href*="/rc/clk"], a[href*="/pagead"]', el).first().attr('href') ||
        (jk ? `/rc/clk?jk=${jk}` : '');

      const fullUrl = relativeUrl.startsWith('http')
        ? relativeUrl
        : `${this.baseUrl}${relativeUrl}`;

      if (!title || !relativeUrl) return null;

      return {
        title,
        company: company || 'Empresa confidencial',
        description: $('[class*="summary"], [class*="snippet"]', el).text().trim() || title,
        salary, modality, location,
        url: fullUrl,
        portal: 'indeed',
      };
    } catch {
      return null;
    }
  }
}

module.exports = { IndeedPortal };
