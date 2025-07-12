package tui

import "github.com/charmbracelet/lipgloss"

// Styles for layout
var (
	AppStyle       = lipgloss.NewStyle().Padding(1, 2)
	ListStyle      = lipgloss.NewStyle().Padding(1, 0)
	TitleStyle     = lipgloss.NewStyle().Bold(true).Underline(true).Padding(0, 1)
	FocusedStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder())
	UnfocusedStyle = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
)
