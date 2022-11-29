package global

import (
	"fmt"
	"os"
	"path"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/conf"
	"github.com/go-gosh/gask/app/model"
)

var (
	_once sync.Once
	_db   *gorm.DB
)

func setupDatabase() (*gorm.DB, error) {
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
	//if database.Debug {
	db = db.Debug()
	//}
	_ = db.AutoMigrate(&model.Milestone{}, &model.Checkpoint{})
	return db, nil
}

func GetDatabase() (*gorm.DB, error) {
	var err error
	_once.Do(func() {
		_db, err = setupDatabase()
	})
	return _db, err
}
