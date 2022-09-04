package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gestful/component/mapper"
	"github.com/go-gosh/gestful/component/service"
	"gorm.io/gorm"
)

type TaskPageRequest struct {
	mapper.CRUDPageResult[TaskViewResp]
	ParentId *uint `json:"parent_id" form:"parent_id"`
}

func (r TaskPageRequest) MakeWrapper() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Order("id DESC")
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

type TaskIdUri struct {
	Id uint `json:"id" uri:"id" binding:"required"`
}

type TaskViewResp struct {
	ID        uint   `json:"id"`
	ParentId  uint   `json:"parent_id"`
	Point     uint8  `json:"point"`
	IsCheck   bool   `json:"is_check"`
	Star      uint8  `json:"star"`
	Category  string `json:"category"`
	Title     string `json:"title"`
	Detail    string `json:"detail"`
	StartAt   int64  `json:"start_at,string"`
	DeadLine  int64  `json:"dead_line,string"`
	CreatedAt int64  `json:"created_at,string"`
	UpdatedAt int64  `json:"updated_at,string"`
}

type TaskCreateRequest struct {
	ParentId uint   `json:"parent_id"`
	Point    uint8  `json:"point"`
	IsCheck  bool   `json:"is_check"`
	Star     uint8  `json:"star"`
	Category string `json:"category"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	StartAt  int64  `json:"start_at,string"`
	DeadLine int64  `json:"dead_line,string"`
}

func (t TaskCreateRequest) ToEntity() (*model.Task, error) {
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

func NewTask(repo repo.TaskRepo) service.RestfulService[TaskViewResp, mapper.CRUDPageResult[TaskViewResp]] {
	return &task{
		repo: repo,
	}
}

type task struct {
	repo repo.TaskRepo
}

func (t task) Create(ctx *gin.Context) error {
	var req TaskCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return err
	}
	et, err := req.ToEntity()
	if err != nil {
		return err
	}
	return t.repo.Create(ctx, et)
}

func (t task) Paginate(ctx *gin.Context) (*mapper.CRUDPageResult[TaskViewResp], error) {
	var req TaskPageRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, err
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}
	if req.PageSize > 500 {
		req.PageSize = 500
	}
	data, err := t.repo.Paginate(ctx, req.CRUDPaginator, req.MakeWrapper())
	if err != nil {
		return nil, err
	}
	view, err := Map(data.Data, t.NewTaskViewResp)
	if err != nil {
		return nil, err
	}
	return &mapper.CRUDPageResult[TaskViewResp]{
		CRUDPaginator: data.CRUDPaginator,
		Total:         data.Total,
		TotalPage:     data.TotalPage,
		Data:          view,
	}, nil
}

func (t task) Retrieve(ctx *gin.Context) (*TaskViewResp, error) {
	var id TaskIdUri
	if err := ctx.ShouldBindUri(&id); err != nil {
		return nil, err
	}
	et, err := t.repo.OneById(ctx, id.Id)
	if err != nil {
		return nil, err
	}

	return t.NewTaskViewResp(et)
}

func Map[T, U any](data []T, fn func(*T) (*U, error)) ([]U, error) {
	res := make([]U, 0, len(data))
	for i := 0; i < len(data); i++ {
		u, err := fn(&data[i])
		if err != nil {
			return nil, err
		}
		res = append(res, *u)
	}
	return res, nil
}

func (t task) NewTaskViewResp(et *model.Task) (*TaskViewResp, error) {
	if et == nil {
		return nil, nil
	}

	var deadLine int64
	if et.DeadLine != nil {
		deadLine = et.DeadLine.Unix()
	}
	return &TaskViewResp{
		ID:        et.ID,
		ParentId:  et.ParentId,
		Point:     et.Point,
		IsCheck:   et.IsCheck,
		Star:      et.Star,
		Category:  et.Category,
		Title:     et.Title,
		Detail:    et.Detail,
		StartAt:   et.StartAt.Unix(),
		DeadLine:  deadLine,
		CreatedAt: et.CreatedAt.Unix(),
		UpdatedAt: et.UpdatedAt.Unix(),
	}, nil
}

func (t task) Update(ctx *gin.Context) error {
	var id TaskIdUri
	if err := ctx.ShouldBindUri(&id); err != nil {
		return err
	}
	var updated TaskUpdateRequest
	if err := ctx.ShouldBindJSON(&updated); err != nil {
		return err
	}
	m, err := updated.MakeUpdate()
	if err != nil {
		return err
	}
	return t.repo.UpdateById(ctx, id.Id, m)
}

func (t task) Delete(ctx *gin.Context) error {
	var id TaskIdUri
	if err := ctx.ShouldBindUri(&id); err != nil {
		return err
	}
	return t.repo.DeleteById(ctx, id.Id)
}
