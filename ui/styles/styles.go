package ui

import "github.com/charmbracelet/lipgloss"

// HeaderStyle sets the style of our header
var HeaderStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderBottom(true).
	BorderTop(false).BorderRight(false).BorderLeft(false).
	BorderForeground(lipgloss.Color("198")).
	Bold(true).
	Foreground(lipgloss.Color("99")).
	MarginTop(1)

// FooterStyle sets the style of our footer
var FooterStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderTop(true).
	BorderRight(false).BorderLeft(false).BorderBottom(false).
	BorderForeground(lipgloss.Color("198")).
	Bold(true).
	Foreground(lipgloss.Color("99"))

// ViewportStyle sets the style of the viewport
var ViewportStyle = lipgloss.NewStyle().
	PaddingLeft(3).
	BorderStyle(lipgloss.NormalBorder()).
	BorderLeft(true).
	BorderForeground(lipgloss.Color("236"))

// AppListStyle sets the style of our application list
var AppListStyle = lipgloss.NewStyle().
	MarginRight(2).
	MarginLeft(0)
