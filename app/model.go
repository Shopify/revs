package app

import (
	"context"
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
	"github.com/pkg/browser"

	"github.com/campbel/revs/ghutil"
)

type Model struct {
	list   list.Model
	keys   *listKeyMap
	client *github.Client
}

func NewModel(ctx context.Context, client *github.Client) *Model {

	notifications, err := ghutil.GetUnreadPullRequests(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]list.Item, len(notifications))
	for i, notification := range notifications {
		items[i] = item{notification}
	}

	keys := newListKeyMap()

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Github Reviews"
	list.StatusMessageLifetime = 5
	list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.open,
			keys.done,
		}
	}
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.open,
			keys.done,
			keys.unsubscribe,
		}
	}
	return &Model{
		list:   list,
		keys:   keys,
		client: client,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *Model) View() string {
	return docStyle.Render(m.list.View())
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			url := ghutil.GetPullRequestURL(item.notification)
			browser.OpenURL(url)
			return m, nil
		case key.Matches(msg, m.keys.done):
			item := m.list.SelectedItem().(item)
			m.list.ToggleSpinner()
			_, err := m.client.Activity.MarkThreadRead(context.Background(), *item.notification.ID)
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
			_, err := m.client.Activity.DeleteThreadSubscription(context.Background(), *item.notification.ID)
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
