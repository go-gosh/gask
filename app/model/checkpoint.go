package model

import "time"

type checkpoint struct {
	Model
	Point     int       `gorm:"not null"`
	Content   string    `gorm:"not null"`
	JoinedAt  time.Time `gorm:"not null"`
	CheckedAt *time.Time
}
