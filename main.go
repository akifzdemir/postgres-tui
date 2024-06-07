package main

import (
	"go-psql/models"
	"log"

	"github.com/charmbracelet/bubbletea"
)

func main() {
	model := models.InitialLoginModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
