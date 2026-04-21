package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type indeedScraper struct {
	scraperURL string
	http       *http.Client
}

func newIndeedScraper(scraperURL string) *indeedScraper {
	return &indeedScraper{
		scraperURL: scraperURL,
		http:       &http.Client{Timeout: 180 * time.Second},
	}
}

func (i *indeedScraper) Name() Portal { return PortalIndeed }

func (i *indeedScraper) SupportsFilter(f FilterKind) bool {
	return f == FilterSalary || f == FilterModality || f == FilterDate
}

func (i *indeedScraper) BuildURL(params SearchParams) (string, error) {
	return BuildIndeedURL(cliToIndeedQuery(params, 1)), nil
}

func (i *indeedScraper) FetchListings(ctx context.Context, params SearchParams) ([]Listing, error) {
	location := ""
	if params.City != "" && params.State != "" {
		location = params.City + ", " + params.State
	} else if params.City != "" {
		location = params.City
	} else if params.State != "" {
		location = params.State
	}

	maxPages := params.MaxPages
	if maxPages == 0 {
		maxPages = 3
	}

	payload := map[string]any{
		"role":      params.Keywords,
		"location":  location,
		"salaryMin": params.SalaryMin,
		"salaryMax": params.SalaryMax,
		"modality":  string(params.Modality),
		"portals":   []string{"indeed"},
		"maxPages":  maxPages,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, i.scraperURL+"/scrape", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("indeed: creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := i.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("indeed: llamando al scraper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("indeed: scraper respondió %d", resp.StatusCode)
	}

	var result nodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("indeed: parseando respuesta: %w", err)
	}

	listings := make([]Listing, 0, len(result.Jobs))
	for _, j := range result.Jobs {
		listings = append(listings, nodeJobToListing(j, PortalIndeed))
	}
	return listings, nil
}

func (i *indeedScraper) FetchDetail(ctx context.Context, jobkey string) (*Listing, error) {
	u := BuildIndeedDetailURL(jobkey)
	html, err := i.fetchRaw(ctx, u)
	if err != nil {
		return &Listing{ExternalID: jobkey, SourceURL: u, FetchedAt: time.Now(), Portal: PortalIndeed}, nil
	}
	_ = html
	// TODO(detail-enrichment): parse window._initialData for full job description
	return &Listing{ExternalID: jobkey, SourceURL: u, FetchedAt: time.Now(), Portal: PortalIndeed}, nil
}

func (i *indeedScraper) fetchRaw(ctx context.Context, u string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-MX,es;q=0.9,en;q=0.8")
	resp, err := i.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("indeed: HTTP %d para %s", resp.StatusCode, u)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func cliToIndeedQuery(params SearchParams, page int) IndeedQuery {
	q := IndeedQuery{
		Puesto: params.Keywords,
		Page:   page,
	}
	q.Ciudad = params.City
	q.Estado = params.State
	if params.SalaryMin > 0 {
		q.SalaryMinAnnualMXN = params.SalaryMin
	}
	switch params.Modality {
	case ModalityRemoto:
		q.Modality = "remote"
	case ModalityHibrido:
		q.Modality = "hybrid"
	}
	return q
}
