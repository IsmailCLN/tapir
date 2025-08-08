package ui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/IsmailCLN/tapir/internal/domain"
	"github.com/IsmailCLN/tapir/internal/helpers"
	"github.com/IsmailCLN/tapir/internal/parser"
	"github.com/IsmailCLN/tapir/internal/runner"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	lgl "github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
)

type resultView struct {
	rows       [][]string
	results    []runner.Result
	message    string
	suitePaths []string
	isRunning  bool
	lastRerun  time.Time
}

type rerunDoneMsg struct {
	results []runner.Result
	err     error
}

func (rv resultView) Init() tea.Cmd { return nil }

func (rv resultView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m := msg.(type) {
	case rerunDoneMsg:
		rv.isRunning = false
		rv.lastRerun = time.Now()

		if m.err != nil {
			rv.message = checkIOErr("run error", m.err) // kırmızı
			return rv, nil
		}
		rv.results = m.results
		rv.rows = buildRows(m.results)
		rv.message = checkIOErr("Re-run completed at "+rv.lastRerun.Format("15:04:05"), nil) // yeşil
		return rv, nil

	case tea.KeyMsg:
		if h, ok := rv.keyHandlers()[m.String()]; ok {
			return h(rv)
		}
	}
	return rv, nil
}

func (rv resultView) keyHandlers() map[string]func(resultView) (tea.Model, tea.Cmd) {
	return map[string]func(resultView) (tea.Model, tea.Cmd){
		"r":      handleRerun,
		"c":      handleCopy,
		"p":      handleSaveMarkdown,
		"q":      handleQuit,
		"esc":    handleQuit,
		"ctrl+c": handleQuit,
	}
}

func handleQuit(rv resultView) (tea.Model, tea.Cmd) { return rv, tea.Quit }

func handleCopy(rv resultView) (tea.Model, tea.Cmd) {
	err := clipboard.WriteAll(rv.getRawOutput())
	rv.message = checkIOErr("Results copied to clipboard", err)
	return rv, nil
}

func handleSaveMarkdown(rv resultView) (tea.Model, tea.Cmd) {
	filename := "tapir-report-" + time.Now().Format("20060102") + ".md"
	err := os.WriteFile(filename, []byte(rv.getMarkdownOutput()), 0644)
	rv.message = checkIOErr("Markdown saved to "+filename, err)
	return rv, nil
}

// !!! ASENKRON RERUN !!!
// Ağır işleri burada yapmıyoruz; sadece guard'lar ve komut dönüşü var.
func handleRerun(rv resultView) (tea.Model, tea.Cmd) {
	// Çalışıyorsa tekrar başlatma → kırmızı mesaj
	if rv.isRunning {
		rv.message = checkIOErr("Already running, please wait…", errors.New("busy"))
		return rv, nil
	}
	// 3 sn cooldown → kırmızı mesaj
	if !rv.lastRerun.IsZero() && time.Since(rv.lastRerun) < 3*time.Second {
		rv.message = checkIOErr("Please wait a moment before re-running", errors.New("cooldown"))
		return rv, nil
	}

	// Arka planda çalıştır
	rv.isRunning = true
	rv.message = checkIOErr("Re-running…", nil)
	return rv, rerunCmd(rv.suitePaths)
}

// Arka plan işi: YAML'leri yükle, runner'ı çalıştır, sonucu mesajla geri dön.
func rerunCmd(paths []string) tea.Cmd {
	return func() tea.Msg {
		var allSuites []domain.TestSuite
		for _, p := range paths {
			s, err := parser.LoadTestSuite(p)
			if err != nil {
				return rerunDoneMsg{err: fmt.Errorf("reload error: %w", err)}
			}
			allSuites = append(allSuites, s...)
		}
		res, err := runner.Run(context.Background(), allSuites)
		return rerunDoneMsg{results: res, err: err}
	}
}

func (rv resultView) View() string {
	t := ltable.New().
		Border(lgl.NormalBorder()).
		BorderStyle(lgl.NewStyle().Foreground(PurpleColor)).
		StyleFunc(styleCell).
		Headers("✓", "Suite", "Request", "Test", "Error").
		Rows(rv.rows...)

	return lgl.NewStyle().Margin(1, 2).
		Render("🧪 Tapir Test Results\n\n" + t.String() + "\nPress 'c' to copy, 'p' to save as markdown, 'r' to rerun, 'q' to quit.\n\n" + rv.message)
}

// –– Helpers –– //
func styleCell(row, col int) lgl.Style {
	var s lgl.Style
	switch {
	case row == ltable.HeaderRow:
		s = HeaderStyle
	case col == 0:
		s = ExtraSmallCell
	case col == 1, col == 2, col == 3:
		s = LargeCell
	case col == 4:
		s = ExtraLargeCell
	default:
		s = MediumCell
	}

	if row != ltable.HeaderRow {
		if row%2 == 0 {
			s = s.Foreground(lgl.Color("#bfbfbf"))
		} else {
			s = s.Foreground(lgl.Color("#ffffff"))
		}
	}

	return s
}

func buildRows(results []runner.Result) [][]string {
	rows := make([][]string, len(results))
	for i, r := range results {
		icon := green("✓")
		if !r.Passed {
			icon = red("✗")
		}

		errMsg := ""
		if r.Err != nil {
			errMsg = helpers.Sanitize(r.Err.Error())
		}

		rows[i] = []string{
			icon,
			r.Suite,
			r.Request,
			r.TestName,
			errMsg,
		}
	}
	return rows
}

func Render(paths []string, results []runner.Result) error {
	rv := resultView{
		rows:       buildRows(results),
		results:    results,
		suitePaths: paths,
	}
	p := tea.NewProgram(rv)
	_, err := p.Run()
	return err
}

func (rv resultView) getRawOutput() string {
	var b strings.Builder
	const minColWidth = 4
	var passed, failed int

	w := tabwriter.NewWriter(&b, minColWidth, 0, 3, ' ', 0)
	fmt.Fprintln(w, "✓\tSuite\tRequest\tTest\tError")

	for _, r := range rv.results {
		icon := "✗"
		if r.Passed {
			icon = "✓"
			passed++
		} else {
			failed++
		}

		errMsg := ""
		if r.Err != nil {
			errMsg = helpers.Sanitize(r.Err.Error())
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			icon,
			r.Suite,
			r.Request,
			r.TestName,
			errMsg,
		)
	}

	fmt.Fprintln(w)
	fmt.Fprintf(w, "Passed:\t%d\tFailed:\t%d\tTotal:\t%d\n", passed, failed, passed+failed)

	_ = w.Flush()

	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(b.String(), "\n", "\r\n")
	}
	return b.String()
}

func (rv resultView) getMarkdownOutput() string {
	var sb strings.Builder
	var passed, failed int

	sb.WriteString("# 🧪 Tapir Test Results\n\n")
	sb.WriteString("| ✓ | Suite | Request | Test Name | Description |\n")
	sb.WriteString("|---|-------|---------|-----------|-------------|\n")

	for _, r := range rv.results {
		icon := "✓"
		if r.Passed {
			passed++
		} else {
			icon = "✗"
			failed++
		}

		errMsg := ""
		if r.Err != nil {
			errMsg = helpers.Sanitize(r.Err.Error())
		}

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
			icon, r.Suite, r.Request, r.TestName, errMsg))
	}
	sb.WriteString(fmt.Sprintf("\n**Summary:** ✅ %d passed, ❌ %d failed\n", passed, failed))

	return sb.String()
}
