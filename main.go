package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/go-github/v44/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

const (
	ReasonAssigned        = "assigned"
	ReasonComment         = "comment"
	ReasonMention         = "mention"
	ReasonReviewRequested = "review_requested"
	ReasonTeamMention     = "team_mention"
	ReasonAuthor          = "author"
)

var (
	// From lowest to highest priority
	ReasonPriority = []string{ReasonAssigned, ReasonAuthor, ReasonReviewRequested, ReasonTeamMention, ReasonMention, ReasonComment}
)

var (
	docStyle           = lipgloss.NewStyle().Margin(1, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type model struct {
	list list.Model
	keys *listKeyMap
}

type item struct {
	notification *github.Notification
}

func (i item) Title() string       { return *i.notification.Subject.Title }
func (i item) Description() string { return *i.notification.Repository.FullName }
func (i item) FilterValue() string {
	return *i.notification.Subject.Title + " " + *i.notification.Repository.FullName
}

func initialModel(notifications []*github.Notification) model {

	items := make([]list.Item, len(notifications))
	for i, notification := range notifications {
		items[i] = item{notification}
	}

	keys := newListKeyMap()

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Github Notifications"
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
		}
	}
	return model{
		list: list,
		keys: keys,
	}
}

func (m model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var client *github.Client

func main() {

	ctx := context.Background()

	// Authenticate with static token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create a new GitHub client
	client = github.NewClient(tc)

	// list all notifications
	notificationList, _, err := client.Activity.ListNotifications(ctx, nil)
	if err != nil {
		logrus.Fatal(err)
	}

	// filter to only unread pr notifications
	notifications := make([]*github.Notification, 0)
	for _, notification := range notificationList {
		if *notification.Unread && *notification.Subject.Type == "PullRequest" {
			notifications = append(notifications, notification)
		}
	}

	sort.Slice(notifications, func(i, j int) bool {
		if *notifications[i].Repository.FullName != *notifications[j].Repository.FullName {
			return *notifications[i].Repository.FullName < *notifications[j].Repository.FullName
		}
		return slices.Index(ReasonPriority, *notifications[i].Reason) > slices.Index(ReasonPriority, *notifications[j].Reason)
	})

	p := tea.NewProgram(initialModel(notifications), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Printf("Oh no an error: %v", err)
		os.Exit(1)
	}
}

func getPullRequestURL(notification *github.Notification) string {
	return fmt.Sprintf("https://github.com/%s/pull/%d?notification_referrer_id=%s", *notification.Repository.FullName, getPullRequestID(notification), *notification.ID)
}

func getPullRequestID(notification *github.Notification) int {
	parts := strings.Split(*notification.Subject.URL, "/")
	val, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1
	}
	return val
}
