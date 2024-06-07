package models

import (
	"database/sql"
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"go-psql/constants"
	"log"
	"strings"
	"time"
)

type TableRecordsModel struct {
	dataTable    table.Model
	db           *sql.DB
	help         help.Model
	tableList    TableListModel
	columns      []table.Column
	tableName    string
	errorMessage string
	confirm      *ConfirmModel
	rowToDelete  table.Row
}

type ConfirmModel struct {
	message   string
	confirmed bool
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
		Background(constants.BlueViolet).
		Foreground(constants.White)

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
		case key.Matches(msg, constants.GeneralKeys.Up):
			m.dataTable.SetCursor(m.dataTable.Cursor() - 1)
		case key.Matches(msg, constants.GeneralKeys.Down):
			m.dataTable.SetCursor(m.dataTable.Cursor() + 1)
		case key.Matches(msg, constants.GeneralKeys.Back):
			return m.tableList.Update(tea.KeyMsg{})
		case key.Matches(msg, constants.TableKeys.Delete):
			m.confirm = &ConfirmModel{
				message: fmt.Sprintf("Are you sure you want to delete the selected row? (y/n)"),
			}
			m.rowToDelete = m.dataTable.SelectedRow()
		case key.Matches(msg, constants.TableKeys.Yes):
			if m.confirm != nil {
				m.confirm.confirmed = true
				return m.performDelete(), nil
			}
		case key.Matches(msg, constants.TableKeys.No):
			if m.confirm != nil {
				m.confirm = nil
				return m, nil
			}
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
	errorMessage := ""
	confirmationMessage := ""
	if m.errorMessage != "" {
		errorMessage = lipgloss.NewStyle().Foreground(constants.Red).Render(m.errorMessage)
	}

	if m.confirm != nil {
		confirmationMessage = lipgloss.NewStyle().Foreground(constants.Yellow).Render(m.confirm.message)
	}
	return lipgloss.
		JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Padding(1, 2).Render(m.dataTable.View()),
			helpView,
			confirmationMessage,
			errorMessage,
		)
}

func (m TableRecordsModel) deleteRecord(row table.Row) tea.Model {
	var whereClause strings.Builder
	for i, column := range m.columns {
		if row[i] != "<nil>" {
			value := row[i]
			if _, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", value); err == nil {
				continue
			}
			whereClause.WriteString(fmt.Sprintf("\"%s\" = '%s'", column.Title, value))
			if i < len(m.columns)-1 {
				whereClause.WriteString(" AND ")
			}
		}
	}

	whereClauseString := whereClause.String()
	if strings.HasSuffix(whereClauseString, " AND ") {
		whereClauseString = whereClauseString[:len(whereClauseString)-5]
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", m.tableName, whereClauseString)
	_, err := m.db.Exec(query)
	if err != nil {
		m.errorMessage = err.Error()
		return m
	}
	return m.refreshTable()
}

func (m TableRecordsModel) performDelete() tea.Model {
	if m.confirm != nil && m.confirm.confirmed {
		return m.deleteRecord(m.rowToDelete)
	}
	return m
}

func (m TableRecordsModel) refreshTable() tea.Model {
	return InitialTableRecordsModel(m.tableName, m.db, m.tableList)
}
