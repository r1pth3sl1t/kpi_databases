package search_by_name

import (
	"fmt"
	"rgr/queries/search"
)

type SearcherByName struct {
	search.SearchQueryHandler
}

func (s *SearcherByName) Search() string {
	return `
		select user_id, first_name, last_name, count(*) as connections_num 
		from (select * from "user"
			join "connection"
			on "user".user_id = "connection".u1
			union 
			select * from "user"
			join "connection"
			on "user".user_id = "connection".u2) as users
		where first_name like $1 and last_name like $2
		group by user_id, first_name, last_name
	`
}

func (s *SearcherByName) FetchSearchAttributes() []string {
	return []string{"First name", "Last name"}
}

func (s *SearcherByName) Head() string {
	return fmt.Sprintf("%8s | %15s | %15s | %15s", "user_id", "first_name", "last_name", "connections_num")
}
