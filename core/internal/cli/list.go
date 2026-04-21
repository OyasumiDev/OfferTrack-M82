// core/internal/cli/list.go
package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/QUERTY/OfferTrack-M82/internal/config"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/services"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Listar vacantes guardadas en Qdrant",
	RunE:  runList,
}

func init() {
	listCmd.Flags().Int("limit", 50, "Máximo de vacantes a mostrar")
	listCmd.Flags().String("portal", "", "Filtrar por portal: occ, computrabajo, indeed")
	listCmd.Flags().String("modality", "", "Filtrar por modalidad: remote, hybrid, onsite")
	listCmd.Flags().Int("salary-min", 0, "Filtrar por salario mínimo (MXN)")
}

func runList(cmd *cobra.Command, _ []string) error {
	limit, _ := cmd.Flags().GetInt("limit")
	portal, _ := cmd.Flags().GetString("portal")
	modality, _ := cmd.Flags().GetString("modality")
	salaryMin, _ := cmd.Flags().GetInt("salary-min")

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	qdrantClient, err := db.NewQdrantClient(
		cfg.Qdrant.Host, cfg.Qdrant.Port, cfg.Qdrant.Collections["jobs"],
	)
	if err != nil {
		return fmt.Errorf("qdrant: %w\n(¿está corriendo Docker con Qdrant?)", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := qdrantClient.Ping(ctx); err != nil {
		return fmt.Errorf("Qdrant no responde: %w", err)
	}

	// Con filtros: usar FilterJobs; sin filtros: ListJobs vía Orchestrator
	var jobs interface{ error }
	if portal != "" || modality != "" || salaryMin > 0 {
		filtered, e := qdrantClient.FilterJobs(ctx, db.JobFilter{
			Portal:    portal,
			Modality:  modality,
			SalaryMin: salaryMin,
		}, uint64(limit))
		if e != nil {
			return fmt.Errorf("error filtrando vacantes: %w", e)
		}
		if len(filtered) == 0 {
			fmt.Println("No se encontraron vacantes con esos filtros.")
			return nil
		}
		printJobsTable(filtered)
		fmt.Printf("\n%d vacantes mostradas.\n", len(filtered))
		return nil
	}

	_ = jobs
	orch := services.NewOrchestrator(qdrantClient, cfg.Scraper.BaseURL)
	saved, err := orch.ListSaved(ctx, uint64(limit))
	if err != nil {
		return fmt.Errorf("error listando vacantes: %w", err)
	}

	if len(saved) == 0 {
		fmt.Println("No hay vacantes guardadas.")
		fmt.Println("Ejecuta: offertrack search --role=\"<rol>\"")
		return nil
	}

	printJobsTable(saved)
	fmt.Printf("\n%d vacantes guardadas en Qdrant.\n", len(saved))
	return nil
}
