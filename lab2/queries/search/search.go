package search

import "gorm.io/gorm"

type SearchQueryHandler interface {
	Search(db *gorm.DB, params map[string]any) []string
	FetchSearchAttributes() []string
	Head() string
}
