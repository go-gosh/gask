package model

import "time"

type Milestone struct {
	Model
	Point       int       `gorm:"not null"`
	Progress    int       `gorm:"not null"`
	Title       string    `gorm:"not null"`
	StartedAt   time.Time `gorm:"not null"`
	Deadline    *time.Time
	Checkpoints []*Checkpoint
	Tags        []*MilestoneTag
}

func (m Milestone) IsDeleted() bool {
	return m.DeletedAt.Valid
}
