package ai

import "fmt"

func BuildAnalysisPrompt(job, profile, cv string) string {
	return fmt.Sprintf(`Eres un experto en reclutamiento. Analiza esta vacante contra el perfil del candidato.

VACANTE:
%s

PERFIL DEL CANDIDATO:
%s

CV BASE:
%s

Responde SOLO en JSON con este esquema:
{
  "compatibility_score": <0-100>,
  "strengths": ["..."],
  "gaps": ["..."],
  "salary_estimate": "...",
  "recommendation": "apply|consider|discard",
  "raw_analysis": "..."
}`, job, profile, cv)
}

func BuildAdaptPrompt(job, cv, profile string) string {
	return fmt.Sprintf(`Eres un experto en redaccion de CVs. Adapta el CV a la vacante.

REGLAS:
- No inventes experiencias ni habilidades inexistentes
- Toda informacion debe estar respaldada en el CV base
- Preserva la linea de tiempo laboral
- Prioriza palabras clave de la vacante

VACANTE: %s
CV BASE: %s
PERFIL: %s

Responde SOLO en JSON:
{"adapted_cv":"<Markdown>","changes":["..."],"keywords_added":["..."]}`, job, cv, profile)
}
