package repo

import (
	"context"
)

func FindEntityByPage[V any](ctx context.Context, page, limit int) (*Paginator[V], error) {
	if page < 1 {
		page = DefaultPage
	}
	if limit < 1 {
		limit = DefaultPageSize
	}
	if limit > MaxPageSize {
		limit = MaxPageSize
	}
	var c int64
	result := make([]V, 0, limit)
	err := GetDBFromCtx(ctx).Count(&c).Offset((page - 1) * limit).Limit(limit).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return &Paginator[V]{
		Pager: Pager{
			Page:     page,
			PageSize: limit,
		},
		Total: int(c),
		Data:  result,
	}, nil
}
