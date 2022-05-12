package ghutil

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/exp/slices"
	"golang.org/x/oauth2"
)

func GetClientFromToken(ctx context.Context, token string) *github.Client {
	// Authenticate with static token
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Create a new GitHub client
	return github.NewClient(tc)
}

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

type UnreadPullRequest struct {
	Notification *github.Notification
	PR           *github.PullRequest
}

func GetUnreadPullRequests(ctx context.Context, client *github.Client) ([]*UnreadPullRequest, error) {

	// list all notifications
	notificationList, _, err := client.Activity.ListNotifications(ctx, nil)
	if err != nil {
		return nil, err
	}

	// filter to only unread pr notifications
	notifications := make([]*github.Notification, 0)
	for _, notification := range notificationList {
		if *notification.Unread && *notification.Subject.Type == "PullRequest" {
			notifications = append(notifications, notification)
		}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var unreadPullRequests []*UnreadPullRequest
	for _, notification := range notifications {
		wg.Add(1)
		go func(notification *github.Notification) {
			defer wg.Done()
			pr, _, _ := client.PullRequests.Get(ctx, *notification.Repository.Owner.Login, *notification.Repository.Name, GetPullRequestID(notification))
			if *pr.State != "open" {
				return
			}
			mu.Lock()
			unreadPullRequests = append(unreadPullRequests, &UnreadPullRequest{Notification: notification, PR: pr})
			mu.Unlock()
		}(notification)
	}
	wg.Wait()

	sort.Slice(unreadPullRequests, func(i, j int) bool {
		if *unreadPullRequests[i].PR.CreatedAt != *unreadPullRequests[j].PR.CreatedAt {
			return unreadPullRequests[i].PR.CreatedAt.UnixNano() < unreadPullRequests[j].PR.CreatedAt.UnixNano()
		}
		if *unreadPullRequests[i].Notification.Repository.FullName != *unreadPullRequests[j].Notification.Repository.FullName {
			return *unreadPullRequests[i].Notification.Repository.FullName < *unreadPullRequests[j].Notification.Repository.FullName
		}
		return slices.Index(ReasonPriority, *unreadPullRequests[i].Notification.Reason) > slices.Index(ReasonPriority, *unreadPullRequests[j].Notification.Reason)
	})

	return unreadPullRequests, nil
}

func GetPullRequestURL(notification *github.Notification) string {
	return fmt.Sprintf("https://github.com/%s/pull/%d?notification_referrer_id=%s", *notification.Repository.FullName, GetPullRequestID(notification), *notification.ID)
}

func GetPullRequestID(notification *github.Notification) int {
	parts := strings.Split(*notification.Subject.URL, "/")
	val, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return -1
	}
	return val
}
