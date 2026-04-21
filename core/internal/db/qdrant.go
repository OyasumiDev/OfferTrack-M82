// core/internal/db/qdrant.go
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/qdrant/go-client/qdrant"
)

// QdrantClient envuelve el cliente gRPC de Qdrant con helpers para el dominio.
type QdrantClient struct {
	client         *qdrant.Client
	collectionJobs string
}

// NewQdrantClient crea la conexión gRPC a Qdrant.
func NewQdrantClient(host string, port int, collectionJobs string) (*QdrantClient, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:                 host,
		Port:                 port,
		SkipCompatibilityCheck: true,
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: error conectando a %s:%d → %w", host, port, err)
	}
	return &QdrantClient{client: client, collectionJobs: collectionJobs}, nil
}

// RawClient expone el cliente subyacente para operaciones de administración (InitCollections).
func (q *QdrantClient) RawClient() *qdrant.Client {
	return q.client
}

// Ping verifica que Qdrant está vivo.
func (q *QdrantClient) Ping(ctx context.Context) error {
	_, err := q.client.HealthCheck(ctx)
	return err
}

// UpsertJob guarda (o actualiza) una vacante en Qdrant.
func (q *QdrantClient) UpsertJob(ctx context.Context, job *domain.Job) error {
	if len(job.EmbeddingVector) == 0 {
		return fmt.Errorf("qdrant: job %q no tiene vector de embedding", job.ID)
	}

	payload := map[string]any{
		"title":        job.Title,
		"company":      job.Company,
		"description":  job.Description,
		"salary_min":   job.SalaryMin,
		"salary_max":   job.SalaryMax,
		"currency":     job.Currency,
		"modality":     job.Modality,
		"location":     job.Location,
		"portal":       job.Portal,
		"url":          job.URL,
		"posted_at":    job.PostedAt.Format(time.RFC3339),
		"scraped_at":   job.ScrapedAt.Format(time.RFC3339),
		"compat_score": job.CompatScore,
	}

	_, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: q.collectionJobs,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDUUID(job.ID),
				Vectors: qdrant.NewVectors(job.EmbeddingVector...),
				Payload: qdrant.NewValueMap(payload),
			},
		},
	})
	if err != nil {
		return fmt.Errorf("qdrant: error upserting job %q → %w", job.ID, err)
	}
	return nil
}

// SearchSimilar busca las vacantes más parecidas a un vector de consulta.
// scoreThreshold 0 = sin filtro; >0 descarta resultados por debajo del umbral (cosine).
func (q *QdrantClient) SearchSimilar(ctx context.Context, vector []float32, limit uint64, scoreThreshold float32) ([]*domain.Job, error) {
	qp := &qdrant.QueryPoints{
		CollectionName: q.collectionJobs,
		Query:          qdrant.NewQuery(vector...),
		Limit:          qdrant.PtrOf(limit),
		WithPayload:    qdrant.NewWithPayload(true),
	}
	if scoreThreshold > 0 {
		qp.ScoreThreshold = qdrant.PtrOf(scoreThreshold)
	}
	results, err := q.client.Query(ctx, qp)
	if err != nil {
		return nil, fmt.Errorf("qdrant: error en búsqueda similar → %w", err)
	}

	jobs := make([]*domain.Job, 0, len(results))
	for _, r := range results {
		job := pointToJob(r.Id, r.Payload)
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// GetJob recupera una vacante por su ID.
func (q *QdrantClient) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	results, err := q.client.Get(ctx, &qdrant.GetPoints{
		CollectionName: q.collectionJobs,
		Ids:            []*qdrant.PointId{qdrant.NewIDUUID(id)},
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: error obteniendo job %q → %w", id, err)
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("qdrant: job %q no encontrado", id)
	}
	return pointToJob(results[0].Id, results[0].Payload), nil
}

// ListJobs lista vacantes con paginación simple (scroll).
// Limit es uint32 en ScrollPoints.
func (q *QdrantClient) ListJobs(ctx context.Context, limit uint64, _ uint64) ([]*domain.Job, error) {
	results, err := q.client.Scroll(ctx, &qdrant.ScrollPoints{
		CollectionName: q.collectionJobs,
		Limit:          qdrant.PtrOf(uint32(limit)),
		WithPayload:    qdrant.NewWithPayload(true),
	})
	if err != nil {
		return nil, fmt.Errorf("qdrant: error listando jobs → %w", err)
	}

	jobs := make([]*domain.Job, 0, len(results))
	for _, r := range results {
		job := pointToJob(r.Id, r.Payload)
		jobs = append(jobs, job)
	}
	return jobs, nil
}

// CountJobs retorna el total de vacantes guardadas.
func (q *QdrantClient) CountJobs(ctx context.Context) (uint64, error) {
	result, err := q.client.Count(ctx, &qdrant.CountPoints{
		CollectionName: q.collectionJobs,
	})
	if err != nil {
		return 0, fmt.Errorf("qdrant: error contando jobs → %w", err)
	}
	return result, nil
}

// pointToJob convierte un punto de Qdrant a domain.Job.
func pointToJob(id *qdrant.PointId, payload map[string]*qdrant.Value) *domain.Job {
	job := &domain.Job{
		ID:          id.GetUuid(),
		Title:       strVal(payload, "title"),
		Company:     strVal(payload, "company"),
		Description: strVal(payload, "description"),
		Currency:    strVal(payload, "currency"),
		Modality:    strVal(payload, "modality"),
		Location:    strVal(payload, "location"),
		Portal:      strVal(payload, "portal"),
		URL:         strVal(payload, "url"),
	}

	if v, ok := payload["salary_min"]; ok {
		job.SalaryMin = int(v.GetIntegerValue())
	}
	if v, ok := payload["salary_max"]; ok {
		job.SalaryMax = int(v.GetIntegerValue())
	}
	if v, ok := payload["compat_score"]; ok {
		job.CompatScore = int(v.GetIntegerValue())
	}
	if s := strVal(payload, "scraped_at"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			job.ScrapedAt = t
		}
	}
	if s := strVal(payload, "posted_at"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			job.PostedAt = t
		}
	}
	return job
}

// UpdateCompatScore actualiza solo el campo compat_score de una vacante en Qdrant.
func (q *QdrantClient) UpdateCompatScore(ctx context.Context, jobID string, score int) error {
	_, err := q.client.SetPayload(ctx, &qdrant.SetPayloadPoints{
		CollectionName: q.collectionJobs,
		Payload:        qdrant.NewValueMap(map[string]any{"compat_score": score}),
		PointsSelector: &qdrant.PointsSelector{
			PointsSelectorOneOf: &qdrant.PointsSelector_Points{
				Points: &qdrant.PointsIdsList{
					Ids: []*qdrant.PointId{qdrant.NewIDUUID(jobID)},
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("qdrant: UpdateCompatScore %q → %w", jobID, err)
	}
	return nil
}

func strVal(payload map[string]*qdrant.Value, key string) string {
	if v, ok := payload[key]; ok {
		return v.GetStringValue()
	}
	return ""
}
