// core/internal/services/exporter.go
package services

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/QUERTY/OfferTrack-M82/internal/domain"
)

const defaultPandocPath = `C:\Program Files\Pandoc\pandoc.exe`

// ExporterService convierte análisis y CVs a archivos Markdown / PDF / DOCX.
type ExporterService struct {
	pandocPath string
	outputDir  string
}

// NewExporterService crea un ExporterService.
// outputDir es el directorio donde se guardan los archivos generados.
// pandocPath es la ruta al ejecutable pandoc ("" usa el valor por defecto).
func NewExporterService(outputDir, pandocPath string) (*ExporterService, error) {
	if outputDir == "" {
		outputDir = "./output"
	}
	if pandocPath == "" {
		pandocPath = defaultPandocPath
	}
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, fmt.Errorf("exporter: no se pudo crear %q → %w", outputDir, err)
	}
	return &ExporterService{pandocPath: pandocPath, outputDir: outputDir}, nil
}

// ExportAnalysis escribe el análisis como Markdown y, si format != "md",
// invoca Pandoc para convertirlo al formato solicitado ("pdf" o "docx").
// Retorna la ruta del archivo generado.
func (e *ExporterService) ExportAnalysis(analysis *domain.Analysis, job *domain.Job, format string) (string, error) {
	md := buildAnalysisMD(analysis, job)

	slug := sanitizeSlug(job.Title)
	date := time.Now().Format("2006-01-02")
	mdPath := filepath.Join(e.outputDir, fmt.Sprintf("%s_%s_analysis.md", date, slug))

	if err := os.WriteFile(mdPath, []byte(md), 0o644); err != nil {
		return "", fmt.Errorf("exporter: escribir MD → %w", err)
	}

	if format == "md" || format == "" {
		return mdPath, nil
	}
	return e.convertWithPandoc(mdPath, format)
}

// ExportCV escribe el CV adaptado como Markdown y, si format != "md",
// invoca Pandoc para convertirlo. Retorna la ruta del archivo generado.
func (e *ExporterService) ExportCV(adaptedCV, jobTitle, format string) (string, error) {
	slug := sanitizeSlug(jobTitle)
	date := time.Now().Format("2006-01-02")
	mdPath := filepath.Join(e.outputDir, fmt.Sprintf("%s_%s_cv.md", date, slug))

	if err := os.WriteFile(mdPath, []byte(adaptedCV), 0o644); err != nil {
		return "", fmt.Errorf("exporter: escribir CV MD → %w", err)
	}

	if format == "md" || format == "" {
		return mdPath, nil
	}
	return e.convertWithPandoc(mdPath, format)
}

// convertWithPandoc convierte un archivo .md al formato pedido usando Pandoc.
func (e *ExporterService) convertWithPandoc(mdPath, format string) (string, error) {
	ext := strings.ToLower(format)
	if ext != "pdf" && ext != "docx" {
		return "", fmt.Errorf("exporter: formato no soportado %q (use md, pdf o docx)", format)
	}
	outPath := strings.TrimSuffix(mdPath, ".md") + "." + ext

	// #nosec G204 — mdPath y outPath son rutas internas generadas por el propio servicio
	cmd := exec.Command(e.pandocPath, mdPath, "-o", outPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if _, err := os.Stat(e.pandocPath); os.IsNotExist(err) {
		return "", fmt.Errorf("exporter: Pandoc no encontrado en %q\nInstálalo desde https://pandoc.org/installing.html", e.pandocPath)
	}

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exporter: pandoc falló → %w", err)
	}
	return outPath, nil
}

// buildAnalysisMD construye el contenido Markdown del análisis.
func buildAnalysisMD(a *domain.Analysis, job *domain.Job) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Análisis: %s — %s\n\n", job.Title, job.Company))
	sb.WriteString(fmt.Sprintf("**Portal:** %s | **Modalidad:** %s | **Ubicación:** %s\n\n", job.Portal, job.Modality, job.Location))
	if job.SalaryMin > 0 {
		sb.WriteString(fmt.Sprintf("**Salario:** %d – %d %s\n\n", job.SalaryMin, job.SalaryMax, job.Currency))
	}
	sb.WriteString(fmt.Sprintf("**URL:** %s\n\n", job.URL))
	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("## Compatibilidad: %d/100  (%s)\n\n", a.CompatibilityScore, strings.ToUpper(a.Recommendation)))
	sb.WriteString("### Fortalezas\n\n")
	for _, s := range a.Strengths {
		sb.WriteString(fmt.Sprintf("- %s\n", s))
	}
	sb.WriteString("\n### Brechas\n\n")
	for _, g := range a.Gaps {
		sb.WriteString(fmt.Sprintf("- %s\n", g))
	}
	sb.WriteString(fmt.Sprintf("\n**Estimado salarial:** %s\n\n", a.SalaryEstimate))
	sb.WriteString("### Análisis completo\n\n")
	sb.WriteString(a.RawAnalysis + "\n")
	return sb.String()
}

// sanitizeSlug convierte un título en un slug seguro para nombres de archivo.
func sanitizeSlug(title string) string {
	r := strings.NewReplacer(" ", "_", "/", "-", "\\", "-", ":", "", "\"", "", "'", "")
	slug := r.Replace(strings.ToLower(title))
	if len(slug) > 40 {
		slug = slug[:40]
	}
	return slug
}
