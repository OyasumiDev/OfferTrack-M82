// core/internal/services/analyzer.go
package services

import (
	"context"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/ai"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/domain"
)

// AnalyzerService analiza vacantes contra el perfil del usuario usando IA.
type AnalyzerService struct {
	provider ai.AIProvider
	db       *db.QdrantClient
}

// NewAnalyzerService crea un AnalyzerService listo para usar.
func NewAnalyzerService(provider ai.AIProvider, db *db.QdrantClient) *AnalyzerService {
	return &AnalyzerService{provider: provider, db: db}
}

// AnalyzeJob analiza la compatibilidad de una vacante con el perfil del usuario.
// Actualiza compat_score en Qdrant si el análisis tiene éxito.
func (s *AnalyzerService) AnalyzeJob(ctx context.Context, job *domain.Job, profile, cv string) (*domain.Analysis, error) {
	req := domain.AnalysisRequest{
		JobDescription: fmt.Sprintf("%s — %s\n\n%s", job.Title, job.Company, job.Description),
		UserProfile:    profile,
		UserCV:         cv,
	}

	result, err := s.provider.Analyze(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("analyzer: %w", err)
	}

	analysis := &domain.Analysis{
		JobID:              job.ID,
		CompatibilityScore: result.CompatibilityScore,
		Strengths:          result.Strengths,
		Gaps:               result.Gaps,
		SalaryEstimate:     result.SalaryEstimate,
		Recommendation:     result.Recommendation,
		RawAnalysis:        result.RawAnalysis,
	}

	// Persistir compat_score en Qdrant de forma no bloqueante
	if updateErr := s.db.UpdateCompatScore(ctx, job.ID, result.CompatibilityScore); updateErr != nil {
		fmt.Printf("[analyzer] warning: no se pudo actualizar compat_score para %q → %v\n", job.ID, updateErr)
	}

	return analysis, nil
}

// AnalyzeBatch analiza varias vacantes en secuencia y retorna solo las que superan minScore.
func (s *AnalyzerService) AnalyzeBatch(ctx context.Context, jobs []*domain.Job, profile, cv string, minScore int) ([]*domain.Analysis, error) {
	var results []*domain.Analysis
	for _, job := range jobs {
		a, err := s.AnalyzeJob(ctx, job, profile, cv)
		if err != nil {
			fmt.Printf("[analyzer] error en %q: %v — se omite\n", job.Title, err)
			continue
		}
		if a.CompatibilityScore >= minScore {
			results = append(results, a)
		}
	}
	return results, nil
}
