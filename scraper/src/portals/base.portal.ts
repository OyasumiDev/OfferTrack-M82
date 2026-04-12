export interface JobRaw {
  title: string;
  company: string;
  description: string;
  salary?: string;
  modality?: string;
  location?: string;
  url: string;
  portal: string;
}

export interface SearchParams {
  role: string;
  location?: string;
  salaryMin?: number;
  modality?: string;
}

export abstract class BasePortal {
  abstract name: string;
  abstract scrape(params: SearchParams): Promise<JobRaw[]>;

  protected randomDelay(): Promise<void> {
    const min = parseInt(process.env.SCRAPER_DELAY_MIN_MS || "1500");
    const max = parseInt(process.env.SCRAPER_DELAY_MAX_MS || "4000");
    const ms = Math.floor(Math.random() * (max - min + 1)) + min;
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
}
