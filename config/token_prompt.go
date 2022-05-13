package config

import (
	"fmt"

	"github.com/Shopify/revs/bubble/text"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	promptMessage = `
Add a Github Personal Access Token (PAT)
github.com/settings/tokens/new?scopes=repo,notifications

This is used to access your notifications and information about pull requests you've been asked to review.
`
)

type model struct {
	textInput textinput.Model
	text      text.Model
	err       error
}

func initialModel() model {

	ti := textinput.New()
	ti.Placeholder = "GitHub Personal Access Token (PAT)"
	ti.Focus()
	ti.CharLimit = 156

	txt := text.NewModel(promptMessage)

	return model{
		textInput: ti,
		text:      txt,
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		tCmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	case error:
		m.err = msg
		return m, nil
	}

	m.textInput, tiCmd = m.textInput.Update(msg)
	m.text, tCmd = m.text.Update(msg)

	return m, tea.Batch(tiCmd, tCmd)
}

func (m model) View() string {
	return fmt.Sprintf("%s\n\n%s\n", m.text.View(), m.textInput.View())
}
