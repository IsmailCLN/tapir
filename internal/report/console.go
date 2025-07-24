package report

import (
	"fmt"
	"os"

	"github.com/IsmailCLN/tapir/internal/assert"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	var passed, failed int

	// Sadece başlığa stil uyguluyoruz
	output := lipgloss.NewStyle().Margin(1, 2).Render("🧪 Tapir Test Results:") + "\n\n"

	for _, r := range m.results {
		if r.Result.Pass {
			passed++
			output += green(fmt.Sprintf("✓ %s", r.Name)) + "\n"
		} else {
			failed++
			output += red(fmt.Sprintf("✗ %s: %s", r.Name, r.Result.Description)) + "\n"
		}
	}

	output += "\n"
	output += bold(fmt.Sprintf("Summary: ✅ %d passed, ❌ %d failed", passed, failed))
	output += "\nPress 'q' to quit."

	return output
}

// Stil tanımları
var green = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Render
var red = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render
var bold = lipgloss.NewStyle().Bold(true).Render
