package search_by_connections_num_range

import (
	"fmt"
	"gorm.io/gorm"
	"rgr/entity"
	"rgr/queries/search/search_by_name"
)

type SearcherByConnectionsNumRange struct {
}

func (s *SearcherByConnectionsNumRange) Search(db *gorm.DB, params map[string]any) []string {
	var users []search_by_name.SearchByNameDTO
	var usersAsStrings []string

	db.Select("user_id, first_name, last_name, count(u1) as connections_num").
		Table("(? UNION ?) as SearchByNameDTO",
			db.Select("*").Model(&entity.User{}).
				Joins("left join \"connection\" on \"user\".user_id = \"connection\".u1"),

			db.Select("*").Model(&entity.User{}).
				Joins("left join \"connection\" on \"user\".user_id = \"connection\".u2")).
		Group("user_id, first_name, last_name").
		Having("count(u1) between ? and ?", params["min"], params["max"]).
		Find(&users)

	for _, user := range users {
		usersAsStrings = append(usersAsStrings, fmt.Sprintf("%8d | %15s | %15s | %15d", user.UserId, user.FirstName, user.LastName, user.ConnectionsNum))
	}
	return usersAsStrings
}

func (s *SearcherByConnectionsNumRange) FetchSearchAttributes() []string {
	return []string{"min", "max"}
}

func (s *SearcherByConnectionsNumRange) Head() string {
	return fmt.Sprintf("%8s | %15s | %15s | %15s", "user_id", "first_name", "last_name", "connections_num")
}
