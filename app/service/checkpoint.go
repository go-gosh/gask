package service

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type ICheckpoint interface {
	Create(ctx context.Context, create *CheckpointCreate) (*CheckpointView, error)
	FindByPage(ctx context.Context, query *CheckpointQuery) (*repo.Paginator[CheckpointView], error)
	DeleteById(ctx context.Context, id uint, ids ...uint) error
	OneById(ctx context.Context, id uint) (*CheckpointView, error)
	UpdateById(ctx context.Context, id uint, updated *CheckpointUpdate) error
}

type checkpoint struct {
	db *gorm.DB
}

func (c checkpoint) Create(ctx context.Context, create *CheckpointCreate) (*CheckpointView, error) {
	if err := global.Validate.Struct(create); err != nil {
		return nil, err
	}
	var entity model.Checkpoint
	if err := copier.Copy(&entity, create); err != nil {
		return nil, err
	}
	m := model.Milestone{Model: model.Model{ID: create.MilestoneId}}
	db := c.db.WithContext(ctx)
	err := db.First(&m).Error
	if err != nil {
		return nil, err
	}
	return &CheckpointView{
			Checkpoint: entity,
		}, c.db.Transaction(func(tx *gorm.DB) error {
			err := tx.Create(&entity).Error
			if err != nil {
				return err
			}
			if entity.CheckedAt != nil {
				return tx.Model(&m).Update("progress", gorm.Expr("`progress`+?", entity.Point)).Error
			}
			return nil
		})
}

func (c checkpoint) FindByPage(ctx context.Context, query *CheckpointQuery) (*repo.Paginator[CheckpointView], error) {
	//TODO implement me
	panic("implement me")
}

func (c checkpoint) DeleteById(ctx context.Context, id uint, ids ...uint) error {
	//TODO implement me
	panic("implement me")
}

func (c checkpoint) OneById(ctx context.Context, id uint) (*CheckpointView, error) {
	//TODO implement me
	panic("implement me")
}

func (c checkpoint) UpdateById(ctx context.Context, id uint, updated *CheckpointUpdate) error {
	//TODO implement me
	panic("implement me")
}

func NewCheckpoint(db *gorm.DB) ICheckpoint {
	return &checkpoint{db: db}
}
