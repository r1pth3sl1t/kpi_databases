package search_by_name

import (
	"fmt"
	"gorm.io/gorm"
	"rgr/entity"
)

type SearcherByName struct {
}
type SearchByNameDTO struct {
	UserId         int    `gorm:"column:user_id"`
	FirstName      string `gorm:"column:first_name"`
	LastName       string `gorm:"column:last_name"`
	ConnectionsNum int    `gorm:"column:connections_num"`
}

func (s *SearcherByName) Search(db *gorm.DB, params map[string]any) []string {
	var users []SearchByNameDTO
	var usersAsStrings []string

	db.Select("user_id, first_name, last_name, count(u1) as connections_num").
		Table("(? UNION ?) as SearchByNameDTO",
			db.Select("*").Model(&entity.User{}).
				Joins("left join \"connection\" on \"user\".user_id = \"connection\".u1"),

			db.Select("*").Model(&entity.User{}).
				Joins("left join \"connection\" on \"user\".user_id = \"connection\".u2")).
		Group("user_id, first_name, last_name").
		Find(&users, params)

	for _, user := range users {
		usersAsStrings = append(usersAsStrings, fmt.Sprintf("%8d | %15s | %15s | %15d", user.UserId, user.FirstName, user.LastName, user.ConnectionsNum))
	}

	return usersAsStrings

}

func (s *SearcherByName) FetchSearchAttributes() []string {
	return []string{"first_name", "last_name"}
}

func (s *SearcherByName) Head() string {
	return fmt.Sprintf("%8s | %15s | %15s | %15s", "user_id", "first_name", "last_name", "connections_num")
}
