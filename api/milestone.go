package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-gosh/gask/api/views"
	"github.com/go-gosh/gask/app/service"
)

type Milestone struct {
	svc service.IMilestone
}

func (m Milestone) Paginate(ctx *gin.Context) {
	views.PaginateApi(ctx, m.svc.FindByPage)
}

func (m Milestone) Create(ctx *gin.Context) {
	views.CreateApi(ctx, m.svc.Create)
}

func (m Milestone) Delete(ctx *gin.Context) {
	views.DeleteByIdApi(ctx, m.svc.DeleteById)
}

func (m Milestone) Retrieve(ctx *gin.Context) {
	views.OneByIdApi(ctx, m.svc.OneById)
}

func (m Milestone) Update(ctx *gin.Context) {
	views.UpdateByIdApi(ctx, m.svc.UpdateById)
}
