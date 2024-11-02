package search_by_connections_num_range

import (
	"fmt"
	"rgr/queries/search"
)

type SearcherByConnectionsNumRange struct {
	search.SearchQueryHandler
}

func (s *SearcherByConnectionsNumRange) Search() string {
	return `
		select user_id, first_name, last_name, count(u1) as connections_num 
        from (select * from "user"
            left join "connection"
            on "user".user_id = "connection".u1
            union 
            select * from "user"
            left join "connection"
            on "user".user_id = "connection".u2) as users
        group by user_id, first_name, last_name
        having count(u1) between $1 and $2
	`
}

func (s *SearcherByConnectionsNumRange) FetchSearchAttributes() []string {
	return []string{"min", "max"}
}

func (s *SearcherByConnectionsNumRange) Head() string {
	return fmt.Sprintf("%8s | %15s | %15s | %15s", "user_id", "first_name", "last_name", "connections_num")
}
