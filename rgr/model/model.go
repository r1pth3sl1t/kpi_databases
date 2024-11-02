package model

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"rgr/queries"
)

type DatabaseCredentials struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type Model struct {
	db *sql.DB
}

func New() (*Model, error) {
	config, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	var dbc DatabaseCredentials
	decoder := json.NewDecoder(config)
	err = decoder.Decode(&dbc)

	if err != nil {
		return nil, err
	}

	m := new(Model)
	connStr := fmt.Sprintf("user=%s host=%s port=%s dbname=%s password=%s sslmode=disable",
		dbc.User, dbc.Host, dbc.Port, dbc.Dbname, dbc.Password)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	m.db = db

	return m, err
}

func (m *Model) Close() {
	m.db.Close()
}

func (m *Model) FetchTableData() map[string][]string {
	rows, err := m.db.Query(queries.GetFetchingTablesDataQuery())
	if err != nil {
		fmt.Println(err)
	}
	tables := make(map[string][]string)
	for rows.Next() {
		var tableName, column string
		rows.Scan(&tableName, &column)
		tables[tableName] = append(tables[tableName], column)
	}
	return tables
}

func (m *Model) Insert(table string, data map[string]string) error {

	query, values := queries.PrepareInsertQuery(table, data)

	_, err := m.db.Query(query, values...)

	return err
}

func (m *Model) FetchTablePrimaryKeys() map[string][]string {
	rows, err := m.db.Query(queries.GetFetchingPrimaryKeysQuery())
	if err != nil {
		fmt.Println(err)
	}
	tables := make(map[string][]string)

	for rows.Next() {

		var tableName, column string
		rows.Scan(&tableName, &column, nil)

		tables[tableName] = append(tables[tableName], column)
	}
	return tables
}

func (m *Model) Update(tableName string, data map[string]string, pkey map[string]string) error {
	query, values := queries.PrepareUpdateQuery(tableName, data, pkey)

	_, err := m.db.Query(query, values...)

	return err
}

func (m *Model) Delete(table string, pkey map[string]string) error {
	query, values := queries.PrepareDeleteQuery(table, pkey)
	fmt.Println(query)
	_, err := m.db.Query(query, values...)
	return err
}

func (m *Model) GenerateDataSet(size int) error {
	_, err := m.db.Query(queries.GetUserGeneratingQuery(), size)
	if err != nil {
		return err
	}
	_, err = m.db.Query(queries.GenerateSkillsQuery(), size)
	if err != nil {
		return err
	}
	_, err = m.db.Query(queries.GenerateConnectionQuery(), size)
	if err != nil {
		return err
	}
	return err
}

func (m *Model) Search(query string, data map[string]string, attrOrder []string) ([]string, error) {
	var args []any
	for _, attr := range attrOrder {
		args = append(args, data[attr])
	}

	rows, err := m.db.Query(query, args...)
	var rowsToReturn []string
	userId := 0
	connectionsNum := 0
	var firstName, lastName string
	var row string
	for rows.Next() {
		err := rows.Scan(&userId, &firstName, &lastName, &connectionsNum)
		if err != nil {
			return rowsToReturn, err
		}
		row = fmt.Sprintf("%8d | %15s | %15s | %15d", userId, firstName, lastName, connectionsNum)
		rowsToReturn = append(rowsToReturn, row)
	}

	return rowsToReturn, err
}
