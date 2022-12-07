package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/go-gosh/gask/api/resource"
	"github.com/go-gosh/gask/app/conf"
	"github.com/go-gosh/gask/app/service"
)

func New(milestone service.IMilestone, checkpoint service.ICheckpoint, milestoneTag service.IMilestoneTag) *gin.Engine {
	if conf.GetConfig().Database.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.Default()
	_ = engine.SetTrustedProxies(nil)
	engine.Use(cors.Default())
	resource.Setup(engine)
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

	mt := &MilestoneTag{svc: milestoneTag}
	apiV1.GET("/milestone-tag", mt.Paginate)
	apiV1.POST("/milestone-tag", mt.Create)
	apiV1.DELETE("/milestone/:id/tag/:name", mt.Delete)
	return engine
}
