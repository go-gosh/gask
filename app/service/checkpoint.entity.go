package service

import (
	"time"

	"gorm.io/gorm"

	tk "github.com/go-gosh/gask/app/common/toolkit"
	"github.com/go-gosh/gask/app/repo"
)

type CheckpointUpdate struct {
	IsChecked *bool      `json:"isChecked"`
	Point     *int       `json:"point"`
	Content   *string    `json:"content"`
	CheckedAt *time.Time `json:"checkedAt"`
	JoinedAt  *time.Time `json:"joinedAt"`
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
	Point       int        `binding:"gt=0" json:"point"`
	Content     string     `binding:"required" json:"content"`
	JoinedAt    time.Time  `binding:"required" json:"joinedAt"`
	CheckedAt   *time.Time `binding:"omitempty,gtefield=JoinedAt" json:"checkedAt"`
	MilestoneId uint       `binding:"required" json:"milestoneId"`
}

type CheckpointView struct {
	ID          uint          `json:"id"`
	Point       int           `json:"point"`
	MilestoneId uint          `json:"milestoneId"`
	Content     string        `json:"content"`
	JoinedAt    time.Time     `json:"joinedAt"`
	CheckedAt   *time.Time    `json:"checkedAt"`
	Milestone   MilestoneView `json:"milestone"`
	Diff        float64       `json:"diff"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	DeletedAt   *time.Time    `json:"deletedAt"`
}

func (CheckpointView) TableName() string {
	return "checkpoints"
}

type CheckpointQuery struct {
	repo.Pager
	Timestamp     *time.Time `form:"timestamp"`
	MilestoneId   uint       `form:"milestoneId"`
	WithMilestone bool       `form:"withMilestone"`
}

func (q CheckpointQuery) injectDB(db *gorm.DB) *gorm.DB {
	if q.Timestamp != nil {
		db = db.Select("*, julianday(joined_at) - julianday(?) as diff", q.Timestamp)
	} else {
		db = db.Omit("diff")
	}
	if q.MilestoneId != 0 {
		db = db.Where("milestone_id = ?", q.MilestoneId)
	}
	if q.WithMilestone {
		db = db.Preload("Milestone")
	}
	return db
}
