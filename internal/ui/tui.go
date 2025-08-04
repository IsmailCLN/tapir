package ui

import (
	"fmt"
	"runtime"
	"strings"
	"text/tabwriter"

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
		Render("🧪 Tapir Test Results\n\n" + t.String() + "\nPress 'c' to copy, 'q' to quit.\n\n" + rv.message)
}

// –– Helpers –– //
func checkIOErr(successMsg string, err error) string {
	if err != nil {
		return red("Error: " + err.Error())
	}
	return green(successMsg)
}

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
			errMsg = r.Err.Error()
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
		status := "✗"
		if r.Passed {
			status = "✓"
			passed++
		} else {
			failed++
		}

		errMsg := ""
		if r.Err != nil {
			errMsg = r.Err.Error()
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			status,
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
