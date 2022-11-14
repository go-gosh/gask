package milestone

import (
	"errors"
	"time"

	"github.com/go-gosh/gask/app/global"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/query"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Create struct {
	Point     int    `validate:"gt=0"`
	Title     string `validate:"required"`
	Content   string
	StartedAt time.Time  `validate:"required"`
	Deadline  *time.Time `validate:"omitempty,gtefield=StartedAt"`
}

type CheckpointCreate struct {
	Point     int
	Content   string
	JoinedAt  time.Time
	CheckedAt *time.Time
}

type View struct {
	model.Milestone
	Progress int
}

type Service struct {
	q *query.Query
}

func New(db *gorm.DB) *Service {
	return &Service{
		q: query.Use(db),
	}
}

func (s Service) CreateMilestone(create Create) (model.Milestone, error) {
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

func (s Service) SplitMilestoneById(id uint, first CheckpointCreate, checkpoints ...CheckpointCreate) (error, []*model.Checkpoint) {
	m, err := s.q.Milestone.Select(s.q.Milestone.Point, s.q.Milestone.Progress).Where(s.q.Milestone.ID.Eq(id)).First()
	if err != nil {
		return err, nil
	}
	convert := func(c CheckpointCreate) *model.Checkpoint {
		return &model.Checkpoint{
			Point:       c.Point,
			MilestoneId: id,
			Content:     c.Content,
			JoinedAt:    c.JoinedAt,
			CheckedAt:   c.CheckedAt,
		}
	}
	results := make([]*model.Checkpoint, 0, 1+len(checkpoints))
	results = append(results, convert(first))
	sum := first.Point
	for _, checkpoint := range checkpoints {
		sum += checkpoint.Point
		results = append(results, convert(checkpoint))
	}
	if sum > m.Point-m.Progress {
		return errors.New("total points less than checkpoints"), nil
	}
	return s.q.Checkpoint.Create(results...), results
}

func (s Service) CompleteCheckpointById(id uint, timestamp time.Time) error {
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

func (s Service) updateMilestoneProgress(tx *query.Query, id uint, point int) error {
	_, err := tx.Milestone.
		Where(query.Milestone.ID.Eq(id)).
		Update(query.Milestone.Progress, query.Milestone.Progress.Add(point))
	return err
}

func (s Service) ViewAllMilestone() ([]*model.Milestone, error) {
	return s.q.Milestone.Order(query.Milestone.ID.Desc()).Find()
}

func (s Service) Paginate(page int, limit int) ([]*model.Milestone, int64, error) {
	offset := 0
	if page > 1 {
		offset = limit * (page - 1)
	}
	return s.q.Milestone.Order(s.q.Milestone.ID.Desc()).FindByPage(offset, limit)
}
