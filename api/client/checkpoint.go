package client

import (
	"github.com/go-gosh/gask/app/repo"
	"github.com/go-gosh/gask/app/service"
)

func PaginateCheckpoint(params map[string]string) (*repo.Paginator[service.CheckpointView], error) {
	return Paginate[*repo.Paginator[service.CheckpointView]](ApiRoute("checkpoint"), params)
}

func CreateCheckpoint(body service.CheckpointCreate) (service.CheckpointView, error) {
	return Create[service.CheckpointView](ApiRoute("checkpoint"), body)
}

func DeleteCheckpoint(id uint) error {
	_, err := Delete(ApiRoute("checkpoint/{id}"), id)
	return err
}

func RetrieveCheckpoint(id uint) (service.CheckpointView, error) {
	return Retrieve[service.CheckpointView](ApiRoute("checkpoint/{id}"), id)
}

func UpdateCheckpoint(id uint, update service.CheckpointUpdate) error {
	_, err := Update(ApiRoute("checkpoint/{id}"), id, update)
	return err
}
