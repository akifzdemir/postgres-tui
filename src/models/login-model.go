package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoginModel struct {
	inputs     []textinput.Model
	focusIndex int
}

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	red      = lipgloss.Color("160")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("160"))
)

func InitialLoginModel() LoginModel {
	m := LoginModel{
		inputs:     make([]textinput.Model, 4),
		focusIndex: 0,
	}
	var t textinput.Model
	for i := range m.inputs {
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
		m.inputs[i] = t
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
				return dbModel.Update(msg)
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
	builder := strings.Builder{}
	for i := range m.inputs {
		builder.WriteString(lipgloss.NewStyle().
			Border(lipgloss.ThickBorder()).
			Width(30).
			Render(m.getInputName(i)+m.inputs[i].View()) + "\n")
	}
	view := builder.String()
	return view
}

func (m LoginModel) getInputName(i int) string {
	builder := strings.Builder{}
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
	return builder.String()
}
