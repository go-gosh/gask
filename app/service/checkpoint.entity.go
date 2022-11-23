package service

import (
	"time"

	"github.com/go-gosh/gask/app/model"
)

type CheckpointUpdate struct {
}

type CheckpointCreate struct {
	Point     int `validate:"gt=0"`
	Content   string
	JoinedAt  time.Time  `validate:"required"`
	CheckedAt *time.Time `validate:"omitempty,gtefield=JoinedAt"`
}

type CheckpointView struct {
	model.Checkpoint
	Diff float64
}

type CheckpointQuery struct {
	Timestamp   time.Time
	MilestoneId uint
}
