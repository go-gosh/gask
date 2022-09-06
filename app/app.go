package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/middleware"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
	gservice "github.com/go-gosh/gestful/component/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() (*gin.Engine, error) {
	db, err := gorm.Open(sqlite.Open("./data.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.Task{})
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())
	api := engine.Group("/api/v1")
	{
		taskService := service.NewTask(repo.NewTaskRepo(db))
		gservice.RegisterGroupRoute(api, "task", taskService)
	}

	return engine, nil
}
