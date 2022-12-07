package service

import (
	"context"

	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type IMilestoneTag interface {
	Create(ctx context.Context, create model.MilestoneTag) (*model.MilestoneTag, error)
	FindByPage(ctx context.Context, query *MilestoneTagQuery) (*repo.Paginator[model.MilestoneTag], error)
	DeleteByMilestoneIdAndName(ctx context.Context, id uint, name string) error
}

func NewMilestoneTag(db *gorm.DB) IMilestoneTag {
	return &milestoneTag{db: db}
}

type milestoneTag struct {
	db *gorm.DB
}

func (m milestoneTag) DeleteByMilestoneIdAndName(ctx context.Context, id uint, name string) error {
	r := m.db.WithContext(ctx).Where("milestone_id = ? AND name = ?", id, name).Delete(&model.MilestoneTag{})
	err := r.Error
	if err != nil {
		return err
	}
	if r.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m milestoneTag) Create(ctx context.Context, create model.MilestoneTag) (*model.MilestoneTag, error) {
	err := m.db.WithContext(ctx).
		Create(&create).
		Error
	return &create, err
}

func (m milestoneTag) FindByPage(ctx context.Context, query *MilestoneTagQuery) (*repo.Paginator[model.MilestoneTag], error) {
	db := m.db.Session(&gorm.Session{}).Model(&model.MilestoneTag{}).Distinct("name")
	// build query conditions.
	if query != nil {
		db = db.Where(query.MilestoneTag)
	}
	return repo.FindEntityByPage[model.MilestoneTag](repo.CtxWithDB(ctx, db), query.Page, query.PageSize)
}
