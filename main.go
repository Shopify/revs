package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

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

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type model struct {
	list list.Model
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

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Github Notifications"
	return model{
		list: list,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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

	// for _, notification := range notifications {

	// 	id := getPullRequestID(notification)
	// 	merged, _, err := client.PullRequests.IsMerged(ctx, *notification.Repository.Owner.Login, *notification.Repository.Name, id)
	// 	if err != nil {
	// 		continue
	// 	}

	// 	if merged {
	// 		_, _ = client.Activity.MarkThreadRead(ctx, *notification.ID)
	// 		continue
	// 	}

	// 	prURL := getPullRequestURL(notification)
	// 	fmt.Println(notification.GetSubject().GetTitle())
	// 	fmt.Println("Opening", prURL)
	// 	browser.OpenURL(prURL)
	// 	fmt.Println("Press enter to continue...")
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')

	// 	_, _ = client.Activity.MarkThreadRead(ctx, *notification.ID)
	// }
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
