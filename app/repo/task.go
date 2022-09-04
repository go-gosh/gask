package repo

import (
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gestful/component/mapper"
	"gorm.io/gorm"
)

type TaskRepo mapper.CRUDMapper[model.Task]

func NewTaskRepo(db *gorm.DB) TaskRepo {
	return mapper.NewCRUDMapper[model.Task](db)
}
