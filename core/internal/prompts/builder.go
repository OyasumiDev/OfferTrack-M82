// core/internal/prompts/builder.go
package prompts

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
)

const systemAnalysis = `Eres un consultor experto en selección de personal para el mercado laboral de México.
Analiza la compatibilidad entre la vacante y el perfil del candidato.
Responde SOLO con JSON válido, sin texto adicional ni bloques de código markdown:
{"compatibility_score":0,"strengths":[],"gaps":[],"salary_estimate":"","recommendation":"apply|consider|discard","raw_analysis":""}`

const systemAdaptCV = `Eres un redactor experto en CVs para el mercado laboral de México.
Adapta el CV del candidato para maximizar su compatibilidad con la vacante indicada.
Responde SOLO con JSON válido, sin texto adicional ni bloques de código markdown:
{"adapted_cv":"","changes":[],"keywords_added":[]}`

const systemSummarize = `Eres un asistente que resume descripciones de vacantes de forma concisa en español.
Responde SOLO con el resumen en 2-3 oraciones, sin preámbulo.`

// BuildAnalysisMessages retorna (system, user) para analizar compatibilidad.
func BuildAnalysisMessages(req domain.AnalysisRequest) (system, user string) {
	return systemAnalysis, fmt.Sprintf(
		"## Descripción de la vacante\n%s\n\n## Perfil del candidato\n%s\n\n## CV del candidato\n%s",
		req.JobDescription, req.UserProfile, req.UserCV,
	)
}

// BuildAdaptMessages retorna (system, user) para adaptar un CV.
func BuildAdaptMessages(req domain.AdaptRequest) (system, user string) {
	return systemAdaptCV, fmt.Sprintf(
		"## Descripción de la vacante\n%s\n\n## Perfil del candidato\n%s\n\n## CV base\n%s",
		req.JobDescription, req.UserProfile, req.BaseCV,
	)
}

// BuildSummarizeMessages retorna (system, user) para resumir texto.
func BuildSummarizeMessages(text string) (system, user string) {
	return systemSummarize, text
}

// analysisJSON es el esquema JSON que devuelve el LLM para análisis.
type analysisJSON struct {
	CompatibilityScore int      `json:"compatibility_score"`
	Strengths          []string `json:"strengths"`
	Gaps               []string `json:"gaps"`
	SalaryEstimate     string   `json:"salary_estimate"`
	Recommendation     string   `json:"recommendation"`
	RawAnalysis        string   `json:"raw_analysis"`
}

// ParseAnalysisResult parsea la respuesta JSON del LLM para análisis de compatibilidad.
func ParseAnalysisResult(raw string) (*domain.AnalysisResult, error) {
	cleaned := extractJSON(raw)
	var j analysisJSON
	if err := json.Unmarshal([]byte(cleaned), &j); err != nil {
		return nil, fmt.Errorf("ParseAnalysisResult: %w\nrespuesta: %.300s", err, raw)
	}
	if j.RawAnalysis == "" {
		j.RawAnalysis = raw
	}
	return &domain.AnalysisResult{
		CompatibilityScore: j.CompatibilityScore,
		Strengths:          j.Strengths,
		Gaps:               j.Gaps,
		SalaryEstimate:     j.SalaryEstimate,
		Recommendation:     j.Recommendation,
		RawAnalysis:        j.RawAnalysis,
	}, nil
}

// adaptJSON es el esquema JSON que devuelve el LLM para adaptación de CV.
type adaptJSON struct {
	AdaptedCV     string   `json:"adapted_cv"`
	Changes       []string `json:"changes"`
	KeywordsAdded []string `json:"keywords_added"`
}

// ParseAdaptResult parsea la respuesta JSON del LLM para adaptación de CV.
func ParseAdaptResult(raw string) (*domain.AdaptResult, error) {
	cleaned := extractJSON(raw)
	var j adaptJSON
	if err := json.Unmarshal([]byte(cleaned), &j); err != nil {
		return nil, fmt.Errorf("ParseAdaptResult: %w\nrespuesta: %.300s", err, raw)
	}
	return &domain.AdaptResult{
		AdaptedCV:     j.AdaptedCV,
		Changes:       j.Changes,
		KeywordsAdded: j.KeywordsAdded,
	}, nil
}

// extractJSON extrae el primer objeto JSON completo de un string.
// Tolera bloques de código markdown (```json ... ```).
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	// Quitar bloque de código markdown si existe
	if strings.HasPrefix(s, "```") {
		end := strings.LastIndex(s, "```")
		if end > 3 {
			s = s[3:end]
		}
		s = strings.TrimPrefix(s, "json")
		s = strings.TrimSpace(s)
	}
	// Extraer primer { ... }
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
