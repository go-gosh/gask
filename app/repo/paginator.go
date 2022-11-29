package repo

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 500
)

type Pager struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type Paginator[T any] struct {
	Pager
	Total int `json:"total"`
	Data  []T `json:"data"`
}
