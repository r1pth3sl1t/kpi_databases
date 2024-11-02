package search

type SearchQueryHandler interface {
	Search() string
	FetchSearchAttributes() []string
	Head() string
}
