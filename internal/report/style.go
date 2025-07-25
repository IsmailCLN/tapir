package report

import "github.com/charmbracelet/lipgloss"

var (
	green = lipgloss.NewStyle().Foreground(lipgloss.Color("#22c55e")).Render
	red   = lipgloss.NewStyle().Foreground(lipgloss.Color("#ef4444")).Render
	bold  = lipgloss.NewStyle().Bold(true).Render

	PurpleColor = lipgloss.Color("99")
	HeaderStyle = lipgloss.NewStyle().Foreground(PurpleColor).Bold(true).Align(lipgloss.Center)
	TinyCell    = lipgloss.NewStyle().Padding(0, 1).Width(3)
	MediumCell  = lipgloss.NewStyle().Padding(0, 1).Width(15).Align(lipgloss.Center)
	LargeCell   = lipgloss.NewStyle().Padding(0, 1).Width(30)
	SizeCell    = lipgloss.NewStyle().Padding(0, 1).Width(10).Align(lipgloss.Right)
	NumCell     = lipgloss.NewStyle().Padding(0, 1).Width(8).Align(lipgloss.Right)
)
