package model

import "time"

type Milestone struct {
	Model
	Point    int    `gorm:"not null"`
	Progress int    `gorm:"not null"`
	Title    string `gorm:"not null"`
	//Deprecated
	Content     string    `gorm:"not null"`
	StartedAt   time.Time `gorm:"not null"`
	Deadline    *time.Time
	Checkpoints []*Checkpoint
}
