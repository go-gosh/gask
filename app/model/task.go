package model

import "time"

type Task struct {
	Model
	ParentId uint       `json:"parent_id" gorm:"not null"`
	Point    uint8      `json:"point" gorm:"default=100;not null"`
	IsCheck  bool       `json:"is_check" gorm:"not null"`
	Star     uint8      `json:"star" gorm:"not null"`
	Category string     `json:"category" gorm:"not null"`
	Title    string     `json:"title" gorm:"not null"`
	Detail   string     `json:"detail" gorm:"not null"`
	StartAt  time.Time  `json:"start_at" gorm:"not null"`
	DeadLine *time.Time `json:"dead_line"`

	SubTask []Task `json:"sub_task" gorm:"foreignKey:ParentId"`
}

func (Task) TableName() string {
	return "task"
}
