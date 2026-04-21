// core/internal/scraper/computrabajo.go
// Computrabajo — implementación de Scraper, aislada de los otros portales
package scraper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type computrabajoScraper struct {
	scraperURL string
	http       *http.Client
}

func newComputrabajoScraper(scraperURL string) *computrabajoScraper {
	return &computrabajoScraper{
		scraperURL: scraperURL,
		http:       &http.Client{Timeout: 180 * time.Second},
	}
}

func (s *computrabajoScraper) Name() Portal { return PortalComputrabajo }

func (s *computrabajoScraper) BuildURL(params SearchParams) (string, error) {
	return buildComputrabajoURL(params), nil
}

func (s *computrabajoScraper) SupportsFilter(f FilterKind) bool {
	return f == FilterSalary
}

func (s *computrabajoScraper) FetchListings(ctx context.Context, params SearchParams) ([]Listing, error) {
	location := params.City
	if location == "" {
		location = params.State
	}
	maxPages := params.MaxPages
	if maxPages == 0 {
		maxPages = 3
	}

	payload := map[string]any{
		"role":     params.Keywords,
		"location": location,
		"portals":  []string{"computrabajo"},
		"maxPages": maxPages,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.scraperURL+"/scrape", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("computrabajo: creando request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("computrabajo: llamando al scraper: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("computrabajo: scraper respondió %d", resp.StatusCode)
	}

	var result nodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("computrabajo: parseando respuesta: %w", err)
	}

	listings := make([]Listing, 0, len(result.Jobs))
	for _, j := range result.Jobs {
		listings = append(listings, nodeJobToListing(j, PortalComputrabajo))
	}
	return listings, nil
}

func (s *computrabajoScraper) FetchDetail(_ context.Context, _ string) (*Listing, error) {
	return nil, fmt.Errorf("computrabajo: FetchDetail no implementado")
}

// buildComputrabajoURL construye la URL canónica: /trabajo-de-{puesto}-en-{ciudad}
func buildComputrabajoURL(p SearchParams) string {
	const base = "https://mx.computrabajo.com"
	roleSlug := ctSlugify(p.Keywords)
	city := p.City
	if city == "" {
		city = p.State
	}
	if city == "" {
		return fmt.Sprintf("%s/trabajo-de-%s", base, roleSlug)
	}
	return fmt.Sprintf("%s/trabajo-de-%s-en-%s", base, roleSlug, ctSlugify(city))
}

// ctSlugify — algoritmo propio de Computrabajo (idéntico al de OCC pero independiente)
func ctSlugify(s string) string {
	return occSlugify(s) // mismo algoritmo — en producción puede divergir si Computrabajo cambia
}

// Timestamp dummy para evitar import cycle en time.Now
var _ = time.Now
