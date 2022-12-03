package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/repo"
)

type MilestoneCreate struct {
	Point     int        `binding:"gt=0"`
	Title     string     `binding:"required"`
	StartedAt time.Time  `binding:"required"`
	Deadline  *time.Time `binding:"omitempty,gtefield=StartedAt"`
}

type MilestoneQuery struct {
	repo.Pager
	OrderBy  []string `form:"orderBy"`
	HideDone bool     `form:"withUndo"`
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
	ID        uint       `json:"id"`
	Point     int        `json:"point"`
	Progress  int        `json:"progress"`
	Title     string     `json:"title"`
	StartedAt time.Time  `json:"startedAt"`
	Deadline  *time.Time `json:"deadline"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	IsDeleted bool       `json:"isDeleted" gorm:"-"`
}

func (MilestoneView) TableName() string {
	return "milestones"
}

type MilestoneUpdate struct {
	Point     int
	Title     string
	StartedAt time.Time
	Deadline  *time.Time
}
