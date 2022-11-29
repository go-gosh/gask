package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-gosh/gask/api/views"
	"github.com/go-gosh/gask/app/service"
)

type Checkpoint struct {
	svc service.ICheckpoint
}

func (c Checkpoint) Paginate(ctx *gin.Context) {
	views.PaginateApi(ctx, c.svc.FindByPage)
}

func (c Checkpoint) Create(ctx *gin.Context) {
	views.CreateApi(ctx, c.svc.Create)
}

func (c Checkpoint) Delete(ctx *gin.Context) {
	views.DeleteByIdApi(ctx, c.svc.DeleteById)
}

func (c Checkpoint) Update(ctx *gin.Context) {
	views.UpdateByIdApi(ctx, c.svc.UpdateById)
}

func (c Checkpoint) Retrieve(ctx *gin.Context) {
	views.OneByIdApi(ctx, c.svc.OneById)
}
