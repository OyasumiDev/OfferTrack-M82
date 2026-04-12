package ai

import "context"

// AIProvider es el contrato que todos los proveedores deben cumplir.
// El nucleo nunca importa un proveedor concreto, solo esta interfaz.
type AIProvider interface {
	Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResult, error)
	AdaptCV(ctx context.Context, req AdaptRequest) (*AdaptResult, error)
	Summarize(ctx context.Context, text string) (string, error)
	Name() string
}

type AnalysisRequest struct {
	JobDescription string
	UserProfile    string
	UserCV         string
}

type AnalysisResult struct {
	CompatibilityScore int
	Strengths          []string
	Gaps               []string
	SalaryEstimate     string
	Recommendation     string // "apply" | "consider" | "discard"
	RawAnalysis        string
}

type AdaptRequest struct {
	JobDescription string
	BaseCV         string
	UserProfile    string
}

type AdaptResult struct {
	AdaptedCV     string
	Changes       []string
	KeywordsAdded []string
}
