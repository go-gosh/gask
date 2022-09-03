package service

import (
	"errors"
	"time"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gestful/component/mapper"
	"github.com/go-gosh/gestful/component/service"
	"gorm.io/gorm"
)

type TaskCreateRequest struct {
	ParentId uint   `json:"parent_id"`
	Point    uint8  `json:"point"`
	IsCheck  bool   `json:"is_check"`
	Star     uint8  `json:"star"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	StartAt  int64  `json:"start_at"`
	DeadLine int64  `json:"dead_line"`
}

func (t TaskCreateRequest) MakeCreate() (*model.Task, error) {
	req := &model.Task{
		ParentId: t.ParentId,
		Point:    t.Point,
		IsCheck:  t.IsCheck,
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
	if t.DeadLine != 0 {
		deadline := time.Unix(t.DeadLine, 0)
		req.DeadLine = &deadline
	}

	return req, nil
}

type TaskPageRequest struct {
	service.BasePageRequest
	ParentId *uint `json:"parent_id" form:"parent_id"`
}

func (r TaskPageRequest) MakeWrapper() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if r.ParentId == nil {
			return db
		}
		return db.Where("parent_id=?", *r.ParentId)
	}
}

type TaskUpdateRequest struct {
	ParentId *uint   `json:"parent_id"`
	Point    *uint8  `json:"point"`
	IsCheck  *bool   `json:"is_check"`
	Star     *uint8  `json:"star"`
	Category *string `json:"category"`
	Title    *string `json:"title"`
	Detail   *string `json:"detail"`
	StartAt  *int64  `json:"start_at"`
	DeadLine *int64  `json:"dead_line"`
}

func (t TaskUpdateRequest) MakeUpdate() (map[string]interface{}, error) {
	updated := make(map[string]interface{})
	if t.ParentId != nil {
		updated["parent_id"] = *t.ParentId
	}
	if t.Point != nil {
		updated["point"] = *t.Point
	}
	if t.IsCheck != nil {
		updated["is_check"] = *t.IsCheck
	}
	if t.Star != nil {
		updated["start"] = *t.Star
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
	if t.DeadLine != nil {
		if *t.DeadLine == 0 {
			updated["deadline"] = nil
		} else {
			updated["deadline"] = time.Unix(*t.DeadLine, 0)
		}
	}
	if len(updated) == 0 {
		return nil, errors.New("noting changed")
	}
	return updated, nil
}

type Task struct {
	service.BaseRestfulService[model.Task, TaskCreateRequest, TaskPageRequest, TaskUpdateRequest]
}

func NewTask(mapper mapper.Mapper[model.Task]) *Task {
	return &Task{
		BaseRestfulService: service.NewBaseService[model.Task, TaskCreateRequest, TaskPageRequest, TaskUpdateRequest](mapper),
	}
}
