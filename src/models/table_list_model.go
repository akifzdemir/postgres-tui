package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"go-psql/src/config"
	"go-psql/src/constants"
	"log"
)

type TableListModel struct {
	tableList list.Model
	item
	db      *sql.DB
	connStr string
}

func InitialTableListModel(connStr string, dbName string) TableListModel {

	newConnStr := fmt.Sprintf("%s dbname=%s", connStr, dbName)
	db, err := config.ConnectDb(newConnStr)
	if err != nil {
		log.Fatal(err)
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
			constants.GeneralKeys.Enter,
			constants.GeneralKeys.Back,
		}
	}
	return TableListModel{db: db, tableList: tblist, connStr: connStr}
}

func (m TableListModel) Init() tea.Cmd {
	return nil
}

func (m TableListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, constants.GeneralKeys.Back):
			dbListModel := InitialDatabaseModel(m.connStr)
			return dbListModel.Update(tea.KeyMsg{})
		case key.Matches(msg, constants.GeneralKeys.Enter):
			var selectedItem item
			i, ok := m.tableList.SelectedItem().(item)
			if ok {
				selectedItem = i
				tableRecordList := InitialTableRecordsModel(selectedItem.title, m.db, m)
				return tableRecordList.Update(tea.KeyMsg{})
			}
		}
	}
	var cmd tea.Cmd
	m.tableList, cmd = m.tableList.Update(msg)
	return m, cmd
}

func (m TableListModel) View() string {
	return appStyle.Render(m.tableList.View())
}
