// core/internal/services/orchestrator.go
package services

import (
	"context"
	"fmt"
	"time"

	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/scraper"

	// embed endpoint — sigue en Node.js
	"bytes"
	"encoding/json"
	"net/http"
)

// Orchestrator coordina el scraper y Qdrant.
type Orchestrator struct {
	db         *db.QdrantClient
	scraperURL string
	httpClient *http.Client
}

func NewOrchestrator(db *db.QdrantClient, scraperBaseURL string) *Orchestrator {
	return &Orchestrator{
		db:         db,
		scraperURL: scraperBaseURL,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

// SearchParams son los parámetros que vienen del CLI.
type SearchParams struct {
	Role       string
	Location   string   // "Ciudad" o "Ciudad, Estado"
	SalaryMin  int
	SalaryMax  int
	Modality   string
	Portals    []string
	MaxResults int
	MaxPages   int
}

// Search llama al scraper vía la interfaz Go, guarda en Qdrant y devuelve resultados re-rankeados.
func (o *Orchestrator) Search(ctx context.Context, p SearchParams) ([]*domain.Job, error) {
	if p.MaxResults == 0 {
		p.MaxResults = 20
	}
	if len(p.Portals) == 0 {
		p.Portals = []string{"occ"}
	}

	// Convertir SearchParams del CLI a scraper.SearchParams
	sp := cliToScraperParams(p)

	var allListings []scraper.Listing

	for _, portalName := range p.Portals {
		portal := scraper.Portal(portalName)
		s, err := scraper.NewScraper(portal, o.scraperURL)
		if err != nil {
			fmt.Printf("[orchestrator] Portal no reconocido: %q — %v\n", portalName, err)
			continue
		}

		fmt.Printf("[orchestrator] Scrapeando %s — keywords: %q, city: %q\n", portalName, sp.Keywords, sp.City)
		listings, err := s.FetchListings(ctx, sp)
		if err != nil {
			fmt.Printf("[orchestrator] Warning: %s falló (%v) — omitiendo\n", portalName, err)
			continue
		}
		allListings = append(allListings, listings...)
	}

	// Convertir Listing → domain.Job y guardar en Qdrant
	allJobs := make([]*domain.Job, 0, len(allListings))
	for _, l := range allListings {
		job := listingToDomain(l)
		allJobs = append(allJobs, job)
		if err := o.db.UpsertJob(ctx, job); err != nil {
			fmt.Printf("[orchestrator] Warning: no se pudo guardar %q → %v\n", job.Title, err)
		}
	}

	// Re-ranking semántico en memoria usando embeddings de los listings
	roleVec, err := o.embedQuery(ctx, p.Role)
	if err != nil {
		fmt.Printf("[orchestrator] Warning: embed fallido (%v) — devolviendo orden original\n", err)
		return allJobs, nil
	}

	type scored struct {
		job   *domain.Job
		score float32
		emb   []float32
	}
	var candidates []scored
	for i, l := range allListings {
		if len(l.Embedding) == 0 {
			continue
		}
		sim := cosineSim(roleVec, l.Embedding)
		if sim >= 0.4 {
			candidates = append(candidates, scored{allJobs[i], sim, l.Embedding})
		}
	}

	for i := 1; i < len(candidates); i++ {
		for j := i; j > 0 && candidates[j].score > candidates[j-1].score; j-- {
			candidates[j], candidates[j-1] = candidates[j-1], candidates[j]
		}
	}

	ranked := make([]*domain.Job, 0, len(candidates))
	for _, c := range candidates {
		ranked = append(ranked, c.job)
	}

	fmt.Printf("[orchestrator] Re-ranking: %d scrapeadas → %d relevantes (cosine ≥ 0.4)\n",
		len(allJobs), len(ranked))

	if len(ranked) == 0 {
		fmt.Println("[orchestrator] Sin resultados relevantes — devolviendo todas las scrapeadas")
		return allJobs, nil
	}
	return ranked, nil
}

// cliToScraperParams convierte SearchParams del CLI a scraper.SearchParams.
// Parsea "Ciudad, Estado" en los campos separados State/City.
func cliToScraperParams(p SearchParams) scraper.SearchParams {
	sp := scraper.SearchParams{
		Keywords:  p.Role,
		SalaryMin: p.SalaryMin,
		SalaryMax: p.SalaryMax,
		MaxPages:  p.MaxPages,
	}

	switch p.Modality {
	case "remoto", "remote":
		sp.Modality = scraper.ModalityRemoto
	case "hibrido", "hybrid":
		sp.Modality = scraper.ModalityHibrido
	case "presencial", "onsite":
		sp.Modality = scraper.ModalityPresencial
	}

	if p.Location != "" {
		parts := splitLocation(p.Location)
		if len(parts) >= 2 {
			sp.City = parts[0]
			sp.State = parts[1]
		} else {
			sp.City = parts[0]
		}
	}

	return sp
}

func splitLocation(loc string) []string {
	for i, c := range loc {
		if c == ',' {
			city := trim(loc[:i])
			state := trim(loc[i+1:])
			return []string{city, state}
		}
	}
	return []string{trim(loc)}
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// listingToDomain convierte un scraper.Listing a domain.Job.
func listingToDomain(l scraper.Listing) *domain.Job {
	modality := string(l.Modality)
	if modality == "" || modality == string(scraper.ModalityUnknown) {
		modality = "onsite"
	}
	postedAt := l.PostedAt
	if postedAt.IsZero() {
		postedAt = time.Now()
	}
	return &domain.Job{
		ID:              l.ID,
		Title:           l.Title,
		Company:         l.Company,
		Description:     l.Description,
		SalaryMin:       l.SalaryMin,
		SalaryMax:       l.SalaryMax,
		Currency:        "MXN",
		Modality:        modality,
		Location:        l.Location,
		Portal:          string(l.Portal),
		URL:             l.SourceURL,
		PostedAt:        postedAt,
		ScrapedAt:       time.Now(),
		EmbeddingVector: l.Embedding,
	}
}

// cosineSim calcula similitud coseno entre dos vectores de igual dimensión.
func cosineSim(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float32
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (sqrt32(normA) * sqrt32(normB))
}

func sqrt32(x float32) float32 {
	if x <= 0 {
		return 0
	}
	g := x / 2
	for i := 0; i < 8; i++ {
		g = (g + x/g) / 2
	}
	return g
}

// embedQuery obtiene el vector del rol desde el scraper Node.js (/embed/query).
func (o *Orchestrator) embedQuery(ctx context.Context, text string) ([]float32, error) {
	body, _ := json.Marshal(map[string]string{"text": text})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.scraperURL+"/embed/query", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embed/query: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("embed/query respondió %d", resp.StatusCode)
	}

	var result struct {
		Embedding []float32 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("embed/query: parse error → %w", err)
	}
	return result.Embedding, nil
}

// ListSaved lista vacantes ya guardadas en Qdrant.
func (o *Orchestrator) ListSaved(ctx context.Context, limit uint64) ([]*domain.Job, error) {
	return o.db.ListJobs(ctx, limit, 0)
}
