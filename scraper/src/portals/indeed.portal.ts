import { BasePortal, JobRaw, SearchParams } from "./base.portal";

export class IndeedPortal extends BasePortal {
  name = "indeed";

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    // TODO: implementar scraping de Indeed Mexico
    return [];
  }
}
