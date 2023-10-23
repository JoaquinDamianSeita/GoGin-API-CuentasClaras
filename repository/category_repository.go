package repository

import (
	"GoGin-API-Base/dao"

	"gorm.io/gorm"
)

type CategoryRepository interface {
}

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func CategoryRepositoryInit(db *gorm.DB) *CategoryRepositoryImpl {
	db.AutoMigrate(&dao.Category{})
	return &CategoryRepositoryImpl{
		db: db,
	}
}
