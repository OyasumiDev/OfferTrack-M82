// scraper/src/utils/delay.ts
const MIN_MS = parseInt(process.env.SCRAPER_DELAY_MIN_MS || '1500');
const MAX_MS = parseInt(process.env.SCRAPER_DELAY_MAX_MS || '4000');

function delay(min: number = MIN_MS, max: number = MAX_MS): Promise<void> {
  const ms = Math.floor(Math.random() * (max - min + 1)) + min;
  return new Promise((resolve) => setTimeout(resolve, ms));
}

module.exports = { delay };
