package app

import "github.com/google/go-github/github"

type item struct {
	notification *github.Notification
}

func (i item) Title() string       { return *i.notification.Subject.Title }
func (i item) Description() string { return *i.notification.Repository.FullName }
func (i item) FilterValue() string {
	return *i.notification.Subject.Title + " " + *i.notification.Repository.FullName
}
