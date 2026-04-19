package domain

// Analysis es el resultado de un análisis guardado en Qdrant.
type Analysis struct {
	JobID              string   `json:"job_id"`
	CompatibilityScore int      `json:"compatibility_score"`
	Strengths          []string `json:"strengths"`
	Gaps               []string `json:"gaps"`
	SalaryEstimate     string   `json:"salary_estimate"`
	Recommendation     string   `json:"recommendation"` // "apply" | "consider" | "discard"
	RawAnalysis        string   `json:"raw_analysis"`
}

// AnalysisRequest es el payload que se envía al proveedor IA para analizar.
type AnalysisRequest struct {
	JobDescription string
	UserProfile    string
	UserCV         string
}

// AnalysisResult es la respuesta estructurada del proveedor IA.
type AnalysisResult struct {
	CompatibilityScore int
	Strengths          []string
	Gaps               []string
	SalaryEstimate     string
	Recommendation     string // "apply" | "consider" | "discard"
	RawAnalysis        string
}

// AdaptRequest es el payload que se envía al proveedor IA para adaptar el CV.
type AdaptRequest struct {
	JobDescription string
	BaseCV         string
	UserProfile    string
}

// AdaptResult es la respuesta del proveedor IA con el CV adaptado.
type AdaptResult struct {
	AdaptedCV     string
	Changes       []string
	KeywordsAdded []string
}
