package service

import "gorm.io/gorm"

type Filter interface {
	injectDB(db *gorm.DB) *gorm.DB
}

type Updater interface {
	updateDB(db *gorm.DB) *gorm.DB
}
