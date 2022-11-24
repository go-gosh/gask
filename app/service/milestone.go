package service

import (
	"context"
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/query"
	"github.com/go-gosh/gask/app/repo"
)

type IMilestone interface {
	Create(ctx context.Context, create MilestoneCreate) (*MilestoneView, error)
	FindByPage(ctx context.Context, query *MilestoneQuery) (*repo.Paginator[MilestoneView], error)
	DeleteById(ctx context.Context, id uint, ids ...uint) error
	OneById(ctx context.Context, id uint) (*MilestoneView, error)
	UpdateById(ctx context.Context, id uint, updated *MilestoneUpdate) error
}

func NewMilestoneV2(db *gorm.DB) IMilestone {
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
	return &MilestoneView{Milestone: entity}, m.db.WithContext(ctx).Create(&entity).Error
}

func (m milestone) FindByPage(ctx context.Context, query *MilestoneQuery) (*repo.Paginator[MilestoneView], error) {
	db := m.db.Session(&gorm.Session{}).Model(&model.Milestone{})
	// build query conditions.
	if query != nil {
		db = query.injectDB(db)
	}
	return repo.FindEntityByPage[MilestoneView](repo.CtxWithDB(ctx, db), query.Page, query.PageSize)
}

func (m milestone) DeleteById(ctx context.Context, id uint, ids ...uint) error {
	return repo.WhereInIds(m.db.WithContext(ctx), id, ids...).Delete(&model.Milestone{}).Error
}

func (m milestone) OneById(ctx context.Context, id uint) (*MilestoneView, error) {
	result := MilestoneView{}
	err := m.db.WithContext(ctx).Model(&model.Milestone{}).First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (m milestone) UpdateById(ctx context.Context, id uint, updated *MilestoneUpdate) error {
	return m.db.WithContext(ctx).Model(&model.Milestone{}).Where("`id` = ?", id).Updates(updated.updateDB()).Error
}

func NewMilestone(q *query.Query) *Milestone {
	db, _ := global.GetDatabase()
	return &Milestone{
		q:  q,
		db: db,
	}
}

type Milestone struct {
	q  *query.Query
	db *gorm.DB
}

func (s Milestone) CompleteCheckpointById(id uint, timestamp time.Time) error {
	c, err := s.q.Checkpoint.Where(query.Checkpoint.ID.Eq(id)).First()
	if err != nil {
		return err
	}
	if c.CheckedAt != nil {
		return errors.New("checkpoint already completed")
	}
	return s.q.Transaction(func(tx *query.Query) error {
		r, err := tx.Checkpoint.
			Where(query.Checkpoint.ID.Eq(id), query.Checkpoint.CheckedAt.IsNull()).
			Update(query.Checkpoint.CheckedAt, &timestamp)
		if err != nil {
			return err
		}
		if r.RowsAffected != 1 {
			return errors.New("checkpoint already completed")
		}
		return s.updateMilestoneProgress(tx, c.MilestoneId, c.Point)
	})
}

func (s Milestone) updateMilestoneProgress(tx *query.Query, id uint, point int) error {
	_, err := tx.Milestone.
		Where(query.Milestone.ID.Eq(id)).
		Update(query.Milestone.Progress, query.Milestone.Progress.Add(point))
	return err
}
