package ui

import "github.com/charmbracelet/lipgloss"

var headerStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderBottom(true).
	BorderTop(false).BorderRight(false).BorderLeft(false).
	BorderForeground(lipgloss.Color("198")).
	Bold(true).
	Foreground(lipgloss.Color("99")).
	MarginRight(1).
	MarginLeft(1).
	PaddingTop(1).
	PaddingLeft(4)

var footerStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderTop(true).
	BorderRight(false).BorderLeft(false).BorderBottom(false).
	BorderForeground(lipgloss.Color("198")).
	Bold(true).
	Foreground(lipgloss.Color("99")).
	MarginRight(1).
	MarginLeft(1)
	// PaddingLeft(4).Align(lipgloss.Right)

var viewportStyle = lipgloss.NewStyle().PaddingLeft(3).
	BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).BorderForeground(lipgloss.Color("236"))

var appListStyle = lipgloss.NewStyle().MarginRight(3).MarginLeft(2)
