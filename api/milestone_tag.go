package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/api/views"
	"github.com/go-gosh/gask/app/service"
)

type MilestoneTag struct {
	svc service.IMilestoneTag
}

// Paginate MilestoneTag data
//
// @Summary     Paginate MilestoneTag
// @Description Paginate data of MilestoneTag
// @Tags        milestone-tag
// @Accept      json
// @Produce     json
// @Param       request query    service.MilestoneTagQuery true "query params"
// @Success     200  {object} views.JsonResult{data=repo.Paginator[model.MilestoneTag]{data=[]string}}
// @Router      /milestone-tag [get]
func (m MilestoneTag) Paginate(ctx *gin.Context) {
	views.PaginateApi(ctx, m.svc.FindByPage)
}

// PaginateByMilestone MilestoneTag data
//
// @Summary     Paginate MilestoneTag by milestone id
// @Description Paginate MilestoneTag by milestone id
// @Tags        milestone-tag
// @Accept      json
// @Produce     json
// @Param       id   path     uint   true "milestone id"
// @Param       name query    string true "name"
// @Success     200     {object} views.JsonResult{data=repo.Paginator[model.MilestoneTag]{data=[]string}}
// @Router      /milestone/{id}/tag [get]
func (m MilestoneTag) PaginateByMilestone(ctx *gin.Context) {
	err := ctx.Request.ParseForm()
	if err != nil {
		views.Error(ctx, http.StatusInternalServerError, err)
		return
	}
	id := struct {
		ID int `uri:"id" binding:"gt=0"`
	}{}
	err = ctx.ShouldBindUri(&id)
	if err != nil {
		views.Error(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.Request.Form.Set("milestoneId", strconv.Itoa(id.ID))
	m.Paginate(ctx)
}

// Create MilestoneTag
//
// @Summary     Create MilestoneTag
// @Description Create a MilestoneTag
// @Tags        milestone-tag
// @Accept      json
// @Produce     json
// @Param       request body     model.MilestoneTag true "data"
// @Success     200     {object} views.JsonResult{data=string}
// @Router      /milestone-tag [post]
func (m MilestoneTag) Create(ctx *gin.Context) {
	views.CreateApi(ctx, m.svc.Create)
}

// Delete MilestoneTag
//
// @Summary     Delete MilestoneTag
// @Description Delete a MilestoneTag
// @Tags        milestone-tag
// @Accept      json
// @Produce     json
// @Param       id   path     uint true "milestone id"
// @Param       name path     uint true "tag name"
// @Success     200  {object} views.JsonResult
// @Router      /milestone/{id}/tag/{name} [delete]
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
