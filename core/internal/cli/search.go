// core/internal/cli/search.go
package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/QUERTY/OfferTrack-M82/internal/config"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/services"
	"github.com/QUERTY/OfferTrack-M82/internal/tui"
)

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Buscar vacantes en portales de empleo",
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().String("role", "", "Puesto o palabras clave (se pide interactivamente si se omite)")
	searchCmd.Flags().String("location", "", "Ciudad o estado (se pide interactivamente si se omite)")
	searchCmd.Flags().Int("salary-min", 0, "Salario mínimo mensual MXN")
	searchCmd.Flags().Int("salary-max", 0, "Salario máximo mensual MXN")
	searchCmd.Flags().String("modality", "", "Modalidad: presencial, hibrido, remoto")
	searchCmd.Flags().String("portal", "", "Salta el selector: occ, all")
	searchCmd.Flags().StringSlice("portals", []string{}, "Legacy — usa --portal")
	searchCmd.Flags().Int("max", 20, "Máximo de vacantes a mostrar")
}

func runSearch(cmd *cobra.Command, _ []string) error {

	// ── PASO 1: elegir portal ─────────────────────────────────────────────────

	portalFlag, _ := cmd.Flags().GetString("portal")

	// Compatibilidad con --portals (legacy, usado por start.ps1)
	if portalFlag == "" {
		legacyPortals, _ := cmd.Flags().GetStringSlice("portals")
		if len(legacyPortals) > 0 && legacyPortals[0] != "" {
			portalFlag = legacyPortals[0]
		}
	}

	selectedPortal := portalFlag
	if selectedPortal == "" {
		// Sin flag → mostrar selector interactivo de Bubble Tea
		result, err := runPortalSelector()
		if err != nil {
			return fmt.Errorf("error en el selector de portal: %w", err)
		}
		if result.Quitting() || result.Selected() == "" {
			fmt.Println("Búsqueda cancelada.")
			return nil
		}
		selectedPortal = result.Selected()
	}

	// Validar que el portal sea reconocido y esté disponible
	validPortals := map[string]bool{"occ": true, "all": true, "indeed": true}
	if !validPortals[selectedPortal] {
		return fmt.Errorf(
			"portal %q no reconocido o aún no disponible.\nPortales disponibles: occ, all",
			selectedPortal,
		)
	}

	// ── PASO 2: resolver portales activos ────────────────────────────────────

	var activePortals []string
	switch selectedPortal {
	case "all":
		activePortals = []string{"occ", "indeed"}
	default:
		activePortals = []string{selectedPortal}
	}

	// ── PASO 3: pedir parámetros de búsqueda ─────────────────────────────────

	role, _ := cmd.Flags().GetString("role")
	location, _ := cmd.Flags().GetString("location")
	salaryMin, _ := cmd.Flags().GetInt("salary-min")
	salaryMax, _ := cmd.Flags().GetInt("salary-max")
	modality, _ := cmd.Flags().GetString("modality")
	maxResults, _ := cmd.Flags().GetInt("max")

	if role == "" {
		role = promptLine("  Puesto (ej: desarrollador .NET): ")
		if strings.TrimSpace(role) == "" {
			return fmt.Errorf("el puesto es requerido")
		}
	}
	if location == "" {
		location = promptLine("  Ciudad (ej: monterrey — Enter para omitir): ")
	}

	role = strings.TrimSpace(role)
	location = strings.TrimSpace(location)

	// ── PASO 4: ejecutar búsqueda ─────────────────────────────────────────────

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("error cargando configuración: %w", err)
	}

	qdrantClient, err := db.NewQdrantClient(cfg.Qdrant.Host, cfg.Qdrant.Port, cfg.Qdrant.Collections["jobs"])
	if err != nil {
		return fmt.Errorf("error conectando Qdrant: %w\n(¿está corriendo Docker con Qdrant?)", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	if err := qdrantClient.Ping(ctx); err != nil {
		return fmt.Errorf("Qdrant no responde: %w", err)
	}

	if err := db.InitCollections(ctx, qdrantClient.RawClient(), db.CollectionNames{
		Jobs:    cfg.Qdrant.Collections["jobs"],
		Profile: cfg.Qdrant.Collections["profile"],
		CVs:     cfg.Qdrant.Collections["cvs"],
		Memory:  cfg.Qdrant.Collections["memory"],
	}); err != nil {
		return fmt.Errorf("error inicializando colecciones: %w", err)
	}

	orch := services.NewOrchestrator(qdrantClient, cfg.Scraper.BaseURL)

	fmt.Printf("\nBuscando \"%s\" en %v...\n\n", role, activePortals)

	jobs, err := orch.Search(ctx, services.SearchParams{
		Role:       role,
		Location:   location,
		SalaryMin:  salaryMin,
		SalaryMax:  salaryMax,
		Modality:   modality,
		Portals:    activePortals,
		MaxResults: maxResults,
		MaxPages:   3,
	})
	if err != nil {
		return fmt.Errorf("error en búsqueda: %w", err)
	}

	if len(jobs) == 0 {
		fmt.Println("No se encontraron vacantes. Intenta con otro rol o ubicación.")
		return nil
	}

	printJobsTable(jobs)
	fmt.Printf("\n%d vacantes guardadas en Qdrant.\n", len(jobs))
	return nil
}

// runPortalSelector ejecuta el selector Bubble Tea y devuelve el modelo final.
func runPortalSelector() (tui.PortalSelectorModel, error) {
	m := tui.InitialPortalSelector()
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return tui.PortalSelectorModel{}, err
	}
	return final.(tui.PortalSelectorModel), nil
}

// promptLine muestra un label y lee una línea de stdin.
func promptLine(label string) string {
	fmt.Print(label)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func printJobsTable(jobs []*domain.Job) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "#\tTÍTULO\tEMPRESA\tSALARIO\tMODALIDAD\tPORTAL")
	fmt.Fprintln(w, "─\t──────\t───────\t───────\t─────────\t──────")
	for i, job := range jobs {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			i+1, truncate(job.Title, 35), truncate(job.Company, 20),
			formatSalary(job), job.Modality, job.Portal)
	}
	w.Flush()
}

func formatSalary(job *domain.Job) string {
	if job.SalaryMin == 0 && job.SalaryMax == 0 {
		return "N/D"
	}
	if job.SalaryMin == job.SalaryMax {
		return fmt.Sprintf("$%d", job.SalaryMin)
	}
	return fmt.Sprintf("$%d–%d", job.SalaryMin, job.SalaryMax)
}

func truncate(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	return string([]rune(s)[:max-1]) + "…"
}
