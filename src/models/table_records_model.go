package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-psql/src/constants"
	"log"
)

type TableRecordsModel struct {
	dataTable table.Model
	db        *sql.DB
	help      help.Model
	tableList TableListModel
	columns   []table.Column
}

func InitialTableRecordsModel(tableName string, db *sql.DB, tableList TableListModel) TableRecordsModel {
	query := fmt.Sprintf(`SELECT * FROM %s`, tableName)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	dTable := table.New()
	tableColumns := make([]table.Column, 0)

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	for _, column := range columns {
		tableColumns = append(tableColumns, table.Column{Title: column, Width: 10})
	}

	dTable.SetColumns(tableColumns)

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)
		if err != nil {
			log.Fatal(err)
		}

		row := make([]string, len(columns))
		for i, val := range values {
			if b, ok := val.([]byte); ok {
				row[i] = string(b)
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		dTable.SetRows(append(dTable.Rows(), row))
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57"))

	dTable.SetStyles(s)
	help := help.New()
	return TableRecordsModel{dataTable: dTable, db: db, tableList: tableList, help: help, columns: tableColumns}
}

func (m TableRecordsModel) Init() tea.Cmd {
	return nil
}

func (m TableRecordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.GeneralKeys.Enter):
		case key.Matches(msg, constants.GeneralKeys.Up):
			m.dataTable.SetCursor(m.dataTable.Cursor() - 1)
		case key.Matches(msg, constants.GeneralKeys.Down):
			m.dataTable.SetCursor(m.dataTable.Cursor() + 1)
		case key.Matches(msg, constants.GeneralKeys.Back):
			return m.tableList.Update(tea.KeyMsg{})
		case key.Matches(msg, constants.TableKeys.Create):
		case key.Matches(msg, constants.TableKeys.Delete):

		case key.Matches(msg, constants.GeneralKeys.Quit):
			return m, tea.Quit

		}
	}

	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(msg)
	return m, cmd
}

func (m TableRecordsModel) View() string {
	helpView := m.help.View(constants.TableKeys)
	return lipgloss.
		JoinVertical(
			lipgloss.
				Left, lipgloss.
				NewStyle().
				Padding(1, 2).
				Render(m.dataTable.View()), helpView)
}
