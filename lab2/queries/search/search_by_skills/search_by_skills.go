package search_by_skills

import (
	"fmt"
	"gorm.io/gorm"
	"rgr/entity"
)

type SearcherBySkills struct {
}

func (s *SearcherBySkills) Search(db *gorm.DB, params map[string]any) []string {
	var users []entity.User
	var usersAsStrings []string
	db.Model(&entity.User{}).Preload("Skills", params).Find(&users)

	for _, user := range users {
		if len(user.Skills) > 0 {
			usersAsStrings = append(usersAsStrings, fmt.Sprintf("%8d | %15s | %15s | %15d", user.UserId, user.FirstName, user.LastName, len(user.Skills)))
		}
	}
	return usersAsStrings
}

func (s *SearcherBySkills) FetchSearchAttributes() []string {
	return []string{"skill_type"}
}

func (s *SearcherBySkills) Head() string {
	return fmt.Sprintf("%8s | %15s | %15s | %15s", "user_id", "first_name", "last_name", "skill_count")
}
