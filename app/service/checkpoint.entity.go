package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type CheckpointUpdate struct {
}

type CheckpointCreate struct {
	Point       int        `validate:"gt=0"`
	Content     string     `validate:"required"`
	JoinedAt    time.Time  `validate:"required"`
	CheckedAt   *time.Time `validate:"omitempty,gtefield=JoinedAt"`
	MilestoneId uint       `validate:"required"`
}

type CheckpointView struct {
	model.Checkpoint
	Diff float64
}

type CheckpointQuery struct {
	repo.Pager
	Timestamp     *time.Time
	MilestoneId   uint
	WithMilestone bool
}

func (q CheckpointQuery) injectDB(db *gorm.DB) *gorm.DB {
	if q.Timestamp != nil {
		db = db.Select("*, julianday(joined_at) - julianday(?) as diff", q.Timestamp)
	}
	if q.MilestoneId != 0 {
		db = db.Where("milestone_id = ?", q.MilestoneId)
	}
	if q.WithMilestone {
		db = db.Preload("Milestone")
	}
	return db
}
