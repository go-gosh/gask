package milestone

import (
	"errors"
	"time"

	"github.com/go-gosh/gask/app/model"
	"gorm.io/gorm"
)

type Create struct {
	Point     int
	Title     string
	Content   string
	StartedAt time.Time
	Deadline  *time.Time
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
	db *gorm.DB
}

func New(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s Service) CreateMilestone(create Create) (model.Milestone, error) {
	m := model.Milestone{
		Point:     create.Point,
		Title:     create.Title,
		Content:   create.Content,
		StartedAt: create.StartedAt,
		Deadline:  create.Deadline,
	}
	return m, s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Milestone{}).Create(&m).Error
	})
}

func (s Service) SplitMilestoneById(id uint, first CheckpointCreate, checkpoints ...CheckpointCreate) (error, []model.Checkpoint) {
	m := struct {
		Id    uint
		Point int
	}{}
	err := s.db.Model(&model.Milestone{}).Find(&m, id).Error
	if err != nil {
		return err, nil
	}
	convert := func(c CheckpointCreate) model.Checkpoint {
		return model.Checkpoint{
			Point:       c.Point,
			MilestoneId: id,
			Content:     c.Content,
			JoinedAt:    c.JoinedAt,
			CheckedAt:   c.CheckedAt,
		}
	}
	results := make([]model.Checkpoint, 0, 1+len(checkpoints))
	results = append(results, convert(first))
	sum := first.Point
	for _, checkpoint := range checkpoints {
		sum += checkpoint.Point
		results = append(results, convert(checkpoint))
	}
	if sum > m.Point {
		return errors.New("total points less than checkpoints"), nil
	}
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Checkpoint{}).Create(results).Error
	}), results
}

func (s Service) CompleteCheckpointById(id uint, timestamp time.Time) error {
	c := struct {
		Id        uint
		CheckedAt *time.Time
	}{}
	err := s.db.Model(&model.Checkpoint{}).First(&c, id).Error
	if err != nil {
		return err
	}
	if c.CheckedAt != nil {
		return errors.New("checkpoint already completed")
	}
	c.CheckedAt = &timestamp
	return s.db.Transaction(func(tx *gorm.DB) error {
		return tx.Model(&model.Checkpoint{}).Updates(c).Error
	})
}

func (s Service) ViewAllMilestone() ([]View, error) {
	m := make([]View, 0)
	q := s.db.Model(&model.Checkpoint{}).Select("SUM(point) AS progress, milestone_id").Group("milestone_id")
	err := s.db.Model(&model.Milestone{}).Select("milestone.*, q.progress").Joins("LEFT JOIN (?) AS q ON q.milestone_id=milestone.id", q).Find(&m).Error
	return m, err
}
