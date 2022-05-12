package config

import (
	"fmt"
	"io/ioutil"
	"os"

	tea "github.com/charmbracelet/bubbletea"
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

	p := tea.NewProgram(initialModel())
	teaModelOut, _ := p.StartReturningModel()
	modelOut := teaModelOut.(model)
	token := modelOut.textInput.Value()
	if token == "" {
		return "", fmt.Errorf("no token provided")
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
