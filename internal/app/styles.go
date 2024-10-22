package app

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color definitions
	primaryBlue   = lipgloss.Color("#2A9FD6")
	successGreen  = lipgloss.Color("#04B575")
	warningYellow = lipgloss.Color("#FFCC00")
	errorRed      = lipgloss.Color("#FF0033")
	textWhite     = lipgloss.Color("#FFFDF5")
	subtleGray    = lipgloss.Color("#666666")

	// appStyle defines the main application padding
	appStyle = lipgloss.NewStyle().
			Padding(1, 2)

	// titleStyle defines the style for the application title
	titleStyle = lipgloss.NewStyle().
			Foreground(textWhite).
			Background(primaryBlue).
			Padding(0, 1).
			Bold(true)

	// viewportStyle defines the style for the message viewport
	viewportStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(subtleGray).
			PaddingLeft(1).
			MarginTop(1).
			MarginBottom(1)

	// Message styles
	successStyle = lipgloss.NewStyle().
			Foreground(successGreen)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorRed)

	blockedDomainStyle = lipgloss.NewStyle().
				Foreground(warningYellow)

	// Help style for instructions
	helpStyle = lipgloss.NewStyle().
			Foreground(subtleGray).
			Italic(true)
)

// Expose styles for main package
func ErrorStyle() lipgloss.Style {
	return errorStyle
}

func SuccessStyle() lipgloss.Style {
	return successStyle
}
