package ui

import "github.com/charmbracelet/lipgloss"

var (
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Render
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render

	PurpleColor    = lipgloss.Color("99")
	HeaderStyle    = lipgloss.NewStyle().Foreground(PurpleColor).Bold(true).Align(lipgloss.Center)
	ExtraSmallCell = lipgloss.NewStyle().Padding(0, 1).Width(3)
	SmallCell      = lipgloss.NewStyle().Padding(0, 1).Width(10)
	MediumCell     = lipgloss.NewStyle().Padding(0, 1).Width(20)
	LargeCell      = lipgloss.NewStyle().Padding(0, 1).Width(30)
	ExtraLargeCell = lipgloss.NewStyle().Padding(0, 1).Width(50)
)

func checkIOErr(successMsg string, err error) string {
	if err != nil {
		return red("Error: " + err.Error())
	}
	return green(successMsg)
}
