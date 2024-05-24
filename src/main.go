package main

import (
	"log"

	"github.com/charmbracelet/bubbletea"
	"go-psql/src/models"
)

func main() {
	model := models.InitialModel()
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
