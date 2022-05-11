package main

import (
	"context"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/campbel/revs/app"
	"github.com/campbel/revs/config"
	"github.com/campbel/revs/ghutil"
)

func main() {

	token, err := config.GetToken()
	if err != nil {
		log.Fatal(err)
	}
	if token == "" {
		log.Fatal("error: invalid token")
	}

	ctx := context.Background()

	p := tea.NewProgram(app.NewModel(ctx, ghutil.GetClientFromToken(ctx, token)), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal("error:", err)
	}
}
