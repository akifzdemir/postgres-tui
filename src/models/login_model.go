package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-psql/src/config"
	"log"
	"strings"
)

type model struct {
	inputs     []textinput.Model
	focusIndex int
	connStr    string
}

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	red      = lipgloss.Color("160")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
)

func InitialModel() model {
	m := model{
		inputs:     make([]textinput.Model, 4),
		focusIndex: 0,
	}
	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		switch i {
		case 0:
			t.Placeholder = "Host Name"
			t.Focus()
		case 1:
			t.Placeholder = "Port"
		case 2:
			t.Placeholder = "Username"
		case 3:
			t.Placeholder = "Password"
		}

		m.inputs[i] = t
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyDown:
			m.inputs[m.focusIndex].Blur()
			m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
			m.inputs[m.focusIndex].Focus()
		case tea.KeyUp:
			m.inputs[m.focusIndex].Blur()
			m.focusIndex = (m.focusIndex - 1 + len(m.inputs)) % len(m.inputs)
			m.inputs[m.focusIndex].Focus()
		case tea.KeyEnter:
			var host, username, password, port, connStr string
			if m.focusIndex == len(m.inputs)-1 {
				host = m.inputs[0].Value()
				port = m.inputs[1].Value()
				username = m.inputs[2].Value()
				password = m.inputs[3].Value()
				connStr =
					fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
						username, password, host, port)
				_, err := config.ConnectDb(connStr)
				if err != nil {
					log.Fatal(err)
				}

			} else {
				m.inputs[m.focusIndex].Blur()
				m.focusIndex = (m.focusIndex + 1) % len(m.inputs)
				m.inputs[m.focusIndex].Focus()
			}
		}
	}

	m.updateinputs(msg)
	return m, cmd
}

func (m model) updateinputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m model) View() string {
	builder := strings.Builder{}
	for i := range m.inputs {
		switch i {
		case 0:
			builder.WriteString(inputStyle.Width(10).Render("Host name: "))
		case 1:
			builder.WriteString(inputStyle.Width(10).Render("Port: "))
		case 2:
			builder.WriteString(inputStyle.Width(10).Render("Username: "))
		case 3:
			builder.WriteString(inputStyle.Width(10).Render("Password: "))
		}
		builder.WriteString(m.inputs[i].View() + "\n")
	}
	return builder.String()
}
