package report

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/IsmailCLN/tapir/internal/assert"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type TestResult struct {
	Name           string
	Result         assert.AssertionResult
	Duration       string
	ResponseSize   string
	ExpectedStatus int
	ActualStatus   int
}

var testResults []TestResult

func PrintResult(testName string, result assert.AssertionResult, durationMs int64, sizeBytes int64, expected, actual int) {
	testResults = append(testResults, TestResult{
		Name:           testName,
		Result:         result,
		Duration:       formatDuration(durationMs),
		ResponseSize:   formatSize(sizeBytes),
		ExpectedStatus: expected,
		ActualStatus:   actual,
	})
}

func formatDuration(ms int64) string {
	if ms < 1000 {
		return fmt.Sprintf("%d ms", ms)
	}
	return fmt.Sprintf("%.2f s", float64(ms)/1000.0)
}

func formatSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	return fmt.Sprintf("%.2f KB", float64(bytes)/1024.0)
}

func RenderResults() {
	p := tea.NewProgram(ResultView{results: testResults})
	model, err := p.Run()
	if err != nil {
		fmt.Printf("Failed to render results: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Final state:", model)
}

type ResultView struct {
	results  []TestResult
	quitting bool
	message  string
}

func (rv ResultView) Init() tea.Cmd {
	return nil
}

func (rv ResultView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			rv.quitting = true
			return rv, tea.Quit
		case "p":
			err := os.WriteFile("tapir_output.txt", []byte(rv.getRawOutput()), 0644)
			if err != nil {
				rv.message = red("Failed to write to file: " + err.Error())
			} else {
				rv.message = green("Results written to tapir_output.txt")
			}
			return rv, nil
		case "c":
			err := clipboard.WriteAll(rv.getRawOutput())
			if err != nil {
				rv.message = red("Failed to copy to clipboard: " + err.Error())
			} else {
				rv.message = green("Results copied to clipboard")
			}
			return rv, nil
		}
	}
	return rv, nil
}

func (rv ResultView) View() string {
	if rv.quitting {
		return "Bye!\n"
	}

	var passed, failed int

	// Table styling
	purple := lipgloss.Color("99")

	headerStyle := lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	statusCellStyle := lipgloss.NewStyle().Padding(0, 1).Width(3)
	nameCellStyle := lipgloss.NewStyle().Padding(0, 1).Width(30)
	durationAndSizeCellStyle := lipgloss.NewStyle().Padding(0, 1).Width(10).Align(lipgloss.Right)
	numCellStyle := lipgloss.NewStyle().Padding(0, 1).Width(8).Align(lipgloss.Right)

	t := table.New().
		Border(lipgloss.NormalBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(purple)).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			case col == 0:
				return statusCellStyle
			case col == 1:
				return nameCellStyle
			case col == 2:
				return durationAndSizeCellStyle
			case col == 3:
				return durationAndSizeCellStyle
			case col == 4:
				return numCellStyle
			case col == 5:
				return numCellStyle
			default:
				return nameCellStyle
			}
		}).
		Headers("✓", "Test Name", "Duration", "Size", "Expected", "Actual")

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
			r.Duration,
			r.ResponseSize,
			strconv.Itoa(r.ExpectedStatus),
			strconv.Itoa(r.ActualStatus),
		}

		if !r.Result.Pass && r.Result.Description != "" {
			row[1] = row[1] + " — " + r.Result.Description
		}

		t.Row(row...)
	}

	summary := bold(fmt.Sprintf("\nSummary: ✅ %d passed, ❌ %d failed", passed, failed))
	footer := "\nPress 'q' to quit, 'p' to print to file, 'c' to copy to clipboard."

	return lipgloss.NewStyle().Margin(1, 2).Render("🧪 Tapir Test Results:\n\n" + t.String() + summary + footer + "\n\n" + rv.message)
}

func (rv ResultView) getRawOutput() string {
	var sb strings.Builder
	sb.WriteString("🧪 Tapir Test Results:\n\n")

	var passed, failed int
	for _, r := range rv.results {
		if r.Result.Pass {
			passed++
			sb.WriteString(fmt.Sprintf("✓ %s (%s)\n", r.Name, r.Duration))
		} else {
			failed++
			sb.WriteString(fmt.Sprintf("✗ %s: %s (%s)\n", r.Name, r.Result.Description, r.Duration))
		}
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d passed, %d failed\n", passed, failed))
	return sb.String()
}

// Stil tanımları
var green = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Render
var red = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render
var bold = lipgloss.NewStyle().Bold(true).Render
