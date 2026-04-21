// core/internal/tui/theme.go
package tui

import "github.com/charmbracelet/lipgloss"

// Paleta de colores.
var (
	colorPrimary  = lipgloss.Color("#7C3AED")
	colorAccent   = lipgloss.Color("#06B6D4")
	colorSuccess  = lipgloss.Color("#10B981")
	colorWarning  = lipgloss.Color("#F59E0B")
	colorDanger   = lipgloss.Color("#EF4444")
	colorSubtle   = lipgloss.Color("#9CA3AF")
	colorSelected = lipgloss.Color("#3B0764")
	colorText     = lipgloss.Color("#F3F4F6")
	colorDim      = lipgloss.Color("#6B7280")
)

var (
	styleTitle = lipgloss.NewStyle().
			Foreground(colorPrimary).Bold(true)

	styleSubtitle = lipgloss.NewStyle().
			Foreground(colorAccent)

	styleColHeader = lipgloss.NewStyle().
			Foreground(colorText).Bold(true).Underline(true)

	styleSelectedRow = lipgloss.NewStyle().
				Background(colorSelected).Foreground(colorText).Bold(true)

	styleMuted = lipgloss.NewStyle().Foreground(colorDim)

	styleSuccess = lipgloss.NewStyle().Foreground(colorSuccess).Bold(true)
	styleWarning = lipgloss.NewStyle().Foreground(colorWarning).Bold(true)
	styleDanger  = lipgloss.NewStyle().Foreground(colorDanger).Bold(true)

	styleHelp    = lipgloss.NewStyle().Foreground(colorSubtle)
	styleDivider = lipgloss.NewStyle().Foreground(colorDim)
	styleEmphasis = lipgloss.NewStyle().Foreground(colorAccent).Bold(true)
	styleStatus  = lipgloss.NewStyle().Foreground(colorAccent)
)

// scoreStyle devuelve el estilo adecuado según el score de compatibilidad.
func scoreStyle(score int) lipgloss.Style {
	switch {
	case score >= 75:
		return styleSuccess
	case score >= 50:
		return styleWarning
	default:
		return styleDanger
	}
}

// recommendStyle devuelve el estilo según la recomendación del LLM.
func recommendStyle(rec string) lipgloss.Style {
	switch rec {
	case "apply":
		return styleSuccess
	case "consider":
		return styleWarning
	default:
		return styleDanger
	}
}
