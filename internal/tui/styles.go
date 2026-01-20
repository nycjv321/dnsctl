package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("86")  // Cyan
	secondaryColor = lipgloss.Color("245") // Gray
	successColor   = lipgloss.Color("82")  // Green
	errorColor     = lipgloss.Color("196") // Red
	warningColor   = lipgloss.Color("214") // Orange

	// Title style
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Subtitle style
	subtitleStyle = lipgloss.NewStyle().
			Foreground(secondaryColor)

	// Selected item style
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Normal item style
	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252"))

	// Help style
	helpStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			MarginTop(1)

	// Status styles
	statusStyle = lipgloss.NewStyle().
			Padding(0, 1)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	warningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Box style for info panels
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Description style
	descStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	// Key binding style
	keyStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Dimmed style for secondary info
	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240"))
)
