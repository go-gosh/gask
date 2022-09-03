package app

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() (*gin.Engine, error) {
	db, err := gorm.Open(sqlite.Open("./data.sqlite3"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	engine := gin.Default()
	api := engine.Group("/api/v1")
	{
		taskService := service.NewTask(repo.NewTaskRepo(db))
		taskService.RegisterGroupRoute(api, "task")
	}

	return engine, nil
}
