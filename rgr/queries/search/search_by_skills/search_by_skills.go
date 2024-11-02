package search_by_skills

import (
	"fmt"
	"rgr/queries/search"
)

type SearcherBySkills struct {
	search.SearchQueryHandler
}

func (s *SearcherBySkills) Search() string {
	return `
		select "user".user_id, first_name, last_name, count(*) as skill_count
		from "user"
		join users_to_skills
		on "user".user_id = users_to_skills.user_id
		join skill
		on users_to_skills.skill_id = skill.skill_id
		where skill_type ilike $1
		group by "user".user_id, first_name, last_name	
	`
}

func (s *SearcherBySkills) FetchSearchAttributes() []string {
	return []string{"Skill type"}
}

func (s *SearcherBySkills) Head() string {
	return fmt.Sprintf("%8s | %15s | %15s | %15s", "user_id", "first_name", "last_name", "skill_count")
}
