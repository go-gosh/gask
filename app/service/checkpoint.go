package service

import (
	"context"
	"log"

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
	err = c.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&entity).Error
		if err != nil {
			return err
		}
		if entity.CheckedAt != nil {
			return c.callbackMilestone(tx, entity.MilestoneId)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	v := CheckpointView{}
	err = copier.Copy(&v, &entity)
	return &v, err
}

func (c checkpoint) FindByPage(ctx context.Context, query *CheckpointQuery) (*repo.Paginator[CheckpointView], error) {
	db := c.db.Session(&gorm.Session{}).Model(&model.Checkpoint{})
	// build query conditions.
	if query != nil {
		db = query.injectDB(db)
	}
	return repo.FindEntityByPage[CheckpointView](repo.CtxWithDB(ctx, db), query.Page, query.PageSize)
}

func (c checkpoint) DeleteById(ctx context.Context, id uint, ids ...uint) error {
	old, err := c.OneById(ctx, id)
	if err != nil {
		return err
	}
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := repo.WhereInIds(tx, id, ids...).Delete(&model.Checkpoint{}).Error
		if err != nil {
			return err
		}
		return c.callbackMilestone(tx, old.MilestoneId)
	})
}

func (c checkpoint) OneById(ctx context.Context, id uint) (*CheckpointView, error) {
	result := CheckpointView{}
	err := c.db.WithContext(ctx).Model(&model.Checkpoint{}).Omit("diff").First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c checkpoint) UpdateById(ctx context.Context, id uint, updated *CheckpointUpdate) error {
	old, err := c.OneById(ctx, id)
	if err != nil {
		return err
	}
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		log.Printf("%#v", *updated)
		err := tx.Table("checkpoints").Where("`id` = ?", id).Select("*").Updates(updated).Error
		if err != nil {
			return err
		}
		return c.callbackMilestone(tx, old.MilestoneId)
	})
}

func (c checkpoint) callbackMilestone(tx *gorm.DB, milestoneId uint) error {
	sub := tx.Model(&model.Checkpoint{}).Select("SUM(`checkpoints`.`point`) as point").Where("`checkpoints`.`milestone_id` = ? AND `checkpoints`.`checked_at` IS NOT NULL", milestoneId)
	return tx.Model(model.Milestone{}).Where("id = ?", milestoneId).Update("progress", gorm.Expr("IFNULL((?),0)", sub)).Error
}

func NewCheckpoint(db *gorm.DB) ICheckpoint {
	return &checkpoint{db: db}
}
