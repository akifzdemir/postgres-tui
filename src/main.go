package main

import (
	"log"

	"github.com/charmbracelet/bubbletea"
	"go-psql/src/models"
)

func main() {
	model := models.InitialLoginModel()
	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
