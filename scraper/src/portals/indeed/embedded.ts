export async function readEmbeddedResults(page: any): Promise<any[]> {
  return page.evaluate(() => {
    // @ts-ignore
    const providerData = (window as any).mosaic?.providerData?.['mosaic-provider-jobcards']?.metaData;
    if (!providerData) return [];
    for (const key of Object.keys(providerData)) {
      const results = providerData[key]?.results;
      if (Array.isArray(results) && results.length > 0 && results[0]?.jobkey) {
        return results;
      }
    }
    return [];
  });
}
