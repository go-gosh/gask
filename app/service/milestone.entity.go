package service

import (
	"time"

	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type MilestoneCreate struct {
	Title     string     `binding:"required" json:"title"`
	StartedAt time.Time  `binding:"required" json:"startedAt"`
	Deadline  *time.Time `binding:"omitempty,gtefield=StartedAt" json:"deadline"`
}

type MilestoneQuery struct {
	repo.Pager
	OrderBy  []string `form:"orderBy"`
	HideDone bool     `form:"withUndo"`
	Tag      string   `form:"tag"`
}

func (q *MilestoneQuery) injectDB(db *gorm.DB) *gorm.DB {
	if len(q.OrderBy) == 0 {
		q.OrderBy = append(
			q.OrderBy,
			"`milestones`.`progress` < `milestones`.`point` desc",
			"`milestones`.`deadline` is null",
			"`milestones`.`deadline`,`milestones`.`started_at` desc",
			"`milestones`.`id` desc",
		)
	}
	db = db.Order(repo.ArrayToQueryOrder(q.OrderBy))

	if q.HideDone {
		db = db.Where("`milestones`.`point` > `milestones`.`progress`")
	}
	if q.Tag != "" {
		db = db.Select("`milestones`.*").
			Joins("INNER JOIN `milestone_tags` ON `milestones`.`id` = `milestone_tags`.`milestone_id`").
			Where("`milestone_tags`.`name` = ?", q.Tag)
	}
	return db
}

type MilestoneView struct {
	ID        uint                  `json:"id"`
	Point     int                   `json:"point"`
	Progress  int                   `json:"progress"`
	Title     string                `json:"title"`
	StartedAt time.Time             `json:"startedAt"`
	Deadline  *time.Time            `json:"deadline"`
	CreatedAt time.Time             `json:"createdAt"`
	UpdatedAt time.Time             `json:"updatedAt"`
	Tags      []*model.MilestoneTag `json:"tags" gorm:"foreignKey:MilestoneID"`
}

func (MilestoneView) TableName() string {
	return "milestones"
}

type MilestoneUpdate struct {
	Title     string
	StartedAt time.Time
	Deadline  *time.Time
}
