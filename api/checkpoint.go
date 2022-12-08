package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-gosh/gask/api/views"
	"github.com/go-gosh/gask/app/service"
)

type Checkpoint struct {
	svc service.ICheckpoint
}

// Paginate Checkpoint data
//
// @Summary     Paginate Checkpoint
// @Description Paginate data of Checkpoint
// @Tags        checkpoint
// @Accept      json
// @Produce     json
// @Param       request query    service.CheckpointQuery true "query params"
// @Success     200     {object} views.JsonResult{data=repo.Paginator[service.CheckpointView]{data=[]service.CheckpointView}}
// @Router      /checkpoint [get]
func (c Checkpoint) Paginate(ctx *gin.Context) {
	views.PaginateApi(ctx, c.svc.FindByPage)
}

// Create Checkpoint
//
// @Summary     Create Checkpoint
// @Description Create a Checkpoint
// @Tags        checkpoint
// @Accept      json
// @Produce     json
// @Param       request body     service.CheckpointCreate true "data"
// @Success     200     {object} views.JsonResult{data=service.CheckpointView}
// @Router      /checkpoint [post]
func (c Checkpoint) Create(ctx *gin.Context) {
	views.CreateApi(ctx, c.svc.Create)
}

// Delete Checkpoint
//
// @Summary     Delete Checkpoint
// @Description Delete a Checkpoint
// @Tags        checkpoint
// @Accept      json
// @Produce     json
// @Param       id  path     uint true "checkpoint id"
// @Success     200 {object} views.JsonResult
// @Router      /checkpoint/{id} [delete]
func (c Checkpoint) Delete(ctx *gin.Context) {
	views.DeleteByIdApi(ctx, c.svc.DeleteById)
}

// Update Checkpoint
//
// @Summary     Update Checkpoint
// @Description Update a Checkpoint
// @Tags        checkpoint
// @Accept      json
// @Produce     json
// @Param       id      path     uint                     true "checkpoint id"
// @Param       request body     service.CheckpointUpdate true "data"
// @Success     200     {object} views.JsonResult
// @Router      /checkpoint/{id} [put]
func (c Checkpoint) Update(ctx *gin.Context) {
	views.UpdateByIdApi(ctx, c.svc.UpdateById)
}

// Retrieve Checkpoint
//
// @Summary     Retrieve Checkpoint
// @Description Retrieve a Checkpoint
// @Tags        checkpoint
// @Accept      json
// @Produce     json
// @Param       id  path     uint true "checkpoint id"
// @Success     200 {object} views.JsonResult{data=service.CheckpointView}
// @Router      /checkpoint/{id} [get]
func (c Checkpoint) Retrieve(ctx *gin.Context) {
	views.OneByIdApi(ctx, c.svc.OneById)
}
