package repo

import (
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gestful/component/mapper"
	"gorm.io/gorm"
)

type TaskRepo struct {
	mapper.Mapper[model.Task]
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{
		Mapper: mapper.NewBaseMapper[model.Task](db),
	}
}
