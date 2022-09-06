package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/util"
	"github.com/go-gosh/gestful/component/mapper"
	"github.com/go-gosh/gestful/component/service"
	"gorm.io/gorm"
)

type TaskPageRequest struct {
	mapper.CRUDPageResult[TaskViewResp]
	ParentId *uint    `json:"parent_id" form:"parent_id"`
	OrderBy  []string `form:"order_by"`
	Datetime string   `json:"datetime" form:"datetime"`
	Process  string   `json:"process" form:"process"`
}

func (r TaskPageRequest) MakeWrapper() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = r.orderFunc(db)
		if r.ParentId != nil {
			db = db.Where("parent_id=?", *r.ParentId)
		}
		switch len(r.Datetime) {
		case 6:
			// Monthly
			t, err := time.Parse("200601", r.Datetime)
			if err != nil {
				break
			}
			db = db.Scopes(repo.WrapperDateRangeActive(t, t.AddDate(0, 1, 0)))
		case 8:
			// Daily
			t, err := time.Parse("20060102", r.Datetime)
			if err != nil {
				break
			}
			db = db.Scopes(repo.WrapperDateRangeActive(t, t.AddDate(0, 0, 1)))
		}
		if len(r.Process) == 8 {
			t, err := time.Parse("20060102", r.Process)
			if err == nil {
				db = db.Where("complete_at>=? and complete_at<?", t, t.AddDate(0, 0, 1))
			}
		}
		return db
	}
}

func (r TaskPageRequest) orderFunc(db *gorm.DB) *gorm.DB {
	pk := false
	if len(r.OrderBy) > 0 {
		for _, s := range r.OrderBy {
			if s == "" {
				continue
			}
			desc := "ASC"
			if s[0] == '-' {
				desc = "DESC"
			}
			k := strings.ToLower(strings.TrimPrefix(strings.TrimPrefix(s, "+"), "-"))
			if !r.isOrderKey(k) {
				continue
			}
			if k == "id" {
				pk = true
			}
			db = db.Order(fmt.Sprintf("%s %s", k, desc))
		}
	}
	if !pk {
		db = db.Order("id DESC")
	}
	return db
}

var taskOrderKey = map[string]struct{}{
	"id":          {},
	"parent_id":   {},
	"point":       {},
	"star":        {},
	"category":    {},
	"start_at":    {},
	"cancel_at":   {},
	"complete_at": {},
	"deadline":    {},
}

func (r TaskPageRequest) isOrderKey(k string) bool {
	_, ok := taskOrderKey[k]
	return ok
}

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

type TaskViewResp struct {
	ID         uint   `json:"id"`
	ParentId   uint   `json:"parent_id"`
	Point      uint8  `json:"point"`
	Star       uint8  `json:"star"`
	Category   string `json:"category"`
	Title      string `json:"title"`
	Detail     string `json:"detail"`
	StartAt    int64  `json:"start_at,string"`
	CompleteAt int64  `json:"complete_at,string"`
	CancelAt   int64  `json:"cancel_at,string"`
	Deadline   int64  `json:"deadline,string"`
	CreatedAt  int64  `json:"created_at,string"`
	UpdatedAt  int64  `json:"updated_at,string"`

	Parent *TaskViewResp `json:"parent"`
}

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
	view, err := t.NewTaskViewResp(et)
	if err != nil {
		return nil, err
	}
	if view.ParentId == 0 {
		return view, err
	}
	parent, err := t.repo.OneById(ctx, view.ParentId)
	if err != nil {
		return nil, err
	}
	parentView, err := t.NewTaskViewResp(parent)
	if err != nil {
		return nil, err
	}
	view.Parent = parentView
	return view, err
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

	unix := time.Unix(0, 0)
	return &TaskViewResp{
		ID:         et.ID,
		ParentId:   et.ParentId,
		Point:      et.Point,
		Star:       et.Star,
		Category:   et.Category,
		Title:      et.Title,
		Detail:     et.Detail,
		StartAt:    et.StartAt.Unix(),
		CompleteAt: util.NilToDefault(et.CompleteAt, unix).Unix(),
		CancelAt:   util.NilToDefault(et.CancelAt, unix).Unix(),
		Deadline:   util.NilToDefault(et.Deadline, unix).Unix(),
		CreatedAt:  et.CreatedAt.Unix(),
		UpdatedAt:  et.UpdatedAt.Unix(),
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
