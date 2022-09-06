package app

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gask/app/conf"
	"github.com/go-gosh/gask/app/middleware"
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
	gservice "github.com/go-gosh/gestful/component/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() (*gin.Engine, error) {
	database := conf.GetConfig().Database
	dir := path.Dir(database.File)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open(database.File), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("<file:%s> %w", database.File, err)
	}
	if database.Debug {
		db = db.Debug()
	}
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		return nil, err
	}
	engine := gin.Default()
	engine.Use(middleware.CORSMiddleware())
	api := engine.Group("/api/v1")
	{
		taskService := service.NewTask(repo.NewTaskRepo(db))
		gservice.RegisterGroupRoute(api, "task", taskService)
	}

	return engine, nil
}

func Run() error {
	log.Printf("load config:<%+v>", conf.GetConfig())
	engine, err := New()
	if err != nil {
		return err
	}

	return engine.Run(fmt.Sprintf(":%d", conf.GetConfig().Port))
}
