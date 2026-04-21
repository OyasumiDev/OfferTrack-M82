// core/internal/scraper/occ.go
// OCC Mundial — implementación de Scraper que delega el browser al servicio Node.js
package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type occScraper struct {
	scraperURL string
	http       *http.Client
}

func newOccScraper(scraperURL string) *occScraper {
	return &occScraper{
		scraperURL: scraperURL,
		http:       &http.Client{Timeout: 180 * time.Second},
	}
}

func (s *occScraper) Name() Portal { return PortalOCC }

func (s *occScraper) BuildURL(params SearchParams) (string, error) {
	return buildOccURL(params), nil
}

func (s *occScraper) SupportsFilter(f FilterKind) bool {
	return f == FilterSalary || f == FilterModality
}

func (s *occScraper) FetchListings(ctx context.Context, params SearchParams) ([]Listing, error) {
	// Componer location para el Node.js scraper (acepta "Ciudad, Estado" o solo "Ciudad")
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
		"portals":   []string{"occ"},
		"maxPages":  maxPages,
	}

	return s.callNodeScraper(ctx, payload)
}

func (s *occScraper) FetchDetail(_ context.Context, _ string) (*Listing, error) {
	// El detalle se obtiene inline durante FetchListings vía la API oferta.occ.com.mx
	return nil, fmt.Errorf("occ: FetchDetail se resuelve inline en FetchListings")
}

// callNodeScraper envía el payload al servicio Node.js y convierte la respuesta a []Listing.
func (s *occScraper) callNodeScraper(ctx context.Context, payload map[string]any) ([]Listing, error) {
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.scraperURL+"/scrape", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("occ: creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("occ: llamando al scraper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("occ: scraper respondió %d", resp.StatusCode)
	}

	var result nodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("occ: parseando respuesta: %w", err)
	}

	listings := make([]Listing, 0, len(result.Jobs))
	for _, j := range result.Jobs {
		listings = append(listings, nodeJobToListing(j, PortalOCC))
	}
	return listings, nil
}

// ── tipos compartidos para la respuesta del Node.js scraper ────────────────────

type nodeResponse struct {
	Jobs  []nodeJob `json:"jobs"`
	Count int       `json:"count"`
}

type nodeJob struct {
	ID                 string    `json:"id"`
	ExternalID         string    `json:"externalId"`
	Title              string    `json:"title"`
	Company            string    `json:"company"`
	Description        string    `json:"description"`
	SalaryMin          int       `json:"salaryMin"`
	SalaryMax          int       `json:"salaryMax"`
	Salary             string    `json:"salary"`
	ModalityNormalized string    `json:"modalityNormalized"`
	Location           string    `json:"location"`
	URL                string    `json:"url"`
	ScrapedAt          string    `json:"scrapedAt"`
	Embedding          []float32 `json:"embedding"`
}

func nodeJobToListing(j nodeJob, portal Portal) Listing {
	modality := Modality(j.ModalityNormalized)
	if modality == "" {
		modality = ModalityUnknown
	}
	return Listing{
		Portal:      portal,
		ExternalID:  j.ExternalID,
		Title:       j.Title,
		Company:     j.Company,
		Description: j.Description,
		Location:    j.Location,
		Modality:    modality,
		SalaryMin:   j.SalaryMin,
		SalaryMax:   j.SalaryMax,
		SalaryRaw:   j.Salary,
		SourceURL:   j.URL,
		FetchedAt:   time.Now(),
		ID:          j.ID,
		Embedding:   j.Embedding,
	}
}
