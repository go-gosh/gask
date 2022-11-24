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
	OrderBy  []string
	HideDone bool
	scopes   []func(db *gorm.DB) *gorm.DB
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
	if q.HideDone {
		q.scopes = append(q.scopes, func(db *gorm.DB) *gorm.DB {
			return db.Where("point>progress")
		})
	}
	return db.Scopes(q.scopes...)
}

type MilestoneView struct {
	model.Milestone
}

type MilestoneUpdate struct {
	CleanDeadline bool
	Point         *int
	Title         *string
	StartedAt     *time.Time
	Deadline      *time.Time
}

func (u MilestoneUpdate) updateDB() map[string]any {
	m := make(map[string]any)
	if u.CleanDeadline {
		m["deadline"] = nil
	} else if u.Deadline != nil {
		m["deadline"] = *u.Deadline
	}
	if u.Title != nil {
		m["title"] = *u.Title
	}
	if u.Point != nil {
		m["point"] = *u.Point
	}
	if u.StartedAt != nil {
		m["started_at"] = *u.StartedAt
	}
	return m
}
