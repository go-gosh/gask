package repo

import (
	"context"
	"testing"

	"github.com/go-gosh/gask/app/model"
	"github.com/go-gosh/gestful/component/mapper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewTaskRepo(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db = db.Debug()
	assert.NoError(t, err)
	defer func() {
		s, _ := db.DB()
		err := s.Close()
		if err != nil {
			t.Error(err)
		}
	}()
	ctx := context.TODO()
	assert.NoError(t, db.AutoMigrate(&model.Task{}))
	repo := NewTaskRepo(db)
	for i := 0; i < 11; i++ {
		assert.NoError(t, repo.Create(ctx, &model.Task{}))
	}

	res, err := repo.Paginate(ctx, mapper.Paginator{
		StartId: 0,
		Limit:   10,
	}, mapper.EmptyWrapperFunc)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, res.StartId)
	assert.EqualValues(t, 10, res.Limit)
	assert.Equal(t, true, res.More)
	assert.Len(t, res.Data, 10)
	t.Logf("%+v", res)
	for i := 0; i < 10; i++ {
		assert.EqualValues(t, i+1, res.Data[i].ID)
	}

	res, err = repo.Paginate(ctx, mapper.Paginator{
		StartId: res.Data[9].ID,
		Limit:   10,
	}, mapper.EmptyWrapperFunc)
	assert.NoError(t, err)
	assert.EqualValues(t, 10, res.StartId)
	assert.EqualValues(t, 10, res.Limit)
	assert.Equal(t, false, res.More)
	assert.Len(t, res.Data, 1)
	assert.Len(t, res.Data, 1)
	assert.EqualValues(t, 11, res.Data[0].ID)

	assert.NoError(t, repo.DeleteById(ctx, 1))
	res, err = repo.Paginate(ctx, mapper.Paginator{
		StartId: 0,
		Limit:   11,
	}, mapper.EmptyWrapperFunc)
	assert.NoError(t, err)
	assert.EqualValues(t, 0, res.StartId)
	assert.EqualValues(t, 11, res.Limit)
	assert.Equal(t, false, res.More)
	assert.Len(t, res.Data, 10)
	assert.EqualValues(t, 2, res.Data[0].ID)
}
