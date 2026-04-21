import type { Element } from 'domhandler';

export function extractJobKey(el: Element): string {
  return (el as any).attribs?.['data-jk'] || '';
}

export async function extractAllJobKeys(page: any): Promise<string[]> {
  const keys: string[] = await page.evaluate(() => {
    // @ts-ignore
    const els = Array.from(document.querySelectorAll('[data-jk]'));
    return els.map((el: any) => el.getAttribute('data-jk')).filter(Boolean);
  });
  return [...new Set(keys)];
}
