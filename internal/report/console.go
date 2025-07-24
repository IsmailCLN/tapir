package report

import (
	"fmt"
	"os"
	"strings"

	"github.com/IsmailCLN/tapir/internal/assert"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/atotto/clipboard"
)


// Bu struct dışarıdan erişilebilsin diye büyük harfle
type TestResult struct {
	Name   string
	Result assert.AssertionResult
}

// geçici olarak sonuçları burada biriktiriyoruz
var testResults []TestResult

// Anlık olarak test sonucu ekler
func PrintResult(testName string, result assert.AssertionResult) {
	testResults = append(testResults, TestResult{Name: testName, Result: result})
}

// Tüm sonuçları kullanıcıya göster
func RenderResults() {
	p := tea.NewProgram(model{results: testResults})
	if err := p.Start(); err != nil {
		fmt.Printf("Failed to render results: %v\n", err)
		os.Exit(1)
	}
}

type model struct {
	results  []TestResult
	quitting bool
	message  string // Alt mesaj göstermek için
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "p":
			err := os.WriteFile("tapir_output.txt", []byte(m.getRawOutput()), 0644)
			if err != nil {
				m.message = red("Failed to write to file: " + err.Error())
			} else {
				m.message = green("Results written to tapir_output.txt")
			}
			return m, nil
		case "c":
			err := clipboard.WriteAll(m.getRawOutput())
			if err != nil {
				m.message = red("Failed to copy to clipboard: " + err.Error())
			} else {
				m.message = green("Results copied to clipboard")
			}
			return m, nil
		}
	}
	return m, nil
}


func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var passed, failed int
	var output strings.Builder

	output.WriteString(lipgloss.NewStyle().Margin(1, 2).Render("🧪 Tapir Test Results:") + "\n\n")

	for _, r := range m.results {
		if r.Result.Pass {
			passed++
			output.WriteString(green(fmt.Sprintf("✓ %s", r.Name)) + "\n")
		} else {
			failed++
			output.WriteString(red(fmt.Sprintf("✗ %s: %s", r.Name, r.Result.Description)) + "\n")
		}
	}

	output.WriteString("\n")
	output.WriteString(bold(fmt.Sprintf("Summary: ✅ %d passed, ❌ %d failed", passed, failed)) + "\n")
	output.WriteString("Press 'q' to quit, 'p' to print to file, 'c' to copy to clipboard.\n")

	if m.message != "" {
		output.WriteString("\n" + m.message + "\n")
	}

	return output.String()
}

func (m model) getRawOutput() string {
	var sb strings.Builder
	sb.WriteString("🧪 Tapir Test Results:\n\n")

	var passed, failed int
	for _, r := range m.results {
		if r.Result.Pass {
			passed++
			sb.WriteString(fmt.Sprintf("✓ %s\n", r.Name))
		} else {
			failed++
			sb.WriteString(fmt.Sprintf("✗ %s: %s\n", r.Name, r.Result.Description))
		}
	}

	sb.WriteString(fmt.Sprintf("\nSummary: %d passed, %d failed\n", passed, failed))
	return sb.String()
}


// Stil tanımları
var green = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Render
var red = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render
var bold = lipgloss.NewStyle().Bold(true).Render
