package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"go-psql/src/config"
	"go-psql/src/constants"
	"log"
)

type TablesModel struct {
	dataTable table.Model
	tableList list.Model
	item
	db      *sql.DB
	connStr string
}

func InitialTablesModel(connStr string, dbName string) TablesModel {

	newConnStr := fmt.Sprintf("%s dbname=%s", connStr, dbName)
	db, err := config.ConnectDb(newConnStr)
	if err != nil {
		log.Fatal("test-->", err)
	}

	query := fmt.Sprintf(`SELECT tablename 
						  FROM pg_catalog.pg_tables 
						  WHERE schemaname != 'pg_catalog'
						  AND schemaname != 'information_schema';`)

	rows, queryErr := db.Query(query)
	if queryErr != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	items := make([]list.Item, 0)
	for rows.Next() {
		var tbName string
		if err := rows.Scan(&tbName); err != nil {
			log.Fatal(err)
		}
		items = append(items, item{title: tbName, desc: ""})
	}
	tblist := list.New(items, list.NewDefaultDelegate(), 30, 20)
	tblist.Title = "Tables"
	tblist.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Enter,
			constants.Keymap.Back,
		}
	}
	return TablesModel{db: db, tableList: tblist, connStr: connStr}
}

func (m TablesModel) Init() tea.Cmd {
	return nil
}

func (m TablesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, constants.Keymap.Back):
			dbListModel := InitialDatabaseModel(m.connStr)
			return dbListModel.Update(tea.KeyMsg{})
		}

	}
	var cmd tea.Cmd
	m.tableList, cmd = m.tableList.Update(msg)
	return m, cmd
}

func (m TablesModel) View() string {
	return appStyle.Render(m.tableList.View())
}
