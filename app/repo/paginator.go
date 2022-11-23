package repo

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 500
)

type Pager struct {
	Page     int
	PageSize int
}

type Paginator[T any] struct {
	Pager
	Total int
	Data  []T
}
