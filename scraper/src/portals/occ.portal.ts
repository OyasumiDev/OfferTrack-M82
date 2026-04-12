import { BasePortal, JobRaw, SearchParams } from "./base.portal";

export class OccPortal extends BasePortal {
  name = "occ";

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    // TODO: implementar scraping de OCC Mundial
    return [];
  }
}
