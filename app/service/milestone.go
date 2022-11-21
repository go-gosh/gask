package service

import (
	"errors"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gen/field"

	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/query"
)

func NewMilestone(q *query.Query) *Milestone {
	return &Milestone{
		q: q,
	}
}

type Milestone struct {
	q *query.Query
}

func (s Milestone) CreateMilestone(create Create) (model.Milestone, error) {
	err := global.Validate.Struct(create)
	if err != nil {
		return model.Milestone{}, err
	}
	m := model.Milestone{}
	err = copier.Copy(&m, &create)
	if err != nil {
		return model.Milestone{}, err
	}
	return m, s.q.Milestone.Create(&m)
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

func (s Milestone) Paginate(page int, limit int) ([]*model.Milestone, int64, error) {
	offset := 0
	if page > 1 {
		offset = limit * (page - 1)
	}
	return s.q.Milestone.Order(s.q.Milestone.ID.Desc()).
		Preload(s.relatedCheckpointOrder()).
		FindByPage(offset, limit)
}

func (s Milestone) DeleteById(id uint) error {
	_, err := s.q.Milestone.Where(s.q.Milestone.ID.Eq(id)).Delete()
	return err
}

func (s Milestone) RetrieveById(id uint) (*model.Milestone, error) {
	return s.q.Milestone.Where(s.q.Milestone.ID.Eq(id)).Preload(s.relatedCheckpointOrder()).First()
}
