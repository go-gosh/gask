package service

import (
	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gask/app/repo"
)

type MilestoneTagQuery struct {
	repo.Pager
	model.MilestoneTag
}
