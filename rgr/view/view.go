package view

import (
	"errors"
	"fmt"
	"rgr/queries/search"
	"rgr/queries/search/search_by_connections_num_range"
	"rgr/queries/search/search_by_name"
	"rgr/queries/search/search_by_skills"
)

type View struct {
	Tables      map[string][]string
	primaryKeys map[string][]string
	modes       []search.SearchQueryHandler
}

func New(tables map[string][]string, primaryKeys map[string][]string) *View {
	view := &View{}

	view.Tables = tables
	view.primaryKeys = primaryKeys
	view.modes = []search.SearchQueryHandler{nil, &search_by_name.SearcherByName{},
		&search_by_connections_num_range.SearcherByConnectionsNumRange{},
		&search_by_skills.SearcherBySkills{}}
	return view
}

func (v *View) Index(optionsNum int) (int, error) {
	fmt.Println("Main page")
	fmt.Println("---")
	fmt.Println("1. Add new data")
	fmt.Println("2. Edit existing data")
	fmt.Println("3. Delete data")
	fmt.Println("4. Generate dataset")
	fmt.Println("5. Search data")

	var choice int
	_, err := fmt.Scan(&choice)
	if err != nil {
		return 0, err
	}
	if choice >= optionsNum {
		return 0, errors.New("invalid option")
	}
	return choice, nil
}

func (v *View) Error(err error) {
	fmt.Print("There's an error! Error message: ")
	fmt.Println(err)
}

func (v *View) SelectTable() (string, error) {
	fmt.Println("Select table:")
	for tableName := range v.Tables {
		fmt.Println(tableName)
	}
	fmt.Print("Table: ")
	var table string
	_, err := fmt.Scan(&table)
	if err != nil {
		return "", err
	}
	_, ok := v.Tables[table]
	if !ok {
		return "", errors.New("no table found")
	}
	return table, nil
}

func (v *View) FetchAttributes(attributes []string) map[string]string {
	data := make(map[string]string)
	for _, column := range attributes {
		var tmp string
		fmt.Print(column + ": ")
		_, err := fmt.Scan(&tmp)
		if err != nil {
			return nil
		}
		data[column] = tmp
	}
	return data
}

func (v *View) FetchPrimaryKey(tableName string) map[string]string {
	fmt.Println("Enter primary key: ")
	data := make(map[string]string)
	for _, column := range v.primaryKeys[tableName] {
		var tmp string
		fmt.Print(column + ": ")
		_, err := fmt.Scan(&tmp)
		if err != nil {
			return nil
		}
		data[column] = tmp
	}
	return data
}

func (v *View) SelectAttributes(tableName string) []string {
	fmt.Println("Select columns or type '-' to end selecting: ")
	var columns []string
	var column string
	for _, column := range v.Tables[tableName] {
		fmt.Println(column)
	}
	fmt.Println("Column: ")
	for true {
		_, err := fmt.Scan(&column)
		if err != nil {
			return columns
		}
		if column == "-" {
			break
		}

		columns = append(columns, column)
	}

	return columns
}

func (v *View) GetDataSize() int {
	fmt.Println("Type desired size of generated data set:")
	value := 0
	_, err := fmt.Scan(&value)
	if err != nil {
		v.Error(err)
	}
	return value
}

func (v *View) GetSearchingMode() (search.SearchQueryHandler, error) {
	fmt.Println("Select searching mode: ")
	fmt.Println("1. Search users by first name and last name")
	fmt.Println("2. Search users by connections number range")
	fmt.Println("3. Search users and skills amount by specific skill type")
	value := 0
	_, err := fmt.Scan(&value)
	if err != nil {
		v.Error(err)
	}
	if value < 1 || value >= len(v.modes) {
		return nil, errors.New("invalid option")
	}

	return v.modes[value], nil
}

func (v *View) IndexColumns(columns []string, head string) {
	fmt.Println("Columns: ")
	fmt.Println(head)
	for _, column := range columns {
		fmt.Println(column)
	}
}

func (v *View) Success(msg string) {
	fmt.Println(msg)
}
