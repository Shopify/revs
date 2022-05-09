package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v44/github"
	"github.com/pkg/browser"
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

func main() {
	ctx := context.Background()

	// Authenticate with static token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create a new GitHub client
	client := github.NewClient(tc)

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
		return slices.Index(ReasonPriority, *notifications[i].Reason) > slices.Index(ReasonPriority, *notifications[j].Reason)
	})

	fmt.Printf("Starting PR workflow... (%d prs to review)\n", len(notifications))

	for _, notification := range notifications {

		id := getPullRequestID(notification)
		merged, _, err := client.PullRequests.IsMerged(ctx, *notification.Repository.Owner.Login, *notification.Repository.Name, id)
		if err != nil {
			continue
		}

		if merged {
			_, _ = client.Activity.MarkThreadRead(ctx, *notification.ID)
			continue
		}

		prURL := getPullRequestURL(notification)
		fmt.Println(notification.GetSubject().GetTitle())
		fmt.Println("Opening", prURL)
		browser.OpenURL(prURL)
		fmt.Println("Press enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		_, _ = client.Activity.MarkThreadRead(ctx, *notification.ID)
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
