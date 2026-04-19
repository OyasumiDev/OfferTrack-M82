// core/internal/cli/tui.go
package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/QUERTY/OfferTrack-M82/internal/ai"
	"github.com/QUERTY/OfferTrack-M82/internal/config"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/services"
	"github.com/QUERTY/OfferTrack-M82/internal/tui"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interfaz de usuario en terminal (TUI interactiva)",
	RunE:  runTUI,
}

func init() {
	tuiCmd.Flags().String("profile", "", "Ruta al perfil del candidato en Markdown")
	tuiCmd.Flags().String("cv", "", "Ruta al CV base en Markdown")
}

func runTUI(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	qdrantClient, err := db.NewQdrantClient(
		cfg.Qdrant.Host,
		cfg.Qdrant.Port,
		cfg.Qdrant.Collections["jobs"],
	)
	if err != nil {
		return fmt.Errorf("qdrant: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := qdrantClient.Ping(ctx); err != nil {
		return fmt.Errorf("Qdrant no responde en %s:%d → %w\n(¿está corriendo Docker?)",
			cfg.Qdrant.Host, cfg.Qdrant.Port, err)
	}

	// Proveedor IA (opcional — la TUI funciona sin él, solo sin análisis)
	provider, provErr := ai.NewProvider(cfg)
	if provErr != nil {
		fmt.Fprintf(os.Stderr, "[warn] proveedor IA no disponible: %v\n", provErr)
	}

	var analyzerSvc *services.AnalyzerService
	if provider != nil {
		analyzerSvc = services.NewAnalyzerService(provider, qdrantClient)
	}

	profilePath, _ := cmd.Flags().GetString("profile")
	cvPath, _ := cmd.Flags().GetString("cv")

	m := tui.NewModel(qdrantClient, analyzerSvc, cfg, profilePath, cvPath)
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}
