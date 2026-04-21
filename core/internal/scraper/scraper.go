// core/internal/scraper/scraper.go
// Interfaces, tipos y dispatcher — sin lógica de parseo ni scraping
package scraper

import (
	"context"
	"fmt"
	"time"
)

// Portal identifica una bolsa de trabajo soportada.
type Portal string

const (
	PortalOCC          Portal = "occ"
	PortalComputrabajo Portal = "computrabajo"
	PortalIndeed       Portal = "indeed"
	PortalAll          Portal = "all"
)

// AllPortals lista los portales individuales (excluye "all").
var AllPortals = []Portal{PortalOCC, PortalComputrabajo, PortalIndeed}

// Modality representa la modalidad de trabajo de una vacante.
type Modality string

const (
	ModalityPresencial Modality = "presencial"
	ModalityHibrido    Modality = "hibrido"
	ModalityRemoto     Modality = "remoto"
	ModalityUnknown    Modality = "unknown"
)

// FilterKind identifica tipos de filtro que un portal puede o no soportar.
type FilterKind string

const (
	FilterSalary   FilterKind = "salary"
	FilterModality FilterKind = "modality"
	FilterDate     FilterKind = "date"
)

// SearchParams son los parámetros de búsqueda normalizados del CLI.
// Cada scraper decide qué campos puede honrar.
type SearchParams struct {
	Keywords  string   // "desarrollador .NET"
	State     string   // "Nuevo León"
	City      string   // "Monterrey" (opcional)
	SalaryMin int      // 0 = sin límite
	SalaryMax int      // 0 = sin límite
	Modality  Modality // "" = todas
	MaxPages  int      // 0 = default del portal (3)
}

// Listing es el contrato de salida de todos los scrapers.
// Si un campo no está disponible, usa el zero value — nunca nil ni dato inventado.
type Listing struct {
	Portal      Portal    `json:"portal"`
	ExternalID  string    `json:"external_id"`
	Title       string    `json:"title"`
	Company     string    `json:"company"`
	Description string    `json:"description"`
	Location    string    `json:"location"`    // "Monterrey, Nuevo León"
	Modality    Modality  `json:"modality"`
	SalaryMin   int       `json:"salary_min"`  // 0 = no especificado
	SalaryMax   int       `json:"salary_max"`
	SalaryRaw   string    `json:"salary_raw"`  // texto original
	PostedAt    time.Time `json:"posted_at"`
	SourceURL   string    `json:"source_url"`
	FetchedAt   time.Time `json:"fetched_at"`
	// Campos internos del pipeline (no parte del contrato de portal)
	ID        string    `json:"id,omitempty"`
	Embedding []float32 `json:"embedding,omitempty"`
}

// Scraper es la interfaz que toda implementación por portal debe cumplir.
type Scraper interface {
	Name() Portal
	BuildURL(params SearchParams) (string, error)
	FetchListings(ctx context.Context, params SearchParams) ([]Listing, error)
	FetchDetail(ctx context.Context, id string) (*Listing, error)
	SupportsFilter(f FilterKind) bool
}

// NewScraper devuelve la implementación correcta según el portal.
// scraperBaseURL es la URL base del servicio Node.js (ej. "http://localhost:3001").
func NewScraper(p Portal, scraperBaseURL string) (Scraper, error) {
	switch p {
	case PortalOCC:
		return newOccScraper(scraperBaseURL), nil
	case PortalComputrabajo:
		return newComputrabajoScraper(scraperBaseURL), nil
	case PortalIndeed:
		return newIndeedScraper(scraperBaseURL), nil
	}
	return nil, fmt.Errorf("portal no reconocido: %q — válidos: occ, computrabajo, indeed, all", p)
}
