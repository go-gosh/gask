package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type MilestoneCreate struct {
	Point     int        `validate:"gt=0"`
	Title     string     `validate:"required"`
	StartedAt time.Time  `validate:"required"`
	Deadline  *time.Time `validate:"omitempty,gtefield=StartedAt"`
}

type MilestoneQuery struct {
	repo.Pager
	OrderBy []string
	scopes  []func(db *gorm.DB) *gorm.DB
}

func (q *MilestoneQuery) add(scope func(db *gorm.DB) *gorm.DB) *MilestoneQuery {
	q.scopes = append(q.scopes, scope)
	return q
}

func (q *MilestoneQuery) injectDB(db *gorm.DB) *gorm.DB {
	if len(q.OrderBy) == 0 {
		q.add(func(db *gorm.DB) *gorm.DB {
			return db.Order(repo.DefaultOrderBy)
		})
	} else {
		q.add(func(db *gorm.DB) *gorm.DB {
			return db.Order(repo.ArrayToQueryOrder(q.OrderBy))
		})
	}
	return db.Scopes(q.scopes...)
}

type MilestoneView struct {
	model.Milestone
}

type MilestoneUpdate struct {
}
