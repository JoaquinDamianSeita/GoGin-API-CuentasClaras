package repository

import (
	"GoGin-API-Base/dao"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindCategoryByOperation(operation dao.Operation) (dao.Category, error)
}

type CategoryRepositoryImpl struct {
	db *gorm.DB
}

func (u CategoryRepositoryImpl) FindCategoryByOperation(operation dao.Operation) (dao.Category, error) {
	u.db.Preload("Category").First(&operation)
	return operation.Category, nil
}

func CategoryRepositoryInit(db *gorm.DB) *CategoryRepositoryImpl {
	db.AutoMigrate(&dao.Category{})
	return &CategoryRepositoryImpl{
		db: db,
	}
}
