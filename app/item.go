package app

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/go-github/github"
	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(rate.Limit(1), 1)

type item struct {
	notification *github.Notification
	pr           *github.PullRequest
}

func (i *item) Title() string { return *i.notification.Subject.Title }
func (i *item) Description() string {
	if i.pr != nil {
		return fmt.Sprintf("%s • %s • created %s",
			*i.notification.Repository.FullName,
			*i.pr.User.Login,
			humanize.Time(*i.pr.CreatedAt),
		)
	}
	return fmt.Sprintf("%s", *i.notification.Repository.FullName)
}
func (i *item) FilterValue() string {
	return fmt.Sprint(
		*i.notification.Subject.Title,
		*i.notification.Repository.FullName,
		*i.pr.User.Login,
	)
}
