package repo

import (
	"context"

	"gorm.io/gorm"
)

var dbKey struct{}
var defaultDB *gorm.DB

func CtxWithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbKey, db.WithContext(ctx))
}

func GetDBFromCtx(ctx context.Context) *gorm.DB {
	val := ctx.Value(dbKey)
	if val == nil {
		return defaultDB
	}
	db, ok := val.(*gorm.DB)
	if !ok {
		return defaultDB
	}
	return db
}
