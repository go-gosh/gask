package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/api/views"
	"github.com/go-gosh/gask/app/service"
)

type MilestoneTag struct {
	svc service.IMilestoneTag
}

func (m MilestoneTag) Paginate(ctx *gin.Context) {
	views.PaginateApi(ctx, m.svc.FindByPage)
}

func (m MilestoneTag) Create(ctx *gin.Context) {
	views.CreateApi(ctx, m.svc.Create)
}

func (m MilestoneTag) Delete(ctx *gin.Context) {
	var r struct {
		ID   uint   `uri:"id"`
		Name string `uri:"name"`
	}
	err := ctx.ShouldBindUri(&r)
	if err != nil {
		views.Error(ctx, http.StatusBadRequest, err)
		return
	}
	err = m.svc.DeleteByMilestoneIdAndName(ctx, r.ID, r.Name)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		views.Error(ctx, 404, err)
		return
	}
	if err != nil {
		views.Error(ctx, 500, err)
		return
	}
	views.Success(ctx, nil)
}
