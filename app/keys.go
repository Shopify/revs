package app

import "github.com/charmbracelet/bubbles/key"

type listKeyMap struct {
	open        key.Binding
	done        key.Binding
	unsubscribe key.Binding
	sortTitle   key.Binding
	sortAuthor  key.Binding
	sortCreated key.Binding
	sortRepo    key.Binding
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
		sortTitle: key.NewBinding(
			key.WithKeys("T"),
			key.WithHelp("T", "sort by title"),
		),
		sortAuthor: key.NewBinding(
			key.WithKeys("A"),
			key.WithHelp("A", "sort by author"),
		),
		sortCreated: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("C", "sort by created"),
		),
		sortRepo: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "sort by repo"),
		),
	}
}
