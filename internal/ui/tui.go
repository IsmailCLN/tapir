package ui

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/IsmailCLN/tapir/internal/helpers"
	"github.com/IsmailCLN/tapir/internal/runner"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	lgl "github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
)

type resultView struct {
	rows    [][]string
	results []runner.Result
	message string
}

func (rv resultView) Init() tea.Cmd { return nil }

func (rv resultView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return rv, tea.Quit
		case "c":
			err := clipboard.WriteAll(rv.getRawOutput())
			rv.message = checkIOErr("Results copied to clipboard", err)
		case "p":
			filename := "tapir-report-" + time.Now().Format("20060102") + ".md"
			err := os.WriteFile(filename, []byte(rv.getMarkdownOutput()), 0644)
			rv.message = checkIOErr("Markdown saved to "+filename, err)
		}
	}
	return rv, nil
}

func (rv resultView) View() string {
	t := ltable.New().
		Border(lgl.NormalBorder()).
		BorderStyle(lgl.NewStyle().Foreground(PurpleColor)).
		StyleFunc(styleCell).
		Headers("✓", "Suite", "Request", "Test", "Error").
		Rows(rv.rows...)

	return lgl.NewStyle().Margin(1, 2).
		Render("🧪 Tapir Test Results\n\n" + t.String() + "\nPress 'c' to copy, 'p' to save as markdown, 'q' to quit.\n\n" + rv.message)
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

func Render(results []runner.Result) error {
	rows := buildRows(results)
	rv := resultView{rows: rows, results: results}

	p := tea.NewProgram(rv, tea.WithAltScreen())
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
