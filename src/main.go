package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

type model struct {
	textInput textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Database"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	return model{textInput: ti}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}
func (m model) View() string {
	builder := strings.Builder{}
	builder.WriteString(m.textInput.View())
	return builder.String()
}

func main() {

	//connStr := "user=postgres password=1234 host=localhost port=5432 sslmode=disable"
	//db, err := config.ConnectDb(connStr)
	//
	//rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//
	//defer rows.Close()
	//
	//for rows.Next() {
	//	var name string
	//	if err := rows.Scan(&name); err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(name)
	//}

	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
