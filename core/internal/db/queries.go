// core/internal/db/queries.go
package db

import (
	"context"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/qdrant/go-client/qdrant"
)

// JobFilter agrupa los filtros opcionales para buscar vacantes.
type JobFilter struct {
	Portal    string // "occ" | "computrabajo" | "indeed" | "" (sin filtro)
	Modality  string // "remote" | "hybrid" | "onsite" | "" (sin filtro)
	SalaryMin int    // 0 = sin filtro
}

// FilterJobs lista vacantes aplicando filtros de payload en Qdrant.
func (q *QdrantClient) FilterJobs(ctx context.Context, f JobFilter, limit uint64) ([]*domain.Job, error) {
	filter := buildFilter(f)

	results, err := q.client.Scroll(ctx, &qdrant.ScrollPoints{
		CollectionName: q.collectionJobs,
		Filter:         filter,
		Limit:          qdrant.PtrOf(uint32(limit)),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: error filtrando jobs → %w", err)
	}

	jobs := make([]*domain.Job, 0, len(results))
	for _, r := range results {
		jobs = append(jobs, pointToJob(r.Id, r.Payload))
	}
	return jobs, nil
}

// SearchSimilarFiltered busca por similitud vectorial y aplica filtros de payload.
func (q *QdrantClient) SearchSimilarFiltered(ctx context.Context, vector []float32, f JobFilter, limit uint64) ([]*domain.Job, error) {
	filter := buildFilter(f)

	results, err := q.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: q.collectionJobs,
		Query:          qdrant.NewQuery(vector...),
		Filter:         filter,
		Limit:          qdrant.PtrOf(limit),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: error en búsqueda filtrada → %w", err)
	}

	jobs := make([]*domain.Job, 0, len(results))
	for _, r := range results {
		jobs = append(jobs, pointToJob(r.Id, r.Payload))
	}
	return jobs, nil
}

// buildFilter construye el filtro de Qdrant a partir de un JobFilter.
// Retorna nil si no hay condiciones activas.
func buildFilter(f JobFilter) *qdrant.Filter {
	var conditions []*qdrant.Condition

	if f.Portal != "" {
		conditions = append(conditions, qdrant.NewMatchKeyword("portal", f.Portal))
	}
	if f.Modality != "" {
		conditions = append(conditions, qdrant.NewMatchKeyword("modality", f.Modality))
	}
	if f.SalaryMin > 0 {
		conditions = append(conditions, &qdrant.Condition{
			ConditionOneOf: &qdrant.Condition_Field{
				Field: &qdrant.FieldCondition{
					Key: "salary_min",
					Range: &qdrant.Range{
						Gte: qdrant.PtrOf(float64(f.SalaryMin)),
					},
				},
			},
		})
	}

	if len(conditions) == 0 {
		return nil
	}
	return &qdrant.Filter{Must: conditions}
}
