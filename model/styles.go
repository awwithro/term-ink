package model

import "github.com/charmbracelet/lipgloss"

var headerStyle = lipgloss.NewStyle().
	Width(width + 2).
	Align(lipgloss.Center).
	Background(lipgloss.Color("#F25D94")).
	Foreground(lipgloss.Color("#FFF7DB")).
	PaddingBottom(0).
	MarginBottom(0)

var storyTextBoxStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Width(width).Height(15).
	BorderForeground(lipgloss.Color("#874BFD")).
	BorderTop(true).
	BorderLeft(true).
	BorderRight(true).
	BorderBottom(true).
	Align(lipgloss.Center)

var choiceBoxStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFF7DB")).
	Background(lipgloss.Color("#888B7E")).
	Padding(0, 2).
	Margin(0, 1)

var selectedChoiceBoxStyle = choiceBoxStyle.Copy().
	Background(lipgloss.Color("#F25D94")).
	Underline(true)
