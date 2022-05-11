package app

import "github.com/charmbracelet/bubbles/key"

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
