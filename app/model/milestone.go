package model

import "time"

type milestone struct {
	Model
	Point     int       `gorm:"not null"`
	Title     string    `gorm:"not null"`
	Content   string    `gorm:"not null"`
	StartedAt time.Time `gorm:"not null"`
	Deadline  *time.Time
}
