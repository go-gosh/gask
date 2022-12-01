package client

import (
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
)

func PaginateMilestone(params map[string]string) (*repo.Paginator[service.MilestoneView], error) {
	return Paginate[*repo.Paginator[service.MilestoneView]](ApiRoute("milestone"), params)
}

func CreateMilestone(body service.MilestoneCreate) (service.MilestoneView, error) {
	return Create[service.MilestoneView](ApiRoute("milestone"), body)
}

func DeleteMilestone(id uint) error {
	_, err := Delete(ApiRoute("milestone/{id}"), id)
	return err
}

func RetrieveMilestone(id uint) (service.MilestoneView, error) {
	return Retrieve[service.MilestoneView](ApiRoute("milestone/{id}"), id)
}

func UpdateMilestone(id uint, update service.MilestoneUpdate) error {
	_, err := Update(ApiRoute("milestone/{id}"), id, update)
	return err
}
