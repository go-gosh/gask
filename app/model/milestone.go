package model

import "time"

type Milestone struct {
	Model
	Point     int       `gorm:"not null"`
	Type      string    `gorm:"not null"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	StartedAt time.Time `gorm:"not null"`
	Deadline  *time.Time
}
