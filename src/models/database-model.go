package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-psql/src/config"
	"log"
)

type DatabaseModel struct {
	db     *sql.DB
	dbList list.Model
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
	defer db.Close()
	items := make([]list.Item, 0)
	rows, err := db.Query(
		`SELECT d.datname,u.usename 
				FROM pg_user u
				JOIN pg_database d
				ON d.datdba = u.usesysid 
				WHERE d.datistemplate = false;`)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var dba string
		if err := rows.Scan(&name, &dba); err != nil {
			log.Fatal(err)
		}
		items = append(items, item{title: name, desc: dba})
	}
	dbList := list.New(items, list.NewDefaultDelegate(), 20, 20)
	return DatabaseModel{db: db, dbList: dbList}
}

func (m DatabaseModel) Init() tea.Cmd {
	return nil
}

func (m DatabaseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.dbList.SetSize(msg.Width-h, msg.Height-v)
		fmt.Println(msg.Width, msg.Height)
	case tea.KeyMsg:
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
