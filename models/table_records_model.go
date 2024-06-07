package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	constants2 "go-psql/constants"
	"log"
	"strings"
)

type TableRecordsModel struct {
	dataTable table.Model
	db        *sql.DB
	help      help.Model
	tableList TableListModel
	columns   []table.Column
	tableName string
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
		BorderBottom(true)

	s.Selected = s.Selected.
		Background(constants2.BlueViolet).
		Foreground(constants2.White)

	dTable.SetStyles(s)
	help := help.New()
	return TableRecordsModel{
		dataTable: dTable,
		db:        db,
		tableList: tableList,
		help:      help,
		columns:   tableColumns,
		tableName: tableName,
	}
}

func (m TableRecordsModel) Init() tea.Cmd {
	return nil
}

func (m TableRecordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants2.GeneralKeys.Enter):
		case key.Matches(msg, constants2.GeneralKeys.Up):
			m.dataTable.SetCursor(m.dataTable.Cursor() - 1)
		case key.Matches(msg, constants2.GeneralKeys.Down):
			m.dataTable.SetCursor(m.dataTable.Cursor() + 1)
		case key.Matches(msg, constants2.GeneralKeys.Back):
			return m.tableList.Update(tea.KeyMsg{})
		case key.Matches(msg, constants2.TableKeys.Delete):
			return m.deleteRecord(m.dataTable.SelectedRow()), nil
		case key.Matches(msg, constants2.GeneralKeys.Quit):
			return m, tea.Quit

		}
	}

	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(msg)
	return m, cmd
}

func (m TableRecordsModel) View() string {
	helpView := m.help.View(constants2.TableKeys)
	return lipgloss.
		JoinVertical(
			lipgloss.
				Left, lipgloss.
				NewStyle().
				Padding(1, 2).
				Render(m.dataTable.View()), helpView)
}

func (m TableRecordsModel) deleteRecord(row table.Row) tea.Model {
	var whereClause strings.Builder
	for i, column := range m.columns {
		if row[i] != "<nil>" {
			whereClause.WriteString(fmt.Sprintf("%s = '%s'", column.Title, row[i]))
			if i < len(m.columns)-1 {
				whereClause.WriteString(" AND ")
			}
		}

	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s", m.tableName, whereClause.String())
	_, err := m.db.Exec(query)
	if err != nil {
		log.Println(err)
	}
	return m.refreshTable()
}

func (m TableRecordsModel) refreshTable() tea.Model {
	return InitialTableRecordsModel(m.tableName, m.db, m.tableList)
}
