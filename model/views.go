package model

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	width = 96
)

func Header() string {
	return "The Intercept"
}

func (m Model) View() string {
	doc := strings.Builder{}
	// Header
	doc.WriteString(lipgloss.Place(width, 1,
		lipgloss.Center, lipgloss.Center,
		headerStyle.Render(Header())) + "\n")

	// Text Box
	doc.WriteString(lipgloss.Place(width, 20,
		lipgloss.Center, lipgloss.Top,
		storyTextBoxStyle.Render(m.storyText.ToString())) + "\n")

	// Choices
	choiceSelector := lipgloss.Place(width, 1,
		lipgloss.Center, lipgloss.Top,
		lipgloss.JoinHorizontal(lipgloss.Left, m.choiceView()...))
	doc.WriteString(choiceSelector + "\n")

	// Footer
	helpView := lipgloss.Place(width, 4,
		lipgloss.Left, lipgloss.Bottom,
		lipgloss.JoinHorizontal(lipgloss.Left, m.help.View(keys)))
	doc.WriteString(helpView + "\n")
	return lipgloss.JoinHorizontal(lipgloss.Center, doc.String())
}

func (m Model) choiceView() []string {
	choices := []string{}
	for x, choice := range m.choices.choices {
		if x == int(m.choices.cursor) {
			choices = append(choices, selectedChoiceBoxStyle.Render(choice))
		} else {
			choices = append(choices, choiceBoxStyle.Render(choice))
		}
	}
	return choices
}
