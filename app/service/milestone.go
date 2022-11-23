package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gen/field"
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
	//TODO implement me
	panic("implement me")
}

func (m milestone) OneById(ctx context.Context, id uint) (*MilestoneView, error) {
	//TODO implement me
	panic("implement me")
}

func (m milestone) UpdateById(ctx context.Context, id uint, updated *MilestoneUpdate) error {
	//TODO implement me
	panic("implement me")
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

func (s Milestone) SplitMilestoneById(id uint, first CheckpointCreate, checkpoints ...CheckpointCreate) ([]*model.Checkpoint, error) {
	m, err := s.q.Milestone.Where(s.q.Milestone.ID.Eq(id)).First()
	if err != nil {
		return nil, err
	}

	if err := global.Validate.Struct(first); err != nil {
		return nil, err
	}

	var point int
	if first.CheckedAt != nil {
		point += first.Point
	}
	inputs := make([]CheckpointCreate, 0, len(checkpoints)+1)
	inputs = append(inputs, first)
	for i := range checkpoints {
		if err := global.Validate.Struct(checkpoints[i]); err != nil {
			return nil, err
		}
		inputs = append(inputs, checkpoints[i])
		if checkpoints[i].CheckedAt != nil {
			point += checkpoints[i].Point
		}
	}

	results := make([]*model.Checkpoint, 0, 1+len(checkpoints))
	if err := copier.Copy(&results, &inputs); err != nil {
		return nil, err
	}
	err = s.q.Transaction(func(tx *query.Query) error {
		if err := tx.Milestone.Checkpoints.Model(m).Append(results...); err != nil {
			return err
		}
		if point > 0 {
			return s.updateMilestoneProgress(tx, m.ID, point)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return results, err
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

func (s Milestone) relatedCheckpointOrder() field.RelationField {
	return s.q.Milestone.Checkpoints.Order(
		s.q.Checkpoint.CheckedAt.IsNull(),
		s.q.Checkpoint.UpdatedAt.Desc(),
	)
}

func (s Milestone) DeleteById(id uint) error {
	_, err := s.q.Milestone.Where(s.q.Milestone.ID.Eq(id)).Delete()
	return err
}

func (s Milestone) RetrieveById(id uint) (*model.Milestone, error) {
	return s.q.Milestone.Where(s.q.Milestone.ID.Eq(id)).Preload(s.relatedCheckpointOrder()).First()
}

func (s Milestone) uncheckedCheckpoints(t time.Time) query.ICheckpointDo {
	ck := s.q.Checkpoint
	return ck.Where(ck.JoinedAt.Gte(t)).Order(ck.CheckedAt.IsNotNull(), field.NewField("", fmt.Sprintf("DATE(joined_at)-DATE('%s')", t.Format("2006-01-02 15:04:05"))), ck.UpdatedAt.Desc())
}

func (s Milestone) checkedCheckpoints(t time.Time) query.ICheckpointDo {
	ck := s.q.Checkpoint
	return ck.Where(ck.CheckedAt.IsNotNull(), ck.JoinedAt.Lt(t)).Order(ck.CheckedAt, ck.UpdatedAt.Desc())
}

func (s Milestone) PaginateCheckpoints(page int, limit int, q CheckpointQuery) ([]*CheckpointView, int64, error) {
	offset := 0
	if page > 1 {
		offset = limit * (page - 1)
	}
	result := make([]*CheckpointView, 0, limit)
	var count int64
	db := s.db.Model(&model.Checkpoint{}).Preload("Milestone").Select("*, julianday(joined_at) - julianday(?) as diff", q.Timestamp).Order("`checked_at` is not null, `checked_at` desc, `diff`, `updated_at` desc").WithContext(context.Background())
	if q.MilestoneId != 0 {
		db = db.Where("milestone_id = ?", q.MilestoneId)
	}
	err := db.Count(&count).Offset(offset).Limit(limit).Find(&result).Error
	return result, count, err
}
