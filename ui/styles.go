package ui

import "github.com/charmbracelet/lipgloss"

var HeaderStyle = lipgloss.NewStyle().
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

var FooterStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.ThickBorder()).
	BorderTop(true).
	BorderRight(false).BorderLeft(false).BorderBottom(false).
	BorderForeground(lipgloss.Color("198")).
	Bold(true).
	Foreground(lipgloss.Color("99")).
	MarginRight(1).
	MarginLeft(1).
	PaddingLeft(4)

var ViewportStyle = lipgloss.NewStyle().MarginLeft(100)
