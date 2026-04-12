import { BasePortal, JobRaw, SearchParams } from "./base.portal";

export class ComputrabajoPortal extends BasePortal {
  name = "computrabajo";

  async scrape(params: SearchParams): Promise<JobRaw[]> {
    // TODO: implementar scraping de Computrabajo
    return [];
  }
}
