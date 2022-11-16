package global

import (
	"fmt"
	"os"
	"path"

	"github.com/go-gosh/gask/app/conf"
	"github.com/go-gosh/gask/app/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDatabase() (*gorm.DB, error) {
	database := conf.GetConfig().Database
	dir := path.Dir(database.File)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	db, err := gorm.Open(sqlite.Open(database.File), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		QueryFields:                              true,
	})
	if err != nil {
		return nil, fmt.Errorf("<file:%s> %w", database.File, err)
	}
	if database.Debug {
		db = db.Debug()
	}
	if err := db.AutoMigrate(&model.Milestone{}, &model.Checkpoint{}); err != nil {
		return nil, err
	}
	return db, nil
}
