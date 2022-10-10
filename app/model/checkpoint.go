package model

import "time"

type Checkpoint struct {
	Model
	Point       int       `gorm:"not null"`
	MilestoneId uint      `gorm:"not null"`
	Content     string    `gorm:"not null"`
	JoinedAt    time.Time `gorm:"not null"`
	CheckedAt   *time.Time
}
