package views

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/go-gosh/gask/app/repo"
)

func PaginateApi[Q, D any](ctx *gin.Context, fn func(context.Context, *Q) (*repo.Paginator[D], error)) {
	var q Q
	err := ctx.ShouldBind(&q)
	if err != nil {
		Error(ctx, 404, err)
		return
	}
	data, err := fn(ctx, &q)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(ctx, 404, err)
		return
	}
	if err != nil {
		Error(ctx, 500, err)
		return
	}
	Success(ctx, data)
}

func CreateApi[C, V any](ctx *gin.Context, fn func(context.Context, C) (V, error)) {
	var c C
	err := ctx.ShouldBindJSON(&c)
	if err != nil {
		Error(ctx, 400, err)
		return
	}
	data, err := fn(ctx, c)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(ctx, 404, err)
		return
	}
	if err != nil {
		Error(ctx, 500, err)
		return
	}
	Success(ctx, data)
}

func DeleteByIdApi(ctx *gin.Context, fn func(context.Context, uint, ...uint) error) {
	id := struct {
		Id uint `uri:"id"`
	}{}
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		Error(ctx, 400, err)
		return
	}
	err = fn(ctx, id.Id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(ctx, 404, err)
		return
	}
	if err != nil {
		Error(ctx, 500, err)
		return
	}
	Success(ctx, nil)
}

func OneByIdApi[T any](ctx *gin.Context, fn func(context.Context, uint) (*T, error)) {
	id := struct {
		Id uint `uri:"id"`
	}{}
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		Error(ctx, 400, err)
		return
	}
	data, err := fn(ctx, id.Id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(ctx, 404, err)
		return
	}
	if err != nil {
		Error(ctx, 500, err)
		return
	}
	Success(ctx, data)
}

func UpdateByIdApi[T any](ctx *gin.Context, fn func(context.Context, uint, *T) error) {
	id := struct {
		Id uint `uri:"id"`
	}{}
	err := ctx.ShouldBindUri(&id)
	if err != nil {
		Error(ctx, 400, err)
		return
	}
	var t T
	err = ctx.ShouldBindJSON(&t)
	if err != nil {
		Error(ctx, 400, err)
		return
	}
	err = fn(ctx, id.Id, &t)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		Error(ctx, 404, err)
		return
	}
	if err != nil {
		Error(ctx, 500, err)
		return
	}
	Success(ctx, nil)
}
