package service

import "time"

type Create struct {
	Point     int    `validate:"gt=0"`
	Title     string `validate:"required"`
	Content   string
	StartedAt time.Time  `validate:"required"`
	Deadline  *time.Time `validate:"omitempty,gtefield=StartedAt"`
}

type CheckpointCreate struct {
	Point     int `validate:"gt=0"`
	Content   string
	JoinedAt  time.Time  `validate:"required"`
	CheckedAt *time.Time `validate:"omitempty,gtefield=JoinedAt"`
}
