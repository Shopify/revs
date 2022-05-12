package app

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/go-github/github"
	"github.com/pkg/browser"

	"github.com/campbel/revs/ghutil"
)

type Model struct {
	list         list.Model
	keys         *listKeyMap
	client       *github.Client
	sortAsceding bool
}

func NewModel(ctx context.Context, client *github.Client) *Model {

	uprs, err := ghutil.GetUnreadPullRequests(ctx, client)
	if err != nil {
		log.Fatal(err)
	}

	items := make([]list.Item, len(uprs))
	for i, upr := range uprs {
		item := &item{notification: upr.Notification, pr: upr.PR}
		items[i] = item
	}

	keys := newListKeyMap()

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Github Reviews"
	list.StatusMessageLifetime = time.Second * 10
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
			keys.sortAuthor,
			keys.sortCreated,
			keys.sortRepo,
			keys.sortTitle,
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
			item := m.list.SelectedItem().(*item)
			url := ghutil.GetPullRequestURL(item.notification)
			browser.OpenURL(url)
			return m, nil
		case key.Matches(msg, m.keys.done):
			item := m.list.SelectedItem().(*item)
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
			item := m.list.SelectedItem().(*item)
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
		case key.Matches(msg, m.keys.sortTitle):
			m.sortAsceding = !m.sortAsceding
			m.Sort(func(a *item, b *item) bool {
				if m.sortAsceding {
					return *a.notification.Subject.Title > *b.notification.Subject.Title
				}
				return *a.notification.Subject.Title < *b.notification.Subject.Title
			})
		case key.Matches(msg, m.keys.sortAuthor):
			m.sortAsceding = !m.sortAsceding
			m.Sort(func(a *item, b *item) bool {
				if m.sortAsceding {
					return *a.pr.User.Login > *b.pr.User.Login
				}
				return *a.pr.User.Login < *b.pr.User.Login
			})
		case key.Matches(msg, m.keys.sortCreated):
			m.sortAsceding = !m.sortAsceding
			m.Sort(func(a *item, b *item) bool {
				if m.sortAsceding {
					return a.pr.CreatedAt.UnixNano() > b.pr.CreatedAt.UnixNano()
				}
				return a.pr.CreatedAt.UnixNano() < b.pr.CreatedAt.UnixNano()
			})
		case key.Matches(msg, m.keys.sortRepo):
			m.sortAsceding = !m.sortAsceding
			m.Sort(func(a *item, b *item) bool {
				if m.sortAsceding {
					return *a.notification.Repository.FullName > *b.notification.Repository.FullName
				}
				return *a.notification.Repository.FullName < *b.notification.Repository.FullName
			})
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

func (m *Model) Sort(compare func(a, b *item) bool) {
	items := m.list.Items()
	sort.Slice(items, func(i, j int) bool {
		a := items[i].(*item)
		b := items[j].(*item)
		return compare(a, b)
	})
}
