// core/internal/cli/analyze.go
package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/QUERTY/OfferTrack-M82/internal/ai"
	"github.com/QUERTY/OfferTrack-M82/internal/config"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/services"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze <job-id>",
	Short: "Analizar una vacante contra tu perfil con IA",
	Args:  cobra.ExactArgs(1),
	RunE:  runAnalyze,
}

func init() {
	analyzeCmd.Flags().String("profile", "", "Ruta al perfil del candidato (Markdown)")
	analyzeCmd.Flags().String("cv", "", "Ruta al CV base (Markdown)")
}

func runAnalyze(cmd *cobra.Command, args []string) error {
	jobID := args[0]
	profilePath, _ := cmd.Flags().GetString("profile")
	cvPath, _ := cmd.Flags().GetString("cv")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	provider, err := ai.NewProvider(cfg)
	if err != nil {
		return fmt.Errorf("proveedor IA (%s): %w", cfg.AI.Provider, err)
	}

	qdrantClient, err := db.NewQdrantClient(
		cfg.Qdrant.Host, cfg.Qdrant.Port, cfg.Qdrant.Collections["jobs"],
	)
	if err != nil {
		return fmt.Errorf("qdrant: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	job, err := qdrantClient.GetJob(ctx, jobID)
	if err != nil {
		return fmt.Errorf("vacante %q no encontrada: %w\n(usa 'offertrack list' para ver los IDs)", jobID, err)
	}

	profile := readFileCLI(profilePath, "")
	cv := readFileCLI(cvPath, "")

	fmt.Printf("\nAnalizando: %s @ %s\n", job.Title, job.Company)
	fmt.Printf("Proveedor:  %s / %s\n\n", cfg.AI.Provider, cfg.AI.Model)

	analyzer := services.NewAnalyzerService(provider, qdrantClient)
	analysis, err := analyzer.AnalyzeJob(ctx, job, profile, cv)
	if err != nil {
		return fmt.Errorf("análisis fallido: %w", err)
	}

	printAnalysisResult(analysis)
	return nil
}

func printAnalysisResult(a *domain.Analysis) {
	rec := strings.ToUpper(a.Recommendation)
	fmt.Printf("══════════════════════════════════════════\n")
	fmt.Printf(" ANÁLISIS IA\n")
	fmt.Printf("══════════════════════════════════════════\n")
	fmt.Printf(" Compatibilidad : %d / 100\n", a.CompatibilityScore)
	fmt.Printf(" Recomendación  : %s\n", rec)
	fmt.Printf("──────────────────────────────────────────\n")

	if len(a.Strengths) > 0 {
		fmt.Printf("\n Fortalezas:\n")
		for _, s := range a.Strengths {
			fmt.Printf("   • %s\n", s)
		}
	}
	if len(a.Gaps) > 0 {
		fmt.Printf("\n Brechas:\n")
		for _, g := range a.Gaps {
			fmt.Printf("   • %s\n", g)
		}
	}
	if a.SalaryEstimate != "" {
		fmt.Printf("\n Estimado salarial: %s\n", a.SalaryEstimate)
	}
	if a.RawAnalysis != "" {
		fmt.Printf("\n Análisis completo:\n %s\n", a.RawAnalysis)
	}
	fmt.Printf("══════════════════════════════════════════\n\n")
}

// readFileCLI lee un archivo y retorna su contenido, o fallback si no se puede.
func readFileCLI(path, fallback string) string {
	if path == "" {
		return fallback
	}
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[warn] no se pudo leer %q: %v\n", path, err)
		return fallback
	}
	return string(data)
}
