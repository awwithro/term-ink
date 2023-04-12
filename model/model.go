package model

import (
	"context"
	"log"
	"strings"

	"github.com/awwithro/term_ink/pkg/inkgrpc"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Model struct {
	storyText        StoryText
	choices          Choices
	storyId          string
	client           inkgrpc.StoryClient
	latestStoryState *inkgrpc.StoryState
	ready            bool
	help             help.Model
}

type StoryText struct {
	text []string
}

func (s StoryText) ToString() string {
	sb := strings.Builder{}
	for _, s := range s.text {
		sb.WriteString(s + "\n")
	}
	return sb.String()
}

type continueStoryMessage struct{ Story *inkgrpc.StoryState }
type errMsg struct{ error }

func NewStoryText() StoryText {
	return StoryText{
		text: []string{},
	}
}

type Choices struct {
	choices []string // items on the to-do list
	cursor  int32    // which to-do list item our cursor is pointing at
}

func NewChoices() Choices {
	return Choices{
		choices: []string{},
		cursor:  0,
	}
}

func (c Choices) HasChoices() bool {
	return len(c.choices) > 0
}
func (c *Choices) Inc() {
	if int(c.cursor) < len(c.choices)-1 {
		c.cursor++
	}
}
func (c *Choices) Dec() {
	if c.cursor > 0 {
		c.cursor--
	}
}
func (c *Choices) AddChoice(choice string) {
	c.choices = append(c.choices, choice)
}

func InitialModel(client inkgrpc.StoryClient) Model {
	m := Model{
		storyText: NewStoryText(),
		choices:   NewChoices(),
		client:    client,
		help:      help.New(),
	}
	m.help.Width = width
	stories, err := client.ListStories(context.TODO(), &emptypb.Empty{})
	if err != nil {
		log.Fatal(err)
	}
	storyTitle := stories.StoryTitles[0]
	start, err := client.StartStory(context.TODO(), &inkgrpc.StartStoryRequest{StoryTitle: storyTitle})
	if err != nil {
		log.Fatal(err)
	}
	m.storyId = start.Id
	r, err := client.Continue(context.TODO(), &inkgrpc.ContinueRequest{Id: m.storyId})
	if err != nil {
		log.Fatal(err)
	}
	m.storyText.text = append(m.storyText.text, r.Story.Text)
	m.choices = NewChoices()
	for _, choice := range r.Story.Choices {
		m.choices.AddChoice(choice.Text)
	}
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Left):
			if m.choices.HasChoices() {
				m.choices.Dec()
			}
		case key.Matches(msg, keys.Right):
			if m.choices.HasChoices() {
				m.choices.Inc()
			}
		case key.Matches(msg, keys.Select):
			if m.choices.HasChoices() {
				// reset the text
				m.storyText = NewStoryText()
				return m, m.sendChoice
			}
		case key.Matches(msg, keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	case continueStoryMessage:
		m.latestStoryState = msg.Story
		m.continueStory()
		if m.latestStoryState.CanContinue {
			m.continueStory()
		}
		return m, nil

	case errMsg:
		log.Fatalf("Error: %v", msg)
	}
	return m, nil
}

func (m Model) sendChoice() tea.Msg {
	_, err := m.client.ChooseChoice(context.TODO(), &inkgrpc.ChooseChoiceRequest{
		ChoiceIndex: m.choices.cursor,
		Id:          m.storyId,
	})
	if err != nil {
		return errMsg{err}
	}

	cont, err := m.client.Continue(context.TODO(), &inkgrpc.ContinueRequest{Id: m.storyId})
	if err != nil {
		return errMsg{err}
	}
	return continueStoryMessage{
		Story: cont.Story,
	}

}

func (m *Model) continueStory() {
	m.storyText.text = append(m.storyText.text, strings.ReplaceAll(m.latestStoryState.Text, "\n", ""))
	m.choices = NewChoices()
	for _, choice := range m.latestStoryState.Choices {
		m.choices.AddChoice(choice.Text)
	}
	if m.latestStoryState.CanContinue {
		r, err := m.client.Continue(context.TODO(), &inkgrpc.ContinueRequest{Id: m.storyId})
		if err != nil {
			log.Fatal(err)
		}
		m.latestStoryState = r.Story
		m.continueStory()
	}
}
