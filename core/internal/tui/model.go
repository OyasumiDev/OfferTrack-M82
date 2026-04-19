// core/internal/tui/model.go
package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/QUERTY/OfferTrack-M82/internal/config"
	"github.com/QUERTY/OfferTrack-M82/internal/db"
	"github.com/QUERTY/OfferTrack-M82/internal/domain"
	"github.com/QUERTY/OfferTrack-M82/internal/services"
)

// ── Estados de la aplicación ──────────────────────────────────────────────────

type appState int

const (
	stateLoading   appState = iota
	stateJobList
	stateJobDetail
	stateAnalyzing
	stateError
)

var spinFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// ── Mensajes Bubble Tea ───────────────────────────────────────────────────────

type (
	jobsLoadedMsg  struct{ jobs []*domain.Job }
	analysisOkMsg  struct {
		jobID  string
		result *domain.Analysis
	}
	analysisErrMsg struct{ err error }
	fatalErrMsg    struct{ err error }
	tickMsg        time.Time
)

// ── Modelo principal ──────────────────────────────────────────────────────────

// Model es el modelo Bubble Tea de OfferTrack.
type Model struct {
	// dependencias
	qdrant      *db.QdrantClient
	analyzer    *services.AnalyzerService
	cfg         *config.Config
	profilePath string
	cvPath      string

	// estado UI
	state     appState
	jobs      []*domain.Job
	cursor    int
	offset    int
	analyses  map[string]*domain.Analysis
	viewJob   *domain.Job
	errText   string
	statusMsg string
	spinFrame int
	width     int
	height    int
}

// NewModel crea el modelo TUI con las dependencias requeridas.
func NewModel(
	qdrant *db.QdrantClient,
	analyzer *services.AnalyzerService,
	cfg *config.Config,
	profilePath, cvPath string,
) Model {
	return Model{
		qdrant:      qdrant,
		analyzer:    analyzer,
		cfg:         cfg,
		profilePath: profilePath,
		cvPath:      cvPath,
		state:       stateLoading,
		analyses:    make(map[string]*domain.Analysis),
	}
}

// ── Interfaz tea.Model ────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return tea.Batch(cmdLoadJobs(m.qdrant), cmdTick())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tickMsg:
		m.spinFrame = (m.spinFrame + 1) % len(spinFrames)
		if m.state == stateLoading || m.state == stateAnalyzing {
			return m, cmdTick()
		}

	case jobsLoadedMsg:
		m.jobs = msg.jobs
		m.state = stateJobList
		m.statusMsg = fmt.Sprintf("%d vacantes cargadas", len(m.jobs))

	case fatalErrMsg:
		m.errText = msg.err.Error()
		m.state = stateError

	case analysisOkMsg:
		m.analyses[msg.jobID] = msg.result
		if m.viewJob != nil && m.viewJob.ID == msg.jobID {
			m.state = stateJobDetail
		} else {
			m.state = stateJobList
		}
		m.statusMsg = fmt.Sprintf("✓ Análisis: %d/100 — %s",
			msg.result.CompatibilityScore, strings.ToUpper(msg.result.Recommendation))

	case analysisErrMsg:
		if m.viewJob != nil {
			m.state = stateJobDetail
		} else {
			m.state = stateJobList
		}
		m.statusMsg = "⚠ " + msg.err.Error()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "ctrl+c" {
		return m, tea.Quit
	}
	switch m.state {
	case stateJobList:
		return m.keyJobList(msg)
	case stateJobDetail:
		return m.keyJobDetail(msg)
	case stateError:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) keyJobList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := m.visibleRows()
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			if m.cursor < m.offset {
				m.offset = m.cursor
			}
		}
	case "down", "j":
		if m.cursor < len(m.jobs)-1 {
			m.cursor++
			if m.cursor >= m.offset+visible {
				m.offset = m.cursor - visible + 1
			}
		}
	case "enter":
		if len(m.jobs) > 0 {
			m.viewJob = m.jobs[m.cursor]
			m.state = stateJobDetail
			m.statusMsg = ""
		}
	case "a":
		if len(m.jobs) == 0 {
			break
		}
		if m.analyzer == nil {
			m.statusMsg = "⚠ Sin proveedor IA — configura tu API key en .env"
			break
		}
		job := m.jobs[m.cursor]
		m.state = stateAnalyzing
		m.statusMsg = "Analizando: " + truncateTUI(job.Title, 40) + "…"
		return m, tea.Batch(cmdAnalyze(m.analyzer, job, m.profilePath, m.cvPath), cmdTick())
	case "r":
		m.state = stateLoading
		m.statusMsg = ""
		return m, tea.Batch(cmdLoadJobs(m.qdrant), cmdTick())
	}
	return m, nil
}

func (m Model) keyJobDetail(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		return m, tea.Quit
	case "esc", "backspace", "b":
		m.state = stateJobList
		m.viewJob = nil
		m.statusMsg = ""
	case "a":
		if m.viewJob == nil || m.analyzer == nil {
			m.statusMsg = "⚠ Sin proveedor IA — configura tu API key en .env"
			break
		}
		m.state = stateAnalyzing
		m.statusMsg = "Analizando: " + truncateTUI(m.viewJob.Title, 40) + "…"
		return m, tea.Batch(cmdAnalyze(m.analyzer, m.viewJob, m.profilePath, m.cvPath), cmdTick())
	}
	return m, nil
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.width == 0 {
		return "Iniciando…\n"
	}

	header := m.renderHeader()
	div := styleDivider.Render(strings.Repeat("─", m.width))

	var body, footerText string
	switch m.state {
	case stateLoading:
		body = m.bodySpinner("Cargando vacantes desde Qdrant")
		footerText = "ctrl+c salir"
	case stateAnalyzing:
		body = m.bodySpinner(m.statusMsg)
		footerText = "ctrl+c salir"
	case stateError:
		body = m.bodyError()
		footerText = "q salir"
	case stateJobDetail:
		body, footerText = m.bodyJobDetail()
	default:
		body, footerText = m.bodyJobList()
	}

	// En la lista de vacantes el statusMsg reemplaza el footer
	footer := styleHelp.Render("  " + footerText)
	if m.statusMsg != "" && m.state == stateJobList {
		footer = styleStatus.Render("  " + m.statusMsg)
	}

	// Rellenar el cuerpo para que el footer quede siempre abajo
	contentH := m.height - 4 // header(1) + div(1) + div(1) + footer(1)
	if contentH < 1 {
		contentH = 1
	}
	body = padToHeight(body, contentH)

	return lipgloss.JoinVertical(lipgloss.Left, header, div, body, div, footer)
}

// ── Secciones de la pantalla ──────────────────────────────────────────────────

func (m Model) renderHeader() string {
	left := styleTitle.Render("  OfferTrack M82  ")
	right := styleMuted.Render("  " + m.cfg.AI.Provider + " / " + m.cfg.AI.Model + "  ")
	mid := styleSubtitle.Render(fmt.Sprintf("%d vacantes", len(m.jobs)))

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right) - lipgloss.Width(mid)
	if gap < 2 {
		gap = 2
	}
	lPad := gap / 2
	rPad := gap - lPad
	return left + strings.Repeat(" ", lPad) + mid + strings.Repeat(" ", rPad) + right
}

func (m Model) bodySpinner(msg string) string {
	spin := styleEmphasis.Render(spinFrames[m.spinFrame])
	return "\n\n  " + spin + "  " + styleMuted.Render(msg)
}

func (m Model) bodyError() string {
	return "\n\n  " + styleDanger.Render("✖ Error fatal:") +
		"\n\n  " + styleMuted.Render(m.errText) + "\n"
}

// bodyJobList renderiza la tabla de vacantes.
func (m Model) bodyJobList() (body, help string) {
	if len(m.jobs) == 0 {
		body = "\n\n  " + styleMuted.Render(`Sin vacantes. Ejecuta:  offertrack search --role="<rol>"`)
		help = "r recargar  q salir"
		return
	}

	const (
		wIdx    = 4
		wTitle  = 32
		wCo     = 22
		wSalary = 13
		wScore  = 6
		wModal  = 8
	)
	colFmt := fmt.Sprintf("%%-%ds  %%-%ds  %%-%ds  %%-%ds  %%-%ds  %%-%ds", wIdx, wTitle, wCo, wSalary, wScore, wModal)

	hdr := styleColHeader.Render("  " + fmt.Sprintf(colFmt, "#", "TÍTULO", "EMPRESA", "SALARIO", "SCORE", "MODAL"))

	visible := m.visibleRows()
	end := m.offset + visible
	if end > len(m.jobs) {
		end = len(m.jobs)
	}

	rows := make([]string, 0, end-m.offset)
	for i := m.offset; i < end; i++ {
		job := m.jobs[i]
		row := fmt.Sprintf(colFmt,
			fmt.Sprintf("%d", i+1),
			truncateTUI(job.Title, wTitle),
			truncateTUI(job.Company, wCo),
			fmtSalary(job),
			renderScore(job.CompatScore),
			truncateTUI(job.Modality, wModal),
		)
		if i == m.cursor {
			rows = append(rows, styleSelectedRow.Render("▶ "+row))
		} else {
			rows = append(rows, "  "+row)
		}
	}

	// Indicador de scroll
	var hints []string
	if m.offset > 0 {
		hints = append(hints, fmt.Sprintf("↑ %d más", m.offset))
	}
	if end < len(m.jobs) {
		hints = append(hints, fmt.Sprintf("↓ %d más", len(m.jobs)-end))
	}
	scrollLine := ""
	if len(hints) > 0 {
		scrollLine = "\n  " + styleMuted.Render(strings.Join(hints, "   "))
	}

	body = hdr + "\n" + strings.Join(rows, "\n") + scrollLine
	help = "↑↓/jk navegar  Enter ver  a analizar  r recargar  q salir"
	return
}

// bodyJobDetail renderiza el detalle de una vacante seleccionada.
func (m Model) bodyJobDetail() (body, help string) {
	job := m.viewJob
	if job == nil {
		return "", "Esc volver"
	}

	var sb strings.Builder

	// ── Cabecera de la vacante ──
	sb.WriteString("\n  ")
	sb.WriteString(styleEmphasis.Render(job.Title))
	sb.WriteString("  " + styleMuted.Render("@ "+job.Company) + "\n\n  ")

	meta := []string{}
	if job.Portal != "" {
		meta = append(meta, styleSubtitle.Render(job.Portal))
	}
	if job.Modality != "" {
		meta = append(meta, job.Modality)
	}
	if job.Location != "" {
		meta = append(meta, job.Location)
	}
	sb.WriteString(strings.Join(meta, styleMuted.Render("  │  ")))
	sb.WriteString("\n  " + styleMuted.Render("Salario: ") + fmtSalary(job))
	if !job.PostedAt.IsZero() {
		sb.WriteString(styleMuted.Render("  │  Publicado: " + job.PostedAt.Format("2006-01-02")))
	}
	sb.WriteString("\n  " + styleMuted.Render(truncateTUI(job.URL, m.width-6)) + "\n\n")

	// ── Análisis IA ──
	divLine := styleDivider.Render(strings.Repeat("─", clamp(m.width-4, 4, 64)))
	sb.WriteString("  " + divLine + "\n")

	if a, ok := m.analyses[job.ID]; ok {
		scoreStr := scoreStyle(a.CompatibilityScore).Render(fmt.Sprintf("%d/100", a.CompatibilityScore))
		recStr := recommendStyle(a.Recommendation).Render("  " + strings.ToUpper(a.Recommendation) + "  ")
		sb.WriteString("  " + styleColHeader.Render("ANÁLISIS IA") +
			"  " + scoreStr + "  " + recStr + "\n\n")

		if len(a.Strengths) > 0 {
			sb.WriteString("  " + styleSuccess.Render("Fortalezas:") + "\n")
			for _, s := range a.Strengths {
				sb.WriteString("    • " + s + "\n")
			}
			sb.WriteString("\n")
		}
		if len(a.Gaps) > 0 {
			sb.WriteString("  " + styleDanger.Render("Brechas:") + "\n")
			for _, g := range a.Gaps {
				sb.WriteString("    • " + g + "\n")
			}
			sb.WriteString("\n")
		}
		if a.SalaryEstimate != "" {
			sb.WriteString("  " + styleMuted.Render("Estimado salarial: "+a.SalaryEstimate) + "\n\n")
		}
	} else {
		sb.WriteString("  " + styleMuted.Render("Sin análisis IA — presiona a para analizar") + "\n\n")
	}

	// ── Descripción truncada ──
	sb.WriteString("  " + divLine + "\n")
	sb.WriteString("  " + styleColHeader.Render("DESCRIPCIÓN") + "\n\n")

	wrapped := wordWrap(job.Description, m.width-6)
	descLines := strings.Split(wrapped, "\n")
	usedLines := strings.Count(sb.String(), "\n") + 1
	maxDesc := m.height - usedLines - 4
	if maxDesc < 2 {
		maxDesc = 2
	}
	if len(descLines) > maxDesc {
		descLines = descLines[:maxDesc]
		descLines = append(descLines, styleMuted.Render("  … (descripción truncada)"))
	}
	for _, l := range descLines {
		sb.WriteString("  " + l + "\n")
	}

	body = sb.String()
	help = "Esc/b volver  a analizar  q salir"
	return
}

// ── Comandos Bubble Tea ───────────────────────────────────────────────────────

func cmdLoadJobs(client *db.QdrantClient) tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		jobs, err := client.ListJobs(ctx, 100, 0)
		if err != nil {
			return fatalErrMsg{err}
		}
		return jobsLoadedMsg{jobs}
	}
}

func cmdAnalyze(svc *services.AnalyzerService, job *domain.Job, profilePath, cvPath string) tea.Cmd {
	return func() tea.Msg {
		profile := readFileFallback(profilePath, "[Perfil no configurado — usa --profile=ruta/perfil.md]")
		cv := readFileFallback(cvPath, "[CV no configurado — usa --cv=ruta/cv.md]")
		ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
		defer cancel()
		a, err := svc.AnalyzeJob(ctx, job, profile, cv)
		if err != nil {
			return analysisErrMsg{err}
		}
		return analysisOkMsg{jobID: job.ID, result: a}
	}
}

func cmdTick() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func (m Model) visibleRows() int {
	// header(1) + div(1) + colheader(1) + div(1) + footer(1) = 5 líneas fijas
	r := m.height - 5
	if r < 3 {
		r = 3
	}
	return r
}

func readFileFallback(path, fallback string) string {
	if path == "" {
		return fallback
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return fallback
	}
	return string(data)
}

func renderScore(score int) string {
	if score == 0 {
		return styleMuted.Render(" -- ")
	}
	return scoreStyle(score).Render(fmt.Sprintf("%3d%%", score))
}

func fmtSalary(job *domain.Job) string {
	if job.SalaryMin == 0 && job.SalaryMax == 0 {
		return styleMuted.Render("N/D")
	}
	if job.SalaryMax == 0 || job.SalaryMin == job.SalaryMax {
		return fmt.Sprintf("$%d", job.SalaryMin)
	}
	return fmt.Sprintf("$%d–%d", job.SalaryMin, job.SalaryMax)
}

func truncateTUI(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "…"
}

func padToHeight(s string, h int) string {
	n := strings.Count(s, "\n") + 1
	if n >= h {
		return s
	}
	return s + strings.Repeat("\n", h-n)
}

func wordWrap(s string, width int) string {
	if width <= 10 {
		width = 80
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return s
	}
	var lines []string
	cur := words[0]
	for _, w := range words[1:] {
		if len(cur)+1+len(w) <= width {
			cur += " " + w
		} else {
			lines = append(lines, cur)
			cur = w
		}
	}
	lines = append(lines, cur)
	return strings.Join(lines, "\n")
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}
