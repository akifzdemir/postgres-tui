package models

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
	constants2 "go-psql/constants"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	inputs     []textinput.Model
	focusIndex int
	help       help.Model
}

func InitialLoginModel() LoginModel {
	inputs := make([]textinput.Model, 4)
	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		switch i {
		case 0:
			t.Placeholder = "localhost"
			t.Focus()
		case 1:
			t.Placeholder = "5432"
		case 2:
			t.Placeholder = "postgres"
		case 3:
			t.Placeholder = "1234"
		}
		inputs[i] = t
	}

	helpModel := help.New()

	m := LoginModel{
		inputs:     inputs,
		focusIndex: 0,
		help:       helpModel,
	}
	return m
}

func (m LoginModel) Init() tea.Cmd {
	return nil
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.focusIndex == len(m.inputs)-1 {
				connStr := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
					m.inputs[2].Value(), m.inputs[3].Value(), m.inputs[0].Value(), m.inputs[1].Value())

				dbModel := InitialDatabaseModel(connStr)
				return dbModel.Update(tea.KeyMsg{})
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

func (m LoginModel) updateinputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m LoginModel) View() string {
	header := lipgloss.
		NewStyle().
		SetString("PostgreSQL TUI").
		Align(lipgloss.Center).
		Width(20).
		PaddingBottom(2).
		Foreground(constants2.BlueViolet).String()

	var views []string
	views = append(views, header)

	for i := range m.inputs {
		inputView := lipgloss.JoinVertical(
			lipgloss.Left,
			m.getInputName(i),
			m.inputs[i].View(),
		)
		views = append(views, inputView)
	}
	helpView := m.help.View(constants2.LoginKeys)
	views = append(views, helpView)
	finalView := lipgloss.JoinVertical(lipgloss.Center, views...)

	return constants2.BorderStyle.Render(finalView)
}

func (m LoginModel) getInputName(i int) string {
	builder := strings.Builder{}
	switch i {
	case 0:
		builder.WriteString(constants2.InputStyle.Render("Host name: "))
	case 1:
		builder.WriteString(constants2.InputStyle.Render("Port: "))
	case 2:
		builder.WriteString(constants2.InputStyle.Render("Username: "))
	case 3:
		builder.WriteString(constants2.InputStyle.Render("Password: "))
	}
	return builder.String()
}
