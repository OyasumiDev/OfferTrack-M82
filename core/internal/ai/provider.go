package ai

import (
	"context"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
)

// AIProvider es el contrato que todos los proveedores deben cumplir.
// El núcleo nunca importa un proveedor concreto, solo esta interfaz.
type AIProvider interface {
	Analyze(ctx context.Context, req domain.AnalysisRequest) (*domain.AnalysisResult, error)
	AdaptCV(ctx context.Context, req domain.AdaptRequest) (*domain.AdaptResult, error)
	Summarize(ctx context.Context, text string) (string, error)
	Name() string
}
