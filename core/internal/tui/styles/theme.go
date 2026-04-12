package styles

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF"))
	SubtitleStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	HighlightStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF88"))
	DimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))
)
