package main

import (
	"context"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pkg/browser"
)

type listKeyMap struct {
	open        key.Binding
	done        key.Binding
	unsubscribe key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		open: key.NewBinding(
			key.WithKeys("o", "enter"),
			key.WithHelp("o/enter", "open"),
		),
		done: key.NewBinding(
			key.WithKeys("e", "I"),
			key.WithHelp("e/I", "mark as read"),
		),
		unsubscribe: key.NewBinding(
			key.WithKeys("U"),
			key.WithHelp("U", "unsubscribe"),
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
			m.list.ToggleSpinner()
			_, err := client.Activity.MarkThreadRead(context.Background(), *item.notification.ID)
			m.list.ToggleSpinner()
			if err != nil {
				statusCmd := m.list.NewStatusMessage(errorMessageStyle("Error while marking " + *item.notification.Subject.Title + " as read"))
				return m, statusCmd
			}
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Marked " + *item.notification.Subject.Title + " as read."))
			m.list.RemoveItem(m.list.Index())
			return m, statusCmd
		case key.Matches(msg, m.keys.unsubscribe):
			item := m.list.SelectedItem().(item)
			m.list.ToggleSpinner()
			_, err := client.Activity.DeleteThreadSubscription(context.Background(), *item.notification.ID)
			m.list.ToggleSpinner()
			if err != nil {
				statusCmd := m.list.NewStatusMessage(errorMessageStyle("Error while unsubscribing from " + *item.notification.Subject.Title + "."))
				return m, statusCmd
			}
			statusCmd := m.list.NewStatusMessage(statusMessageStyle("Unsubscribed from " + *item.notification.Subject.Title + "."))
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
