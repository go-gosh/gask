package api

import (
	"github.com/gin-gonic/gin"

	"github.com/go-gosh/gask/app/service"
)

func New(milestone service.IMilestone, checkpoint service.ICheckpoint) *gin.Engine {
	engine := gin.Default()

	apiV1 := engine.Group("/api/v1")
	m := &Milestone{svc: milestone}
	apiV1.GET("/milestone", m.Paginate)
	apiV1.POST("/milestone", m.Create)
	apiV1.GET("/milestone/:id", m.Retrieve)
	apiV1.PUT("/milestone/:id", m.Update)
	apiV1.DELETE("/milestone/:id", m.Delete)

	c := &Checkpoint{svc: checkpoint}
	apiV1.GET("/checkpoint", c.Paginate)
	apiV1.POST("/checkpoint", c.Create)
	apiV1.GET("/checkpoint/:id", c.Retrieve)
	apiV1.PUT("/checkpoint/:id", c.Update)
	apiV1.DELETE("/checkpoint/:id", c.Delete)
	return engine
}
