package service

import (
	"time"

	"gorm.io/gorm"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type CheckpointUpdate struct {
	IsChecked *bool
	Point     *int
	Content   *string
	CheckedAt *time.Time
	JoinedAt  *time.Time
}

func (u CheckpointUpdate) updateDB() map[string]any {
	m := make(map[string]any)
	if u.IsChecked != nil {
		if u.CheckedAt == nil {
			u.CheckedAt = tk.Pointer(time.Now())
		}
		if *u.IsChecked {
			m["checked_at"] = u.CheckedAt
		} else {
			m["checked_at"] = nil
		}
	}
	if u.Point != nil {
		m["point"] = *u.Point
	}
	if u.Content != nil {
		m["content"] = *u.Content
	}
	if u.JoinedAt != nil {
		m["joined_at"] = *u.JoinedAt
	}
	return m
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
