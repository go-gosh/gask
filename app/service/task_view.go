package service

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/util"
	"github.com/go-gosh/gestful/component/mapper"
	"gorm.io/gorm"
)

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
