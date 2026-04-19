// core/internal/tui/portal_selector.go
package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// PortalChoice representa una opción en el menú de selección de bolsa.
type PortalChoice struct {
	ID        string // "occ", "all", "computrabajo", "indeed"
	Label     string // texto mostrado al usuario
	Available bool   // false = visible pero no seleccionable
	IsSep     bool   // true = separador visual, no opción
}

// portalOptions define el menú en orden de aparición.
// Para habilitar una bolsa cuando esté lista: cambiar Available a true.
var portalOptions = []PortalChoice{
	{ID: "occ", Label: "🔍 OCC Mundial         (disponible)", Available: true},
	{ID: "all", Label: "📋 Todas las bolsas    (usa OCC por ahora)", Available: true},
	{ID: "", Label: "─────────────────────────────────────────────", IsSep: true},
	{ID: "computrabajo", Label: "⏳ Computrabajo         (próximamente)", Available: false},
	{ID: "indeed", Label: "⏳ Indeed México        (próximamente)", Available: false},
}

// PortalSelectorModel es el modelo Bubble Tea del selector de bolsa.
type PortalSelectorModel struct {
	cursor   int
	selected string
	done     bool
	quitting bool
}

// InitialPortalSelector devuelve el modelo con cursor en la primera opción disponible.
func InitialPortalSelector() PortalSelectorModel {
	return PortalSelectorModel{cursor: 0}
}

func (m PortalSelectorModel) Init() tea.Cmd {
	return nil
}

func (m PortalSelectorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "up", "k":
			m.cursor = prevAvailable(m.cursor)

		case "down", "j", "tab":
			m.cursor = nextAvailable(m.cursor)

		case "enter", " ":
			opt := portalOptions[m.cursor]
			if opt.Available && !opt.IsSep {
				m.selected = opt.ID
				m.done = true
				return m, tea.Quit
			}
			// Opción deshabilitada o separador → no hacer nada
		}
	}
	return m, nil
}

func (m PortalSelectorModel) View() string {
	if m.quitting {
		return ""
	}

	header := "\n╔══════════════════════════════════════════════════╗\n" +
		"║       JobSense AI — Búsqueda de Vacantes         ║\n" +
		"╚══════════════════════════════════════════════════╝\n\n"

	body := "  ¿En qué bolsa de trabajo quieres buscar?\n\n"

	for i, opt := range portalOptions {
		if opt.IsSep {
			body += "    " + opt.Label + "\n"
			continue
		}
		cursor := "  "
		if i == m.cursor {
			cursor = "❯ "
		}
		// Dim ANSI para opciones no disponibles
		prefix, suffix := "", ""
		if !opt.Available {
			prefix = "\033[2m"
			suffix = "\033[0m"
		}
		body += "  " + cursor + prefix + opt.Label + suffix + "\n"
	}

	body += "\n  ↑↓ para mover · Enter para seleccionar · Ctrl+C para salir\n"
	return header + body
}

// Selected devuelve el portal elegido. Vacío si el usuario salió sin elegir.
func (m PortalSelectorModel) Selected() string { return m.selected }

// Quitting devuelve true si el usuario presionó Ctrl+C sin seleccionar.
func (m PortalSelectorModel) Quitting() bool { return m.quitting }

// nextAvailable avanza al siguiente portal habilitado (con wrap).
func nextAvailable(current int) int {
	for i := 1; i <= len(portalOptions); i++ {
		next := (current + i) % len(portalOptions)
		opt := portalOptions[next]
		if opt.Available && !opt.IsSep {
			return next
		}
	}
	return current
}

// prevAvailable retrocede al portal habilitado anterior (con wrap).
func prevAvailable(current int) int {
	for i := 1; i <= len(portalOptions); i++ {
		prev := (current - i + len(portalOptions)) % len(portalOptions)
		opt := portalOptions[prev]
		if opt.Available && !opt.IsSep {
			return prev
		}
	}
	return current
}
