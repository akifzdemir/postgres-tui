package models

import (
	"database/sql"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-psql/src/config"
	"go-psql/src/constants"
	"log"
)

type DatabaseModel struct {
	db      *sql.DB
	dbList  list.Model
	connStr string
}

type item struct {
	title, desc string
}

var appStyle = lipgloss.NewStyle().Padding(1, 2)

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func InitialDatabaseModel(connStr string) DatabaseModel {
	db, err := config.ConnectDb(connStr)
	if err != nil {
		log.Fatal(err)
	}
	items := make([]list.Item, 0)
	rows, err := db.Query("SELECT datname FROM pg_database WHERE datistemplate = false;")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		items = append(items, item{title: name, desc: ""})
	}
	dbList := list.New(items, list.NewDefaultDelegate(), 30, 20)
	dbList.Title = "Databases"
	dbList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.GeneralKeys.Enter,
		}
	}
	return DatabaseModel{db: db, dbList: dbList, connStr: connStr}
}

func (m DatabaseModel) Init() tea.Cmd {
	return nil
}

func (m DatabaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.GeneralKeys.Enter):
			var selectedItem item
			i, ok := m.dbList.SelectedItem().(item)
			if ok {
				selectedItem = i
				config.RemoveDb()
				tablesModel := InitialTableListModel(m.connStr, selectedItem.title)
				return tablesModel.Update(tea.KeyMsg{})
			}
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		}

	}
	var cmd tea.Cmd
	m.dbList, cmd = m.dbList.Update(msg)
	return m, cmd
}

func (m DatabaseModel) View() string {
	return appStyle.Render(m.dbList.View())
}
