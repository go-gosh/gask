package service

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/util"
	"github.com/go-gosh/gestful/component/mapper"
	"github.com/go-gosh/gestful/component/service"
)

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
	view, err := util.Map(data.Data, t.NewTaskViewResp)
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
