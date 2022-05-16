package config

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

const (
	promptMessage = `
Add a Github Personal Access Token (PAT)
https://github.com/settings/tokens/new?description=revs-cli-token&scopes=repo,notifications

This is used to access your notifications and information about pull requests you've been asked to review.

> `
)

func promptForToken() (string, error) {
	fmt.Print(promptMessage)

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}
