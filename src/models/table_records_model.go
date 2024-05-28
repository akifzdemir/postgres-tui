package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-psql/src/constants"
	"log"
	"strings"
)

type TableRecordsModel struct {
	dataTable table.Model
	db        *sql.DB
	tableList TableListModel
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
	return TableRecordsModel{dataTable: dTable, db: db, tableList: tableList}
}

func (m TableRecordsModel) Init() tea.Cmd {
	return nil
}

func (m TableRecordsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Enter):
		case key.Matches(msg, constants.Keymap.Up):
			m.dataTable.SetCursor(m.dataTable.Cursor() - 1)
		case key.Matches(msg, constants.Keymap.Down):
			m.dataTable.SetCursor(m.dataTable.Cursor() + 1)
		case key.Matches(msg, constants.Keymap.Back):
			return m.tableList.Update(tea.KeyMsg{})

		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit

		}
	}

	var cmd tea.Cmd
	m.dataTable, cmd = m.dataTable.Update(msg)
	return m, cmd
}

func (m TableRecordsModel) View() string {
	builder := strings.Builder{}

	builder.WriteString(appStyle.Render(m.dataTable.View()) + "\n")
	builder.WriteString(m.renderHelp())

	return builder.String()
}

func (m TableRecordsModel) renderHelp() string {
	helpView := ""
	helpView += constants.Keymap.Up.Help().Desc + " " + constants.Keymap.Up.Help().Key + " | "
	helpView += constants.Keymap.Down.Help().Desc + " " + constants.Keymap.Down.Help().Key + " | "
	helpView += constants.Keymap.Enter.Help().Desc + " " + constants.Keymap.Enter.Help().Key + " | "
	helpView += constants.Keymap.Quit.Help().Desc + " " + constants.Keymap.Quit.Help().Key

	return lipgloss.NewStyle().Padding(1, 2).Foreground(darkGray).Render(helpView)
}
