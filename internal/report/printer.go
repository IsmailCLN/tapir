package report

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func RenderResults() {
	p := tea.NewProgram(ResultView{results: testResults})
	_, err := p.Run()
	if err != nil {
		fmt.Printf("Failed to render results: %v\n", err)
		os.Exit(1)
	}
}

type ResultView struct {
	results  []TestResult
	quitting bool
	message  string
}

func (rv ResultView) Init() tea.Cmd { return nil }

func (rv ResultView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			rv.quitting = true
			return rv, tea.Quit
		case "p":
			err := os.WriteFile("tapir_output.md", []byte(rv.getMarkdownOutput()), 0644)
			rv.message = checkIOErr("Markdown saved to tapir_output.md", err)
		case "c":
			err := clipboard.WriteAll(rv.getRawOutput())
			rv.message = checkIOErr("Results copied to clipboard", err)
		}
	}
	return rv, nil
}

func (rv ResultView) View() string {
	if rv.quitting {
		return "Bye!\n"
	}

	var passed, failed int
	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(PurpleColor)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return HeaderStyle
			case col == 0:
				return TinyCell
			case col == 1:
				return LargeCell
			case col == 2:
				return MediumCell
			case col == 3:
				return SizeCell
			case col == 4:
				return SizeCell
			case col == 5:
				return NumCell
			case col == 6:
				return NumCell
			default:
				return MediumCell
			}
		}).
		Headers("✓", "Test Name", "Check", "Duration", "Size", "Expected", "Actual")

	for _, r := range rv.results {
		icon := "✓"
		if r.Result.Pass {
			passed++
			icon = green(icon)
		} else {
			failed++
			icon = red("✗")
		}

		row := []string{
			icon,
			r.Name,
			r.Category,
			r.Duration,
			r.ResponseSize,
			strconv.Itoa(r.ExpectedStatus),
			strconv.Itoa(r.ActualStatus),
		}
		if !r.Result.Pass && r.Result.Description != "" {
			row[1] += " — " + r.Result.Description
		}
		t.Row(row...)
	}

	return lipgloss.NewStyle().Margin(1, 2).Render(
		"🧪 Tapir Test Results:\n\n" + t.String() +
			bold(fmt.Sprintf("\nSummary: ✅ %d passed, ❌ %d failed", passed, failed)) +
			"\nPress 'q' to quit, 'p' to print to file, 'c' to copy to clipboard.\n\n" +
			rv.message,
	)
}

func (rv ResultView) getRawOutput() string {
	var sb strings.Builder
	var passed, failed int
	sb.WriteString("🧪 Tapir Test Results:\n\n")
	for _, r := range rv.results {
		if r.Result.Pass {
			passed++
			sb.WriteString(fmt.Sprintf("✓ %s: %s (%s)\n", r.Name, r.Category, r.Duration))
		} else {
			failed++
			sb.WriteString(fmt.Sprintf("✗ %s: %s | %s (%s)\n", r.Name, r.Category, r.Result.Description, r.Duration))
		}
	}
	sb.WriteString(fmt.Sprintf("\nSummary: %d passed, %d failed\n", passed, failed))
	return sb.String()
}

func checkIOErr(successMsg string, err error) string {
	if err != nil {
		return red("Error: " + err.Error())
	}
	return green(successMsg)
}

func (rv ResultView) getMarkdownOutput() string {
	var sb strings.Builder
	var passed, failed int

	sb.WriteString("# 🧪 Tapir Test Results\n\n")
	sb.WriteString("| ✓ | Test Name | Check | Duration | Size | Expected | Actual |\n")
	sb.WriteString("|---|-----------|-------|----------|------|----------|--------|\n")

	for _, r := range rv.results {
		icon := "✓"
		if r.Result.Pass {
			passed++
		} else {
			icon = "✗"
			failed++
		}

		name := r.Name
		if !r.Result.Pass && r.Result.Description != "" {
			name += " — " + r.Result.Description
		}

		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %d | %d |\n",
			icon, name, r.Category, r.Duration, r.ResponseSize, r.ExpectedStatus, r.ActualStatus))
	}

	sb.WriteString(fmt.Sprintf("\n**Summary:** ✅ %d passed, ❌ %d failed\n", passed, failed))

	return sb.String()
}
