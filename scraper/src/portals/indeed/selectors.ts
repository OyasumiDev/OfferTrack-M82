export const INDEED_SELECTORS = {
  cardContainer: 'div.job_seen_beacon',
  cardContainerFallback: 'a[data-jk]',
  jobKeyAttr: 'data-jk',
  title: 'h2.jobTitle span[title]',
  titleFallback: 'h2.jobTitle span',
  company: '[data-testid="company-name"]',
  location: '[data-testid="text-location"]',
  salary: '.salary-snippet-container',
  attributeSnippet: '[data-testid*="attribute_snippet_testid"]',
  paginationNext: '[data-testid="pagination-page-next"]',
  detail: {
    title: '[data-testid="simpler-jobTitle"], h1',
    company: '[data-testid="inlineHeader-companyName"]',
    location: '[data-testid="inlineHeader-companyLocation"]',
    description: '#jobDescriptionText',
  },
} as const;
