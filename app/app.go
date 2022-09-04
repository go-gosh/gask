package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
	gservice "github.com/go-gosh/gestful/component/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func New() (*gin.Engine, error) {
	db, err := gorm.Open(sqlite.Open("./data.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.Task{})
	engine := gin.Default()
	engine.Use(CORSMiddleware())
	api := engine.Group("/api/v1")
	{
		taskService := service.NewTask(repo.NewTaskRepo(db))
		gservice.RegisterGroupRoute(api, "task", taskService)
	}

	return engine, nil
}
