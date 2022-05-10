package main

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
)

type listKeyMap struct {
	open key.Binding
	done key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		open: key.NewBinding(
			key.WithKeys("o", "enter"),
			key.WithHelp("o/enter", "open selection"),
		),
		done: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "marke as read"),
		),
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.list.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.open):
			item := m.list.SelectedItem().(item)
			url := getPullRequestURL(item.notification)
			browser.OpenURL(url)
			return m, nil
		case key.Matches(msg, m.keys.done):
			item := m.list.SelectedItem().(item)
			client.Activity.MarkThreadRead(context.Background(), *item.notification.ID)
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Marked " + *item.notification.Subject.Title + " as read."))
			m.list.RemoveItem(m.list.Index())
			return m, statusCmd
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
