package service

import (
	"context"

	"github.com/go-gosh/gask/app/repo"
)

type ICheckpoint interface {
	Create(ctx context.Context, create *CheckpointCreate) (*CheckpointView, error)
	FindByPage(ctx context.Context, query *CheckpointQuery) (*repo.Paginator[CheckpointView], error)
	DeleteById(ctx context.Context, id uint, ids ...uint) error
	OneById(ctx context.Context, id uint) (*CheckpointView, error)
	UpdateById(ctx context.Context, id uint, updated *CheckpointUpdate) error
}
