package service

import (
	"errors"
	"time"

	"github.com/go-gosh/gask/app/util"
)

type TaskUpdateRequest struct {
	ParentId   *uint   `json:"parent_id"`
	Point      *uint8  `json:"point"`
	Star       *uint8  `json:"star"`
	Category   *string `json:"category"`
	Title      *string `json:"title"`
	Detail     *string `json:"detail"`
	StartAt    *int64  `json:"start_at"`
	CompleteAt *int64  `json:"complete_at"`
	CancelAt   *int64  `json:"cancel_at"`
	Deadline   *int64  `json:"deadline"`
}

func (t TaskUpdateRequest) MakeUpdate() (map[string]interface{}, error) {
	updated := make(map[string]interface{})
	if t.ParentId != nil {
		updated["parent_id"] = *t.ParentId
	}
	if t.Point != nil {
		updated["point"] = *t.Point
	}
	if t.CompleteAt != nil {
		if *t.CompleteAt == 0 {
			updated["complete_at"] = nil
		} else {
			updated["complete_at"] = util.Point(time.Unix(*t.CompleteAt, 0))
		}
	}
	if t.CancelAt != nil {
		if *t.CancelAt == 0 {
			updated["cancel_at"] = nil
		} else {
			updated["cancel_at"] = util.Point(time.Unix(*t.CancelAt, 0))
		}
	}
	if t.Star != nil {
		updated["star"] = *t.Star
	}
	if t.Category != nil {
		updated["category"] = *t.Category
	}
	if t.Title != nil {
		updated["title"] = *t.Title
	}
	if t.Detail != nil {
		updated["detail"] = *t.Detail
	}
	if t.StartAt != nil {
		updated["start"] = time.Unix(*t.StartAt, 0)
	}
	if t.Deadline != nil {
		if *t.Deadline == 0 {
			updated["deadline"] = nil
		} else {
			updated["deadline"] = time.Unix(*t.Deadline, 0)
		}
	}
	if len(updated) == 0 {
		return nil, errors.New("noting changed")
	}
	return updated, nil
}

type TaskIdUri struct {
	Id uint `json:"id" uri:"id" binding:"required"`
}
