package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-gosh/gask/api/views"
	"github.com/go-gosh/gask/app/service"
)

type Milestone struct {
	svc service.IMilestone
}

// Paginate Milestone data
//
// @Summary     Paginate Milestone
// @Description Paginate data of Milestone
// @Tags        milestone
// @Accept      json
// @Produce     json
// @Param       request query    service.MilestoneQuery true "query params"
// @Success     200     {object} views.JsonResult{data=repo.Paginator[service.MilestoneView]{data=[]service.MilestoneView}}
// @Router      /milestone [get]
func (m Milestone) Paginate(ctx *gin.Context) {
	views.PaginateApi(ctx, m.svc.FindByPage)
}

// Create Milestone
//
// @Summary     Create Milestone
// @Description Create a Milestone
// @Tags        milestone
// @Accept      json
// @Produce     json
// @Param       request body     service.MilestoneCreate true "data"
// @Success     200     {object} views.JsonResult{data=service.MilestoneView}
// @Router      /milestone [post]
func (m Milestone) Create(ctx *gin.Context) {
	views.CreateApi(ctx, m.svc.Create)
}

// Delete Milestone
//
// @Summary     Delete Milestone
// @Description Delete a Milestone
// @Tags        milestone
// @Accept      json
// @Produce     json
// @Param       id  path     uint true "milestone id"
// @Success     200 {object} views.JsonResult
// @Router      /milestone/{id} [delete]
func (m Milestone) Delete(ctx *gin.Context) {
	views.DeleteByIdApi(ctx, m.svc.DeleteById)
}

// Retrieve Milestone
//
// @Summary     Retrieve Milestone
// @Description Retrieve a Milestone
// @Tags        milestone
// @Accept      json
// @Produce     json
// @Param       id  path     uint true "milestone id"
// @Success     200 {object} views.JsonResult{data=service.MilestoneView}
// @Router      /milestone/{id} [get]
func (m Milestone) Retrieve(ctx *gin.Context) {
	views.OneByIdApi(ctx, m.svc.OneById)
}

// Update Milestone
//
// @Summary     Update Milestone
// @Description Update a Milestone
// @Tags        milestone
// @Accept      json
// @Produce     json
// @Param       id      path     uint                    true "milestone id"
// @Param       request body     service.MilestoneUpdate true "data"
// @Success     200     {object} views.JsonResult
// @Router      /milestone/{id} [put]
func (m Milestone) Update(ctx *gin.Context) {
	views.UpdateByIdApi(ctx, m.svc.UpdateById)
}
