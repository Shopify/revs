package text

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mitchellh/go-wordwrap"
)

type Model struct {
	Value       string
	screenWidth int
}

func NewModel(value string) Model {
	return Model{Value: value}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
	}
	return m, cmd
}

func (m Model) View() string {
	return wordwrap.WrapString(m.Value, uint(m.screenWidth))
}
