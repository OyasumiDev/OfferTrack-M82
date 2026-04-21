// core/internal/services/cv_adapter.go
package services

import (
	"context"
	"fmt"

	"github.com/QUERTY/OfferTrack-M82/internal/ai"
	"github.com/QUERTY/OfferTrack-M82/internal/domain"
)

// CVAdapterService adapta el CV base del usuario a una vacante específica.
type CVAdapterService struct {
	provider ai.AIProvider
}

// NewCVAdapterService crea un CVAdapterService listo para usar.
func NewCVAdapterService(provider ai.AIProvider) *CVAdapterService {
	return &CVAdapterService{provider: provider}
}

// Adapt genera una versión del CV adaptada a la vacante indicada.
// baseCV es el CV en texto plano o Markdown. profile es el perfil del usuario.
func (s *CVAdapterService) Adapt(ctx context.Context, job *domain.Job, baseCV, profile string) (*domain.AdaptResult, error) {
	req := domain.AdaptRequest{
		JobDescription: fmt.Sprintf("%s — %s\n\n%s", job.Title, job.Company, job.Description),
		BaseCV:         baseCV,
		UserProfile:    profile,
	}
	result, err := s.provider.AdaptCV(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cv_adapter: %w", err)
	}
	return result, nil
}
