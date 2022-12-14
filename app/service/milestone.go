package service

import (
	"context"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type IMilestone interface {
	Create(ctx context.Context, create MilestoneCreate) (*MilestoneView, error)
	FindByPage(ctx context.Context, query *MilestoneQuery) (*repo.Paginator[MilestoneView], error)
	DeleteById(ctx context.Context, id uint, ids ...uint) error
	OneById(ctx context.Context, id uint) (*MilestoneView, error)
	UpdateById(ctx context.Context, id uint, updated *MilestoneUpdate) error
}

func NewMilestone(db *gorm.DB) IMilestone {
	return &milestone{db: db}
}

type milestone struct {
	db *gorm.DB
}

func (m milestone) Create(ctx context.Context, create MilestoneCreate) (*MilestoneView, error) {
	if err := global.Validate.Struct(create); err != nil {
		return nil, err
	}
	var entity model.Milestone
	if err := copier.Copy(&entity, &create); err != nil {
		return nil, err
	}
	err := m.db.WithContext(ctx).Create(&entity).Error
	if err != nil {
		return nil, err
	}
	var result MilestoneView
	err = copier.Copy(&result, &entity)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m milestone) FindByPage(ctx context.Context, query *MilestoneQuery) (*repo.Paginator[MilestoneView], error) {
	db := m.db.Session(&gorm.Session{}).Model(&model.Milestone{}).Preload("Tags")
	// build query conditions.
	if query != nil {
		db = query.injectDB(db)
	}
	return repo.FindEntityByPage[MilestoneView](repo.CtxWithDB(ctx, db), query.Page, query.PageSize)
}

func (m milestone) DeleteById(ctx context.Context, id uint, ids ...uint) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		r := repo.WhereInIds(tx, id, ids...).Delete(&model.Milestone{})
		err := r.Error
		if err != nil {
			return err
		}
		if r.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		if len(ids) != 0 {
			return tx.Where("`milestone_tags`.`milestone_id` in (?)", append(ids, id)).Delete(&model.MilestoneTag{}).Error
		}
		return tx.Where("`milestone_tags`.`milestone_id` = ?", id).Delete(&model.MilestoneTag{}).Error
	})
}

func (m milestone) OneById(ctx context.Context, id uint) (*MilestoneView, error) {
	result := MilestoneView{}
	err := m.db.WithContext(ctx).Model(&model.Milestone{}).Preload("Tags").First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m milestone) UpdateById(ctx context.Context, id uint, updated *MilestoneUpdate) error {
	return m.db.WithContext(ctx).Table("milestones").Select("*").Where("`id` = ?", id).Updates(updated).Error
}
