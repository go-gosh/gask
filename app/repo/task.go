package repo

import (
	"time"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gestful/component/mapper"
	"gorm.io/gorm"
)

type TaskRepo mapper.CRUDMapper[model.Task]

func NewTaskRepo(db *gorm.DB) TaskRepo {
	return mapper.NewCRUDMapper[model.Task](db)
}

func WrapperDateRangeActive(start, end time.Time) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("start_at<? and (dead_line is NULL or dead_line>=?)", end, start)
	}
}
