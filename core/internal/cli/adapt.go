// core/internal/cli/adapt.go
package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/QUERTY/OfferTrack-M82/internal/ai"
	"github.com/QUERTY/OfferTrack-M82/internal/config"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/services"
)

var adaptCmd = &cobra.Command{
	Use:   "adapt-cv <job-id>",
	Short: "Generar CV adaptado para una vacante con IA",
	Args:  cobra.ExactArgs(1),
	RunE:  runAdaptCV,
}

func init() {
	adaptCmd.Flags().String("cv", "", "Ruta al CV base en Markdown (requerido)")
	adaptCmd.Flags().String("profile", "", "Ruta al perfil del candidato en Markdown")
	adaptCmd.Flags().String("output", "./exports", "Directorio de salida")
	adaptCmd.Flags().String("format", "md", "Formato de salida: md, pdf, docx")
	_ = adaptCmd.MarkFlagRequired("cv")
}

func runAdaptCV(cmd *cobra.Command, args []string) error {
	jobID := args[0]
	cvPath, _ := cmd.Flags().GetString("cv")
	profilePath, _ := cmd.Flags().GetString("profile")
	outputDir, _ := cmd.Flags().GetString("output")
	format, _ := cmd.Flags().GetString("format")

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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	job, err := qdrantClient.GetJob(ctx, jobID)
	if err != nil {
		return fmt.Errorf("vacante %q no encontrada: %w\n(usa 'offertrack list' para ver los IDs)", jobID, err)
	}

	cv := readFileCLI(cvPath, "")
	if cv == "" {
		return fmt.Errorf("no se pudo leer el CV desde %q", cvPath)
	}
	profile := readFileCLI(profilePath, "")

	fmt.Printf("\nAdaptando CV para: %s @ %s\n", job.Title, job.Company)
	fmt.Printf("Proveedor: %s / %s\n\n", cfg.AI.Provider, cfg.AI.Model)

	adapter := services.NewCVAdapterService(provider)
	result, err := adapter.Adapt(ctx, job, cv, profile)
	if err != nil {
		return fmt.Errorf("adaptación fallida: %w", err)
	}

	exporter, err := services.NewExporterService(outputDir, "")
	if err != nil {
		return fmt.Errorf("exporter: %w", err)
	}

	outPath, err := exporter.ExportCV(result.AdaptedCV, job.Title, format)
	if err != nil {
		return fmt.Errorf("export fallido: %w", err)
	}

	fmt.Printf("✓ CV guardado en: %s\n\n", outPath)

	if len(result.Changes) > 0 {
		fmt.Printf("Cambios realizados:\n")
		for _, c := range result.Changes {
			fmt.Printf("  • %s\n", c)
		}
		fmt.Println()
	}

	if len(result.KeywordsAdded) > 0 {
		fmt.Printf("Keywords añadidas: %s\n", strings.Join(result.KeywordsAdded, ", "))
	}

	return nil
}
