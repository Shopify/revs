package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Shopify/revs/ghutil"
)

const ()

func GetToken() (string, error) {

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	revsPath := dir + "/.config/revs"
	toknPath := revsPath + "/token"

	data, err := ioutil.ReadFile(toknPath)
	if err == nil {
		return string(data), nil
	}

	token, err := promptForToken()
	if err != nil {
		return "", err
	}
	if token == "" {
		return "", fmt.Errorf("no token provided")
	}

	if err := ghutil.ValidateToken(context.TODO(), token); err != nil {
		return "", err
	}

	if err := os.MkdirAll(revsPath, 0700); err != nil {
		return "", err
	}

	fmt.Println("Saving token to", toknPath)
	if err := ioutil.WriteFile(toknPath, []byte(token), 0600); err != nil {
		return "", err
	}

	return token, nil
}
