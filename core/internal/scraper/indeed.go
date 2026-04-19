// core/internal/scraper/indeed.go
// Indeed México — stub pendiente de análisis en vivo del portal
package scraper

import (
	"context"
	"fmt"
)

// ErrIndeedNotImplemented se devuelve hasta que se complete el análisis en vivo de Indeed.
var ErrIndeedNotImplemented = fmt.Errorf("indeed: scraper pendiente de implementación — analizar portal primero")

type indeedScraper struct {
	scraperURL string
}

func newIndeedScraper(scraperURL string) *indeedScraper {
	return &indeedScraper{scraperURL: scraperURL}
}

func (s *indeedScraper) Name() Portal { return PortalIndeed }

func (s *indeedScraper) BuildURL(_ SearchParams) (string, error) {
	return "", ErrIndeedNotImplemented
}

func (s *indeedScraper) SupportsFilter(_ FilterKind) bool { return false }

func (s *indeedScraper) FetchListings(_ context.Context, _ SearchParams) ([]Listing, error) {
	return nil, ErrIndeedNotImplemented
}

func (s *indeedScraper) FetchDetail(_ context.Context, _ string) (*Listing, error) {
	return nil, ErrIndeedNotImplemented
}
