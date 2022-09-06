package service

import (
	"time"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/util"
)

type TaskCreateRequest struct {
	ParentId   uint   `json:"parent_id"`
	Point      uint8  `json:"point"`
	IsCheck    bool   `json:"is_check"`
	Star       uint8  `json:"star"`
	Category   string `json:"category"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	StartAt    int64  `json:"start_at,string"`
	Deadline   int64  `json:"deadline,string"`
	CompleteAt int64  `json:"complete_at,string"`
}

func (t TaskCreateRequest) ToEntity() (*model.Task, error) {
	req := &model.Task{
		ParentId: t.ParentId,
		Point:    t.Point,
		Star:     t.Star,
		Category: t.Category,
		Title:    t.Title,
		Detail:   t.Detail,
	}
	if t.StartAt == 0 {
		req.StartAt = time.Now()
	} else {
		req.StartAt = time.Unix(t.StartAt, 0)
	}
	if t.Deadline != 0 {
		req.Deadline = util.Point(time.Unix(t.Deadline, 0))
	}
	if t.CompleteAt != 0 {
		req.CompleteAt = util.Point(time.Unix(t.CompleteAt, 0))
	}

	return req, nil
}
